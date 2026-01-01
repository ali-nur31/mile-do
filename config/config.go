package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB  Database
	Api Api
	Jwt Jwt
}

type Api struct {
	Port string `env:"API_PORT" env-default:":8080"`
}

type Jwt struct {
	AccessKey      string `env:"JWT_ACCESS_KEY" env-default:"sUp3r_k3y$123"`
	AccessExpMins  int    `env:"JWT_ACCESS_EXP_MINS" env-default:"10"`
	RefreshKey     string `env:"JWT_REFRESH_KEY" env-default:"sUp3r_k3y$321"`
	RefreshExpDays int    `env:"JWT_REFRESH_EXP_DAYS" env-default:"7"`
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
	jwt := jwtLoad()

	cfg.DB = db
	cfg.Api = api
	cfg.Jwt = jwt

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

func jwtLoad() Jwt {
	var jwt Jwt

	err := cleanenv.ReadEnv(&jwt)
	if err != nil {
		slog.Error("failed to load .env vars for JWT, using default values", "error", err)
	}

	return jwt
}

func databaseLoad() Database {
	var db Database

	err := cleanenv.ReadEnv(&db)
	if err != nil {
		slog.Error("failed to load .env vars for DB, using default values", "error", err)
	}

	return db
}
