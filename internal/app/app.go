package app

import (
	"github.com/xxehwuq/go-clean-architecture/config"
	"github.com/xxehwuq/go-clean-architecture/internal/delivery/fasthttp"
	"github.com/xxehwuq/go-clean-architecture/internal/delivery/fasthttp/v1"
	"github.com/xxehwuq/go-clean-architecture/internal/repository"
	"github.com/xxehwuq/go-clean-architecture/internal/usecase"
	"github.com/xxehwuq/go-clean-architecture/pkg/logger"
	"github.com/xxehwuq/go-clean-architecture/pkg/password"
	"github.com/xxehwuq/go-clean-architecture/pkg/postgres"
	"github.com/xxehwuq/go-clean-architecture/pkg/redis"
	"github.com/xxehwuq/go-clean-architecture/pkg/tokens"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level, logger.Console)
	l.Info("🚀 Starting %s", cfg.App.Name)

	// Postgres connection
	l.Info("Connecting to postgres database")

	pg, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		l.Fatal("error connecting to postgres database: %w", err)
	}
	defer pg.Close()

	// Redis
	l.Info("Connecting to redis database")

	rd, err := redis.New(cfg.Redis.URL)
	if err != nil {
		l.Fatal("error connecting to redis database: %w", err)
	}
	defer rd.Close()

	// Packages
	tokensManager := tokens.New(cfg.Tokens.SigningKey, cfg.Tokens.AccessTokenTTL)
	passwordHasher := password.NewHasher(cfg.Password.Salt)

	// Repositories
	userRepository := repository.NewUserRepository(pg, cfg.Postgres.Tables.Users)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepository, tokensManager, passwordHasher, rd, cfg.Tokens.RefreshTokenTTL)

	// FastHTTP
	fasthttpServer := fasthttp.NewServer(fasthttp.ServerConfig{
		Port:            cfg.HTTPServer.Port,
		AppName:         cfg.App.Name,
		ReadTimeout:     cfg.HTTPServer.ReadTimeout,
		WriteTimeout:    cfg.HTTPServer.WriteTimeout,
		IdleTimeout:     cfg.HTTPServer.IdleTimeout,
		ShutdownTimeout: cfg.HTTPServer.ShutdownTimeout,
	})

	v1.NewRouter(fasthttpServer, userUsecase)

	go func() {
		if err := fasthttpServer.Start(); err != nil {
			l.Error("error occurred when starting fasthttp server: %w")
		}
	}()

	l.Info("FastHTTP server was started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := fasthttpServer.Shutdown(); err != nil {
		l.Error("error occurred when stopping fasthttp server: %w")
	}
}
