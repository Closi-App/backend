package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"sync"
	"time"
)

type (
	Config struct {
		App   AppConfig   `yaml:"app" env-required:"true"`
		Log   LogConfig   `yaml:"log" env-required:"true"`
		Mongo MongoConfig `yaml:"mongo" env-required:"true"`
		HTTP  HTTPConfig  `yaml:"http" env-required:"true"`
	}

	AppConfig struct {
		Name string `yaml:"name" env-required:"true"`
	}

	LogConfig struct {
		Level    string `yaml:"level" env-required:"true"`
		Encoding string `yaml:"encoding" env-required:"true"`
	}

	MongoConfig struct {
		URI    string `yaml:"uri" env:"MONGO_URI" env-required:"true"`
		DBName string `yaml:"db_name" env:"MONGO_DB_NAME" env-required:"true"`
	}

	HTTPConfig struct {
		Host         string        `yaml:"host" env:"HTTP_HOST" env-required:"true"`
		Port         int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write_timeout" env-required:"true"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" env-required:"true"`
	}
)

var (
	once   sync.Once
	config Config
)

func Load(configFilePath string) *Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(configFilePath, &config); err != nil {
			panic("error reading config file: " + err.Error())
		}
	})

	return &config
}
