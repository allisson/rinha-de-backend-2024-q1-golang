package rinha

import (
	"log/slog"
	"os"
	"path"

	"github.com/allisson/go-env"
	"github.com/joho/godotenv"
)

func searchup(dir string, filename string) string {
	if dir == "/" || dir == "" {
		return ""
	}

	if _, err := os.Stat(path.Join(dir, filename)); err == nil {
		return path.Join(dir, filename)
	}

	return searchup(path.Dir(dir), filename)
}

func findDotEnv() string {
	directory, err := os.Getwd()
	if err != nil {
		return ""
	}

	filename := ".env"
	return searchup(directory, filename)
}

func loadDotEnv() bool {
	dotenv := findDotEnv()
	if dotenv != "" {
		slog.Info("Found .env", "file", dotenv)
		if err := godotenv.Load(dotenv); err != nil {
			slog.Warn("Can't load .env", "file", dotenv, "error", err)
			return false
		}
		return true
	}
	return false
}

type Config struct {
	ServerHost       string
	ServerPort       uint
	DatabaseURL      string
	DatabaseMinConns uint
	DatabaseMaxConns uint
}

func NewConfig() *Config {
	loadDotEnv()

	return &Config{
		ServerHost:       env.GetString("RINHA_SERVER_HOST", "0.0.0.0"),
		ServerPort:       env.GetUint("RINHA_SERVER_PORT", 8000),
		DatabaseURL:      env.GetString("RINHA_DATABASE_URL", ""),
		DatabaseMinConns: env.GetUint("RINHA_DATABASE_MIN_CONNS", 5),
		DatabaseMaxConns: env.GetUint("RINHA_DATABASE_MAX_CONNS", 25),
	}
}
