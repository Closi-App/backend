package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"sync"
)

const configFilePath = "./config.yml"

type (
	Config struct {
		App AppConfig `yaml:"app" env-required:"true"`
		Log LogConfig `yaml:"log" env-required:"true"`
	}

	AppConfig struct {
		Name string `yaml:"name" env-required:"true"`
	}

	LogConfig struct {
		Level    string `yaml:"level" env-required:"true"`
		Encoding string `yaml:"encoding" env-required:"true"`
	}
)

var (
	once   sync.Once
	config Config
)

func Load() *Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(configFilePath, &config); err != nil {
			panic("error reading config file: " + err.Error())
		}
	})

	return &config
}
