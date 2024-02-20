package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"sync"
)

type (
	Config struct {
		App      App      `yaml:"app" env-required:"true"`
		Log      Log      `env-required:"true"`
		Postgres Postgres `env-required:"true"`
	}

	App struct {
		Name string `yaml:"name" env-required:"true"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true"`
	}

	Postgres struct {
		URL string `env:"POSTGRES_URL" env-required:"true"`
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
