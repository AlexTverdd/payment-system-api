package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит настрорйку приложения
type Config struct {
	DatabaseURL string
}

// LoadConfig загружает конфигурацию из .env файла или переменных окружения
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден. Используется перемееная окружения")

	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Переменная окружения DATABASE_URL не задана")
	}
	return &Config{
		DatabaseURL: dbURL,
	}
}
