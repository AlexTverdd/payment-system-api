package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wallet предстваляет собой модель кошелька в базе данных
type Wallet struct {
	gorm.Model
	// Address уникальный адрес кошелька, используемый при идентефикации
	Address string `gorm:"unique;not null"`
	//Balance текущий баланс кошелька
	Balance float64
}

// Transaction представляет собой модель транзакции в базе данных
type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Amount      float64   `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
	UUID        string    `gorm:"unique;not null" json:"uuid"`
}

// BeforeCreate - метод, который автоматически генерирует уникальный UUID и
// устанавливает временную метку перед сохранением транзакции
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.UUID = uuid.New().String()
	t.Timestamp = time.Now()
	return
}
