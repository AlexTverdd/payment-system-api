package database

import (
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
	gorm.Model
	//FromAddress адрес кошелька отправителя
	FromAddress string
	//ToAddress адрес кошелька получателя
	ToAddress string
	//Amount сумма перевода
	Amount float64
	//UUID уникальный идентификатор транзакции
	UUID string `gorm:"unique;not null"`
}
