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

// generateWalletAddress генерирует адрес
func generateWalletAddress() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read((bytes))
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении случайных байт: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func InitialSetup() {
	var count int64
	DB.Model(&Wallet{}).Count(&count)
	if count == 0 {
		fmt.Println("Создание начальных кошельков")
		wallets := make([]Wallet, 10)
		for i := 0; i < 10; i++ {
			address, err := generateWalletAddress()
			if err != nil {
				log.Fatalf("Не удалось сгенерировать адрес кошелька: %v", err)
			}
			wallets[i] = Wallet{
				Address: address,
				Balance: 100.0,
			}
		}
		if err := DB.Create(&wallets).Error; err != nil {
			log.Fatalf("Не удалось создать начальные кошельки: %v", err)
		}
		fmt.Println("10 начальных кошельков созданы")
	} else {
		fmt.Println("Кошельков уже существует", count)
	}
}
