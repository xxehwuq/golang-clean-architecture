package main

import (
	"github.com/xxehwuq/go-clean-architecture/config"
	"github.com/xxehwuq/go-clean-architecture/internal/app"
)

func main() {
	// Configuration
	cfg := config.New()

	// App
	app.Run(cfg)
}
