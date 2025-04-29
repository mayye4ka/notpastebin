package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	GrpcPort   int    `env:"GRPC_PORT"`
	HttpPort   int    `env:"HTTP_PORT"`
	DbAddr     string `env:"DB_ADDR"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
}

func ReadConfig() (*Config, error) {
	skipEnvLoad := false
	_, err := os.Open(".env")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		skipEnvLoad = true
	}
	if !skipEnvLoad {
		err = godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("load .env: %w", err)
		}
	}
	var config Config
	err = env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}
