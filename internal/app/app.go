package app

import (
	"github.com/xxehwuq/go-clean-architecture/config"
	"github.com/xxehwuq/go-clean-architecture/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level, logger.Console)

	l.Info("🚀 Starting %s", cfg.App.Name)
}
