package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB  Database
	Api Api
}

type Api struct {
	Port         string `env:"API_PORT" env-default:":8080"`
	JWTSecretKey string `env:"JWT_SECRET_KEY" env-default:"sUp3r_k3y$123"`
}

type Database struct {
	Port     string `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Username string `env:"DB_USERNAME" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `env:"DB_NAME" env-default:"mile_do_db"`
}

func MustLoad() *Config {
	var cfg Config

	db := databaseLoad()
	api := apiLoad()

	cfg.DB = db
	cfg.Api = api

	return &cfg
}

func apiLoad() Api {
	var api Api

	err := cleanenv.ReadEnv(&api)
	if err != nil {
		slog.Error("failed to load .env vars for API, using default values", "error", err)
	}

	return api
}

func databaseLoad() Database {
	var db Database

	err := cleanenv.ReadEnv(&db)
	if err != nil {
		slog.Error("failed to load .env vars for DB, using default values", "error", err)
	}

	return db
}
