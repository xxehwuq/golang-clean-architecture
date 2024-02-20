package app

import (
	"github.com/xxehwuq/go-clean-architecture/config"
	"github.com/xxehwuq/go-clean-architecture/pkg/logger"
	"github.com/xxehwuq/go-clean-architecture/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level, logger.Console)
	l.Info("🚀 Starting %s", cfg.App.Name)

	pg, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatal("error connecting to postgres: %w", err)
	}
	defer pg.Close()
}
