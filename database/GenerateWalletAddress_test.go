package database

import (
	"testing"
)

// TestGenerateWalletAddress проверяет работу функции generateWalletAddress.
//
// Тест выполняет следующие проверки:
//   - Генерируются два адреса без ошибок.
//   - Адреса являются уникальными (не совпадают).
//   - Длина каждого адреса соответствует ожидаемой (64 символа).
func TestGenerateWalletAddress(t *testing.T) {
	// Генерируем два адреса
	address1, err1 := generateWalletAddress()
	if err1 != nil {
		t.Fatalf("Не удалось сгенерировать первый адрес: %v", err1)
	}
	t.Logf("Сгенерирован первый адрес: %s", address1)

	address2, err2 := generateWalletAddress()
	if err2 != nil {
		t.Fatalf("Не удалось сгенерировать второй адрес: %v", err2)
	}
	t.Logf("Сгенерирован первый адрес: %s", address2)

	// Проверяем, что оба адреса уникальны
	if address1 == address2 {
		t.Errorf("Ожидалась уникальность адресов, но они одинаковы")
	}

	// Проверяем, что длина адресов правильная (64 символа)
	const expectedLength = 64
	if len(address1) != expectedLength {
		t.Errorf("Неправильная длина первого адреса. Ожидалось %d, получено %d", expectedLength, len(address1))
	}
	if len(address2) != expectedLength {
		t.Errorf("Неправильная длина второго адреса. Ожидалось %d, получено %d", expectedLength, len(address2))
	}
}
