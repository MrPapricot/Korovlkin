package config

import (
	"main/src/utils"
)

type Config struct {
	DatabaseUrl string
	AppPort     uint16
}

func LoadConfig() Config {
	app_port := utils.GetUint16OrDefault("APP_PORT", 8085)

	database_url := utils.GetStringOrDefault("DATABASE_URL", "postgres://postgres:1234@task-service-db:5432/task-db?sslmode=disable")

	return Config {
		AppPort: app_port,
		DatabaseUrl: database_url,
	}
}
