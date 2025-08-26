// Package business содержит бизнес-логику платёжной системы.
// Он предоставляет функции для перевода средств между кошельками,
// получения баланса и списка последних транзакций.
package business

import (
	"errors"
	"time"

	"payment_system_api/database"

	"gorm.io/gorm"
)

// TransactionResponse представляет транзакцию,
// возвращаемую в API. В отличие от модели базы данных,
// сумма хранится в виде float64.
type TransactionResponse struct {
	FromAddress string    `json:"from_address"` // адрес отправителя
	ToAddress   string    `json:"to_address"`   // адрес получателя
	Amount      float64   `json:"amount"`       // сумма перевода
	Timestamp   time.Time `json:"timestamp"`    // время создания транзакции
	UUID        string    `json:"uuid"`         // уникальный идентификатор транзакции
}

// SendMoney выполняет транзакцию перевода средств с одного кошелька на другой.
//
// Проверяет наличие кошельков, достаточность средств и корректность суммы.
// Все операции выполняются в одной транзакции GORM.
// Возможные ошибки:
// - "кошелек отправителя не найден"
// - "кошелек получаетеля не найден"
// - "недостаточно средств"
// - "сумма должна быть положительной"
// - "нельзя отправлять деньги на тот же адрес"
func SendMoney(fromAddress, toAddress string, amount float64) error {
	amountAsInteger := int64(amount * 100)
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var fromWallet, toWallet database.Wallet

		if err := tx.Where("address = ?", fromAddress).First(&fromWallet).Error; err != nil {
			return errors.New("кошелек отправителя не найден")
		}

		if err := tx.Where("address = ?", toAddress).First(&toWallet).Error; err != nil {
			return errors.New("кошелек получателя не найден")
		}

		// Проверка баланса
		if fromWallet.Balance < amountAsInteger {
			return errors.New("недостаточно средств")
		}

		if amountAsInteger <= 0 {
			return errors.New("сумма должна быть положительной")
		}

		if fromAddress == toAddress {
			return errors.New("нельзя отправлять деньги на тот же адрес")
		}

		// Обновление баланса
		fromWallet.Balance -= amountAsInteger
		toWallet.Balance += amountAsInteger

		if err := tx.Save(&fromWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&toWallet).Error; err != nil {
			return err
		}

		// Запись транзакции
		transaction := database.Transaction{
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Amount:      amountAsInteger,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetWalletBalance возвращает текущий баланс кошелька по адресу.
//
// Баланс возвращается в виде float64.
// Возвращает ошибку, если кошелек не найден.
func GetWalletBalance(address string) (float64, error) {
	var wallet database.Wallet
	if err := database.DB.Where("address = ?", address).First(&wallet).Error; err != nil {
		return 0, err
	}
	return float64(wallet.Balance) / 100, nil
}

// GetLastTransactions возвращает последние N транзакций,
//
// отсортированные по времени создания в порядке убывания.
func GetLastTransactions(count int) ([]TransactionResponse, error) {
	var transactionsDB []database.Transaction
	if err := database.DB.Order("timestamp desc").Limit(count).Find(&transactionsDB).Error; err != nil {
		return nil, err
	}

	var transactionsAPI []TransactionResponse
	for _, t := range transactionsDB {
		transactionsAPI = append(transactionsAPI, TransactionResponse{
			FromAddress: t.FromAddress,
			ToAddress:   t.ToAddress,
			Amount:      float64(t.Amount) / 100,
			Timestamp:   t.Timestamp,
			UUID:        t.UUID,
		})
	}

	return transactionsAPI, nil
}
