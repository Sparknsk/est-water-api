package main

import (
	"context"
	"fmt"
	"github.com/ozonmp/est-water-api/internal/app/sender"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ozonmp/est-water-api/internal/app/repo"
	"github.com/ozonmp/est-water-api/internal/app/retranslator"
	"github.com/ozonmp/est-water-api/internal/config"
	"github.com/ozonmp/est-water-api/internal/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	// default: zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

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
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	cfgRetranslator := retranslator.Config{
		ChannelSize:   512,

		ConsumerCount: 10,
		ConsumeSize:   13,
		ConsumeTimeout: time.Millisecond*1000,

		ProducerCount: 1,
		WorkerCount:   1,
		WorkerBatchSize: 10,
		WorkerBatchTimeout: time.Millisecond*5000,

		Repo: repo.NewEventRepo(db),
		Sender: sender.NewEventSender(),
	}

	retranslator.NewRetranslator(cfgRetranslator).Start(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-sigs
}