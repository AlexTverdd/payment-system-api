package database

import (
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
