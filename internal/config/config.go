package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	DB    db
	App   app
	Aws   aws
	Minio minio
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		App:   loadApplicationConfig(),
		Aws:   loadAwsConfig(),
		DB:    loadDatabaseConfig(),
		Minio: loadMinioConfig(),
	}, nil
}
