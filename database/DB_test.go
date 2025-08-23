package database

import (
	"log"
	"os"
	"testing"
)

// dsn - строка подключения к тестовой базе данных
const testDSN = "host=127.0.0.1 user=testuser password=testpass dbname=testdb port=5433 sslmode=disable"

// TestMain используется для настройки и очистки ресурсов для тестов
func TestMain(m *testing.M) {
	ConnectDB(testDSN)
	if DB == nil {
		log.Fatalf("Не удалось подключиться к тестовой базе данных.")
	}

	Migrate()

	code := m.Run()

	os.Exit(code)
}

// TestConnectAndMigrate проверяет, что соединение с БД и миграции работают без ошибок
func TestConnectAndMigrate(t *testing.T) {
	if DB == nil {
		t.Fatalf("Переменная DB не была инициализирована в TestMain.")
	}

	if !DB.Migrator().HasTable(&Wallet{}) {
		t.Fatalf("Таблица Wallet не была создана.")
	}
	if !DB.Migrator().HasTable(&Transaction{}) {
		t.Fatalf("Таблица Transaction не была создана.")
	}
}

// TestInitialSetup проверяет, что функция InitialSetup работает корректно
func TestInitialSetup(t *testing.T) {
	var count int64
	DB.Model(&Wallet{}).Count(&count)
	if count != 0 {
		t.Errorf("Перед запуском теста в базе уже были кошельки.")
	}

	InitialSetup()
	DB.Model(&Wallet{}).Count(&count)
	if count != 10 {
		t.Errorf("Ожидалось 10 кошельков, найдено %d.", count)
	}
}
