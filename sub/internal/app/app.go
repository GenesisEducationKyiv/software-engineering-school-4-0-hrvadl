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
	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/queue"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sender"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sender/formatter"
	subs "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sub"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/validator"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/grpc/clients/ratewatcher"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/grpc/server/sub"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/queue/publisher/mailer"
)

const operation = "app init"

const (
	cronJobHour    = 12
	cronJobMinute  = 0o0
	cronJobTimeout = time.Minute * 1
	sendTimeout    = time.Second * 5
)

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
		sub.NewErrorMappingInterceptor(),
	))

	db, err := db.NewConn(a.cfg.Dsn)
	if err != nil {
		return fmt.Errorf("%s: failed to init db: %w", operation, err)
	}

	sr := subscriber.NewRepo(db)
	v := validator.NewStdlib()
	svc := subs.NewService(sr, v)
	sub.Register(a.srv, svc, a.log.With(slog.String("source", "sub")))

	sg := subscriber.NewRepo(db)
	fmter := formatter.NewWithDate()
	rw, err := ratewatcher.NewClient(
		a.cfg.RateWatcherAddr,
		a.log.With(slog.String("source", "rateWatcher")),
	)
	if err != nil {
		return fmt.Errorf("%s: failed to connect to rate watcher: %w", operation, err)
	}

	if a.nats, err = nats.Connect(a.cfg.NatsURL); err != nil {
		return fmt.Errorf("%s: failed to connect to nats server: %w", operation, err)
	}

	m := mailer.NewClient(
		queue.NewAdapter(a.nats, sendTimeout),
		a.log.With(slog.String("source", "mailer")),
	)

	mailSender := sender.New(
		m,
		sg,
		fmter,
		rw,
		a.log.With(slog.String("source", "cron sender")),
	)

	cronAdapter := sender.NewCronJobAdapter(
		mailSender,
		cronJobTimeout,
		a.log.With(slog.String("source", "adapter")),
	)
	job := cron.NewDailyJob(cronJobHour, cronJobMinute, a.log.With(slog.String("source", "cron")))
	job.Do(cronAdapter)

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(a.srv, healthcheck)
	healthcheck.SetServingStatus(
		pb.RateWatcherService_ServiceDesc.ServiceName,
		healthgrpc.HealthCheckResponse_SERVING,
	)

	l, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to start listener on port %s: %w", operation, a.cfg.Port, err)
	}

	return a.srv.Serve(l)
}

// GracefulStop method gracefully stop the server. It listens to the OS sigals.
// After it receives signal it terminates all currently active servers,
// client, connections (if any) and gracefully exits.
func (a *App) GracefulStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	signal := <-ch
	a.log.Info("Received stop signal. Terminating...", slog.Any("signal", signal))
	a.srv.Stop()
	a.nats.Close()
	a.log.Info("Successfully terminated server. Bye!")
}
