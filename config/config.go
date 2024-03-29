package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"sync"
	"time"
)

type (
	Config struct {
		App        App        `yaml:"app" env-required:"true"`
		Log        Log        `yaml:"log" env-required:"true"`
		Postgres   Postgres   `yaml:"postgres" env-required:"true"`
		Tokens     Tokens     `yaml:"tokens" env-required:"true"`
		Password   Password   `yaml:"password" env-required:"true"`
		Redis      Redis      `yaml:"redis" env-required:"true"`
		HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
	}

	App struct {
		Name string `yaml:"name" env-required:"true"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true"`
	}

	Postgres struct {
		URL    string `env:"POSTGRES_URL" env-required:"true"`
		Tables struct {
			Users string `yaml:"users" env-required:"true"`
		} `yaml:"tables" env-required:"true"`
	}

	Tokens struct {
		SigningKey      string        `env:"TOKENS_SIGNING_KEY" env-required:"true"`
		AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
		RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	}

	Password struct {
		Salt string `env:"PASSWORD_SALT" env-required:"true"`
	}

	Redis struct {
		URL string `env:"REDIS_URL" env-required:"true"`
	}

	HTTPServer struct {
		Port            int           `yaml:"port" env-required:"true"`
		ReadTimeout     time.Duration `yaml:"read_timeout" env-required:"true"`
		WriteTimeout    time.Duration `yaml:"write_timeout" env-required:"true"`
		IdleTimeout     time.Duration `yaml:"idle_timeout" env-required:"true"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-required:"true"`
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
