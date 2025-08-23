package database

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB устанавливает соединение с базой данных
func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	log.Println("Подключение к базе данных выполнено успешно")
}

// Migrate создает таблицы если они не существуют
func Migrate() {
	err := DB.AutoMigrate(&Wallet{}, &Transaction{})
	if err != nil {
		log.Fatalf("Миграция базы данных не удалась %v", err)
	}
	log.Println("Миграция базы данных успешна")
}

// generateWalletAddress генерирует адресс
func generateWalletAddress() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read((bytes))
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении случайных байт: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
