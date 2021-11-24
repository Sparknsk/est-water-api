package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/app/retranslator"
	"github.com/ozonmp/est-water-api/internal/app/sender"
	"github.com/ozonmp/est-water-api/internal/config"
	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.ReadConfigYML("retranslator-config.yml"); err != nil {
		logger.FatalKV(ctx, "Failed init configuration",
			"err", errors.Wrap(err, "config.ReadConfigYML() failed"),
		)
	}
	cfg := config.GetConfigInstance()

	syncLogger, err := logger.NewLogger(cfg)
	if err != nil {
		logger.FatalKV(ctx, "Failed init logger",
			"err", errors.Wrap(err, "logger.NewLogger() failed"),
		)
	}
	defer syncLogger()

	logger.InfoKV(ctx, fmt.Sprintf("Starting service: %s", cfg.Project.Name),
		"version", cfg.Project.Version,
		"commitHash", cfg.Project.CommitHash,
		"debug", cfg.Logging.IsDebug(),
		"environment", cfg.Project.Environment,
		"Starting service: %s", cfg.Project.Name,
	)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(ctx, dsn, cfg.Database.Driver)
	if err != nil {
		logger.FatalKV(ctx, "Failed init postgres",
			"err", errors.Wrap(err, "database.NewPostgres() failed"),
		)
	}
	defer db.Close()

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsAddr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)
	metricsServer := &http.Server{
		Addr: metricsAddr,
		Handler: mux,
	}

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Metrics server is running on %s", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running metrics server",
				"err", errors.Wrap(err, "metricsServer.ListenAndServe() failed"),
			)
			cancel()
		}
	}()

	cfgRetranslator := retranslator.Config{
		ChannelSize: 512,

		ConsumerCount: 10,
		ConsumeSize: 3,
		ConsumeTimeout: time.Millisecond*1000,

		ProducerCount: 1,
		WorkerCount: 1,
		WorkerBatchSize: 4,
		WorkerBatchTimeout: time.Millisecond*5000,

		Repo: repo.NewEventRepo(db),
		Sender: sender.NewEventSender(),
	}

	retranslator.NewRetranslator(cfgRetranslator).Start(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-sigs

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "metricsServer.Shutdown",
			"err", errors.Wrap(err, "metricsServer.Shutdown() failed"),
		)
	} else {
		logger.InfoKV(ctx, "metricsServer shut down correctly")
	}
}