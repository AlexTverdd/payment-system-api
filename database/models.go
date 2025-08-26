// Package database содержит функции и модели для работы с базой данных.
// Подключение к базе данных, миграции и начальную настройку.
package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wallet представляет модель кошелёка в базе данных.
// Содержит уникальный адрес и текущий баланс.
type Wallet struct {
	gorm.Model
	Address string `gorm:"unique;not null"` // Address уникальный адрес кошелька, используемый при идентификации
	Balance int64  //Balance текущий баланс кошелька
}

// Transaction представляет собой модель транзакции в базе данных
// Содержит адрес отправителя, адрес получателя, сумму, временную метку и UUID.
type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`        // идентификатор записи
	FromAddress string    `json:"from_address"`                // адрес отправителя
	ToAddress   string    `json:"to_address"`                  // адрес получателя
	Amount      int64     `json:"amount"`                      // сумма перевода
	Timestamp   time.Time `json:"timestamp"`                   // время создания транзакции
	UUID        string    `gorm:"unique;not null" json:"uuid"` // уникальный идентификатор транзакции
}

// BeforeCreate - метод, который автоматически генерирует уникальный UUID и
// устанавливает временную метку перед сохранением транзакции
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.UUID = uuid.New().String()
	t.Timestamp = time.Now()
	return
}
