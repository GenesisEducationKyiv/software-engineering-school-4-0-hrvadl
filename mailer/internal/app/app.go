package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/grpc/server/mailer"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/pkg/logger"
)

const operation = "app init"

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
	cfg cfg.Config
	log *slog.Logger
	srv *grpc.Server
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
// after that initializes all neccessary domain related services and finally
// starts listening on the provided ports. Could return an error if any of
// described above steps failed
func (a *App) Run() error {
	a.srv = grpc.NewServer(grpc.ChainUnaryInterceptor(
		logger.NewServerGRPCMiddleware(a.log),
	))

	mailer.Register(a.srv, resend.NewClient(a.cfg.MailerToken), a.log.With("source", "mailerSrv"))
	listener, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen on port %s: %w", operation, a.cfg.Port, err)
	}

	return a.srv.Serve(listener)
}

// GracefulStop method gracefully stop the server. It listens to the OS sigals.
// After it recieves signal it terminates all currently active servers,
// client, connections (if any) and gracefully exits.
func (a *App) GracefulStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	signal := <-ch
	a.log.Info("Recieved stop signal. Terminating...", "signal", signal)
	a.srv.Stop()
	a.log.Info("Successfully terminated server. Bye!")
}
