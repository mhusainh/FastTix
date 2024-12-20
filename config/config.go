package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV            string         `env:"ENV" envDefault:"dev"`
	PORT           string         `env:"PORT" envDefault:"8081"`
	JWTConfig      JWTConfig      `envPrefix:"JWT_"`
	PostgresConfig PostgresConfig `envPrefix:"POSTGRES_"`
	SMTPConfig     SMTPConfig     `envPrefix:"SMTP_"`
}

type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"587"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

type JWTConfig struct {
	SecretKey string `env:"SECRET_KEY"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
