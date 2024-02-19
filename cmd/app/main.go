package main

import (
	"fmt"
	"github.com/xxehwuq/go-clean-architecture/config"
)

func main() {
	cfg := config.New()
	fmt.Println(cfg.App.Name)
}
