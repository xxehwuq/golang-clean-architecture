package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"sync"
)

type (
	Config struct {
		App App `yaml:"app" env-required:"true"`
	}

	App struct {
		Name string `yaml:"name" env-required:"true"`
	}
)

var config Config
var once sync.Once

func New() *Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig("./config/config.yml", &config); err != nil {
			panic("error reading config: " + err.Error())
		}
	})

	return &config
}
