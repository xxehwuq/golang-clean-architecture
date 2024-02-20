package app

import (
	"github.com/xxehwuq/go-clean-architecture/config"
	"github.com/xxehwuq/go-clean-architecture/internal/repository"
	"github.com/xxehwuq/go-clean-architecture/pkg/logger"
	"github.com/xxehwuq/go-clean-architecture/pkg/postgres"
	"github.com/xxehwuq/go-clean-architecture/pkg/tokens"
)

func Run(cfg *config.Config) {
	// Logger initializing
	l := logger.New(cfg.Log.Level, logger.Console)
	l.Info("🚀 Starting %s", cfg.App.Name)

	// Postgres connection
	pg, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatal("error connecting to postgres: %w", err)
	}
	defer pg.Close()

	// Packages
	tokens.New(cfg.Tokens.SigningKey, cfg.Tokens.AccessTokenTTL, cfg.Tokens.RefreshTokenTTL)

	// Repositories
	repository.NewUserRepository(pg, cfg.Postgres.Tables.Users)
}
