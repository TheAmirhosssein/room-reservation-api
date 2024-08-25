package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		APP  `yaml:"app"`
		HTTP `yaml:"http"`
	}
	APP struct {
		Name      string `env-required:"true" yaml:"name"`
		Version   string `env-required:"true" yaml:"version"`
		SecretKey string `env-required:"true" env:"SECRET_KEY"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
	}
)

func NewConfig() (*Config, error) {
	conf := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", conf)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(conf)
	if err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}
	return conf, nil
}
