package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Port     string `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Username string `env:"DB_USERNAME" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `env:"DB_NAME" env-default:"taskee_db"`
}

func MustLoad() (Database, error) {
	var cfg Database

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("failed to load .env file, using default values", "error", err)
		return cfg, err
	}

	return cfg, nil
}
