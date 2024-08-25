package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		APP  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}
	APP struct {
		Name      string `env-required:"true" yaml:"name"`
		Version   string `env-required:"true" yaml:"version"`
		SecretKey string `env:"SECRET_KEY,required"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
		Host string `env-required:"true" yaml:"host"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"`
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
