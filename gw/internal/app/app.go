package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/logger"
	rw "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/cfg"
)

const operation = "app init"

const shutdownTimeout = 5 * time.Second

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
	srv *http.Server
	cfg cfg.Config
	log *slog.Logger
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
// described above steps failed
func (a *App) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rm := newResponseMapper(a.log.With(slog.String("source", "responseMapper")))
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(rm.mapGRPCErr),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, newMarshaller()),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := rw.RegisterRateWatcherServiceHandlerFromEndpoint(ctx, mux, a.cfg.RateWatcherAddr, opts); err != nil {
		return fmt.Errorf("%s: failed to create proxy for rw grpc svc: %w", operation, err)
	}

	if err := sub.RegisterSubServiceHandlerFromEndpoint(ctx, mux, a.cfg.SubAddr, opts); err != nil {
		return fmt.Errorf("%s: failed to create proxy for sub grpc svc: %w", operation, err)
	}

	s := newServer(
		middleware.Logger(mux),
		a.cfg.Addr,
		slog.NewLogLogger(a.log.Handler(), logger.MapLevels(a.cfg.LogLevel)),
	)

	return s.ListenAndServe()
}

// GracefulStop method gracefully stop the server. It listens to the OS sigals.
// After it receives signal it terminates all currently active servers,
// client, connections (if any) and gracefully exits.
func (a *App) GracefulStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	signal := <-ch
	a.log.Info("Received stop signal. Terminating...", slog.Any("signal", signal))
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.srv.Shutdown(ctx); err != nil {
		a.log.Error("Failed to gracefully stop the server", slog.Any("err", err))
		return
	}
	a.log.Info("Successfully terminated server. Bye!")
}
