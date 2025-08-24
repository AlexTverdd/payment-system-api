package database

import (
	"log"
	"os"
	"testing"
)

// testDSN содержит строку подключения к тестовой базе данных.
// Используется для инициализации соединения в тестах.
const testDSN = "host=127.0.0.1 user=testuser password=testpass dbname=testdb port=5433 sslmode=disable"

// TestMain выполняет подготовку тестового окружения.
// Она подключается к тестовой базе данных, выполняет миграции
// и затем запускает все тесты.
// После завершения тестов происходит завершение программы с соответствующим кодом выхода.
func TestMain(m *testing.M) {
	ConnectDB(testDSN)
	if DB == nil {
		log.Fatalf("Не удалось подключиться к тестовой базе данных.")
	}

	Migrate()

	code := m.Run()

	os.Exit(code)
}

// TestInitialSetup проверяет работу функции InitialSetup.
// Перед вызовом InitialSetup база данных должна быть пустой по таблице Wallet.
// После вызова InitialSetup должно создаться ровно 10 кошельков.
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
