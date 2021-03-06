package configs

import "os"

const (
	prod = "production"
)

// Config object
type Config struct {
	Env      string         `env:"ENV"`
	Postgres PostgresConfig `json:"postgres"`
	Host     string         `env:"APP_HOST"`
	Port     string         `env:"APP_PORT"`
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Env:      os.Getenv("ENV"),
		Postgres: GetPostgresConfig(),
		Host:     os.Getenv("APP_HOST"),
		Port:     os.Getenv("APP_PORT"),
	}
}
