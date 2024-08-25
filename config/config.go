package config

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
