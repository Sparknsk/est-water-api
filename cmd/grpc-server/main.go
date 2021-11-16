package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ozonmp/est-water-api/internal/config"
	"github.com/ozonmp/est-water-api/internal/database"
	"github.com/ozonmp/est-water-api/internal/logger"
	"github.com/ozonmp/est-water-api/internal/server"
	"github.com/ozonmp/est-water-api/internal/tracer"
)

var (
	batchSize uint = 2
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := config.ReadConfigYML("config.yml"); err != nil {
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

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

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

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.FatalKV(ctx, "Failed migrations",
				"err", errors.Wrap(err, "goose.Up() failed"),
			)
		}
	}

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		logger.FatalKV(ctx, "Failed init tracing",
			"err", errors.Wrap(err, "decoder.Decode() failed"),
		)
	}
	defer tracing.Close()
	logger.InfoKV(ctx, "Tracing started")

	if err := server.NewGrpcServer(&cfg, db, batchSize).Start(); err != nil {
		logger.FatalKV(ctx, "Failed creating gRPC server",
			"err", errors.Wrap(err, "server.NewGrpcServer() failed"),
		)
	}
}
