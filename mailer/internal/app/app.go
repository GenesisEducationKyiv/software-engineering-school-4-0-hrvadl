package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/logger"
	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/gomail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/mailer"
)

const operation = "app init"

const mailerTimeout = time.Second * 5

// New constructs new App with provided arguments.
// NOTE: than neither cfg or log can't be nil or App will panic.
func New(cfg cfg.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

// App is a thin abstraction used to initialize all the dependencies,
// db connections, and GRPC server/clients. Could return an error if any
// of described above steps failed.
type App struct {
	cfg  cfg.Config
	log  *slog.Logger
	srv  *grpc.Server
	nats *nats.Conn
}

// MustRun is a wrapper around App.Run() function which could be handly
// when it's called from the main goroutine and we don't need to handler
// an error.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run method creates new GRPC server then initializes MySQL DB connection,
// after that initializes all necessary domain related services and finally
// starts listening on the provided ports. Could return an error if any of
// described above steps failed.
func (a *App) Run() error {
	a.srv = grpc.NewServer(grpc.ChainUnaryInterceptor(
		logger.NewServerGRPCMiddleware(a.log),
	))

	resend := resend.NewClient(a.cfg.MailerFromFallback, a.cfg.MailerFallbackToken)
	gomail := gomail.NewClient(
		a.cfg.MailerFrom,
		a.cfg.MailerPassword,
		a.cfg.MailerHost,
		a.cfg.MailerPort,
	)

	mailSvc := mail.NewService(gomail)
	mailSvc.SetNext(resend)

	var err error
	if a.nats, err = nats.Connect(a.cfg.NatsURL); err != nil {
		return fmt.Errorf("%s: failed to connect to nats: %w", operation, err)
	}

	js, err := jetstream.New(a.nats)
	if err != nil {
		return fmt.Errorf("%s: failed to connect to jetstream: %w", operation, err)
	}

	if err = mailer.NewSub(js, a.log).Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to sub to deb: %w", operation, err)
	}

	m := mailer.New(a.nats, mailSvc, a.log.With(slog.String("source", "mailerSrv")), mailerTimeout)
	if err = m.Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to subscribe: %w", operation, err)
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(a.srv, healthcheck)
	healthcheck.SetServingStatus(
		pb.MailerService_ServiceDesc.ServiceName,
		healthgrpc.HealthCheckResponse_SERVING,
	)

	listener, err := net.Listen("tcp", net.JoinHostPort(a.cfg.Host, a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen on port %s: %w", operation, a.cfg.Port, err)
	}

	return a.srv.Serve(listener)
}

// GracefulStop method gracefully stop the server. It listens to the OS sigals.
// After it receives signal it terminates all currently active servers,
// client, connections (if any) and gracefully exits.
func (a *App) GracefulStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	signal := <-ch
	a.log.Info("Received stop signal. Terminating...", slog.Any("signal", signal))
	a.Stop()
}

func (a *App) Stop() {
	a.srv.Stop()
	a.nats.Close()
	a.log.Info("Successfully terminated server. Bye!")
}
