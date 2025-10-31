package config

import (
	"log"
	"os"
)

type Config struct {
	Dbname string
	Port   string
}

func Load() Config {
	cfg := Config{
		Dbname: os.Getenv("DB_NAME"),
		Port:   os.Getenv("PORT"),
	}
	if cfg.Dbname == "" {
		log.Fatal("DB_NAME is required")
	}
	if cfg.Port == "" {
		log.Fatal("PORT is required")
	}

	return cfg
}
