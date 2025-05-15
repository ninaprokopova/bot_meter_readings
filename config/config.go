package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
}

func LoadConfig() (*Config, error) {
	return &Config{
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER"),
		DBPassword:    getEnv("DB_PASSWORD"),
		DBName:        getEnv("DB_NAME"),
	}, nil
}

func getEnv(key string, defaults ...string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error: .env file not found")
	}
	value := os.Getenv(key)
	if value == "" && len(defaults) > 0 {
		return defaults[0]
	}
	if value == "" {
		panic(fmt.Sprintf("%v is not set", key))
	}
	return value
}
