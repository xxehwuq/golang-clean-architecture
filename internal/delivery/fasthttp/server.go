package fasthttp

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Server struct {
	*fiber.App

	port            int
	shutdownTimeout time.Duration
}

type ServerConfig struct {
	Port            int
	AppName         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func NewServer(cfg ServerConfig) *Server {
	srv := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		AppName:      cfg.AppName,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	return &Server{
		App:             srv,
		port:            cfg.Port,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Start() error {
	return s.Listen(fmt.Sprintf(":%d", s.port))
}

func (s *Server) Shutdown() error {
	return s.ShutdownWithTimeout(s.shutdownTimeout)
}
