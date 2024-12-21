package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV            string         `env:"ENV" envDefault:"dev"`
	PORT           string         `env:"PORT" envDefault:"8081"`
	JWTConfig      JWTConfig      `envPrefix:"JWT_"`
	MySQLConfig    MySQLConfig    `envPrefix:"MYSQL_"`
	SMTPConfig     SMTPConfig     `envPrefix:"SMTP_"`
	MidtransConfig MidtransConfig `envPreflix:"MIDTRANS_"`
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

type MySQLConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"3306"`
	User     string `env:"USER" envDefault:"root"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type MidtransConfig struct {
	BaseURL   string `env:"BASE_URL"`
	ClientKey string `env:"CLIENT_KEY"`
	ServerKey string `env:"SERVER_KEY"`
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
