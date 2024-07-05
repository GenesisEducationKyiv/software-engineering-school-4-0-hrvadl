package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	runner "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/gomail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
	rateSub "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/rate"
	subSub "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/sub"
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
	db   *db.Conn
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var err error
	if a.db, err = db.NewConn(ctx, a.cfg.MongoURL); err != nil {
		return fmt.Errorf("%s: failed to connect to mongo: %w", operation, err)
	}

	subscriberRepo := subscriber.NewRepo(a.db.GetDB())
	rateRepo := rate.NewRepo(a.db.GetDB())

	resend := resend.NewClient(a.cfg.MailerFromFallback, a.cfg.MailerFallbackToken)
	gomail := gomail.NewClient(
		a.cfg.MailerFrom,
		a.cfg.MailerPassword,
		a.cfg.MailerHost,
		a.cfg.MailerPort,
	)

	mailSvc := mail.NewService(gomail)
	mailSvc.SetNext(resend)

	if a.nats, err = nats.Connect(a.cfg.NatsURL); err != nil {
		return fmt.Errorf("%s: failed to connect to nats: %w", operation, err)
	}

	js, err := jetstream.New(a.nats)
	if err != nil {
		return fmt.Errorf("%s: failed to connect to jetstream: %w", operation, err)
	}

	subSubscriber := subSub.NewSubscriber(js, subscriberRepo, a.log, time.Second*5)
	if err = subSubscriber.Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to sub to CDC: %w", operation, err)
	}

	// TODO: create fucking services
	m := rateSub.NewSubscriber(
		a.nats,
		rateRepo,
		a.log.With(slog.String("source", "mailerSrv")),
		mailerTimeout,
	)
	if err = m.Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to subscribe: %w", operation, err)
	}

	adp := cron.NewAdapter(rateRepo, subscriberRepo, mailSvc, time.Second*5)

	job := runner.NewDailyJob(19, 23, a.log)
	job.Do(adp)

	return nil
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
