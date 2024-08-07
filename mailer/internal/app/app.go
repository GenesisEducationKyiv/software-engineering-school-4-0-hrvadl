package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	runner "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/metrics"
	"github.com/nats-io/nats.go"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/gomail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/cron/formatter"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/subscriber"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	raterepo "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	subrepo "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
	rateSub "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/rate"
	subSub "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/sub"
)

const operation = "app init"

const (
	sendHours   = 12
	sendMinutes = 0o0
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
	cfg     cfg.Config
	log     *slog.Logger
	nats    *nats.Conn
	db      *db.Conn
	metrics *metrics.Engine
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
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ConnectTimeout)
	defer cancel()

	var err error
	if a.db, err = db.NewConn(ctx, a.cfg.MongoURL); err != nil {
		return fmt.Errorf("%s: failed to connect to mongo: %w", operation, err)
	}

	subscriberRepo := subrepo.NewRepo(a.db.GetDB())
	subSvc := subscriber.NewService(subscriberRepo)
	rateRepo := raterepo.NewRepo(a.db.GetDB())
	rateSvc := rate.NewService(rateRepo)

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

	subSubscriber := subSub.NewSubscriber(a.nats, subSvc, a.log, a.cfg.ConnectTimeout)
	if err = subSubscriber.Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to sub to CDC: %w", operation, err)
	}

	m := rateSub.NewSubscriber(
		a.nats,
		rateSvc,
		a.log.With(slog.String("source", "mailerSrv")),
		a.cfg.ConnectTimeout,
	)
	if err = m.Subscribe(); err != nil {
		return fmt.Errorf("%s: failed to subscribe: %w", operation, err)
	}

	a.metrics = metrics.NewEngine(net.JoinHostPort(a.cfg.Host, a.cfg.PrometheusPort))
	if err := a.metrics.Register(runner.GetMetrics()...); err != nil {
		return fmt.Errorf("%s: failed to register metrics: %w", operation, err)
	}

	adp := cron.NewAdapter(
		rateSvc,
		subSvc,
		mailSvc,
		formatter.NewWithDate(),
		a.cfg.ConnectTimeout,
		a.log,
	)
	job := runner.NewDailyJob(sendHours, sendMinutes, a.log)
	job.Do(runner.NewWithMetrics(adp, "mail"))

	go func() {
		if err := a.metrics.Start(); err != nil {
			a.log.Error("Failed to serve metrics", slog.Any("err", err))
		}
	}()

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
	a.nats.Close()

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ConnectTimeout)
	defer cancel()

	if err := a.db.Close(ctx); err != nil {
		a.log.Error("Failed to gracefully stop mongo", slog.Any("err", err))
	}

	a.log.Info("Successfully terminated server. Bye!")
}
