// Package config предоставляет функционал для загрузки
// конфигурации приложения из .env файла или переменных окружения.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит настройки приложения, включая URL базы данных.
type Config struct {
	DatabaseURL string
}

// LoadConfig загружает конфигурацию приложения.
//
// Она выполняет следующие действия:
// 1. Пытается загрузить файл .env (если он существует).
// 2. Считывает переменную окружения DATABASE_URL.
// 3. Если DATABASE_URL не задана, завершает работу с ошибкой.
// Возвращает указатель на структуру Config с загруженными значениями.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден. Используется переменная окружения")

	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("Переменная окружения DATABASE_URL не задана")
	}
	return &Config{
		DatabaseURL: dbURL,
	}
}
