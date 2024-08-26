package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		APP  `yaml:"app"`
		HTTP `yaml:"http"`
		DB   `yaml:"db"`
	}
	APP struct {
		Name      string `env-required:"true" yaml:"name"`
		Version   string `env-required:"true" yaml:"version"`
		SecretKey string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
		Host string `env-required:"true" yaml:"host"`
	}

	DB struct {
		Host     string `env-required:"true" env:"POSTGRES_HOST"`
		Username string `env-required:"true" env:"POSTGRES_USER"`
		Password string `env-required:"true" env:"POSTGRES_PASSWORD"`
		DB       string `env-required:"true" env:"POSTGRES_DB"`
	}
)

func NewConfig() (*Config, error) {
	conf := &Config{}
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	err = cleanenv.ReadConfig("./config/config.yaml", conf)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(conf)
	if err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}
	return conf, nil
}
