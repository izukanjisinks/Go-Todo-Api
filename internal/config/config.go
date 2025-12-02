package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBServer            string
	DBPort              string
	DBName              string
	DBTrustedConnection string
	ServerPort          string
}

func Load() (*Config, error) {
	// Load .env file (ignore error in production where env vars are set directly)
	godotenv.Load()

	return &Config{
		DBServer:            getEnv("DB_SERVER", "localhost"),
		DBPort:              getEnv("DB_PORT", "1433"),
		DBName:              getEnv("DB_NAME", "Todos"),
		DBTrustedConnection: getEnv("DB_TRUSTED_CONNECTION", "yes"),
		ServerPort:          getEnv("SERVER_PORT", "3000"),
	}, nil
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("server=%s,%s;database=%s;trusted_connection=%s",
		c.DBServer, c.DBPort, c.DBName, c.DBTrustedConnection)
}

func getEnv(key, defaultValue string) string {

	value := os.Getenv(key)

	if value != "" {
		return value
	}
	
	return defaultValue
}
