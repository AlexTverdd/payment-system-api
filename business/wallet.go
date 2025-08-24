// Package business содержит бизнес-логику платёжной системы.
package business

import (
	"errors"

	"payment_system_api/database"

	"gorm.io/gorm"
)

// SendMoney выполняет транзакцию перевода средств с одного кошелька на другой.
// Проверяет наличие кошельков, достаточность средств и корректность суммы.
// Все операции выполняются в одной транзакции GORM.
// Возможные ошибки:
// - "кошелек отправителя не найден"
// - "кошелек получаетеля не найден"
// - "недостаточно средств"
// - "сумма должна быть положительной"
// - "нельзя отправлять деньги на тот же адрес"
func SendMoney(fromAddress, toAddress string, amount float64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var fromWallet, toWallet database.Wallet

		if err := tx.Where("address = ?", fromAddress).First(&fromWallet).Error; err != nil {
			return errors.New("кошелек отправителя не найден")
		}
		if err := tx.Where("address = ?", toAddress).First(&toWallet).Error; err != nil {
			return errors.New("кошелек получателя не найден")
		}
		// Проверка баланса
		if fromWallet.Balance < amount {
			return errors.New("недостаточно средств")
		}
		if amount <= 0 {
			return errors.New("сумма должна быть положительной")
		}
		if fromAddress == toAddress {
			return errors.New("нельзя отправлять деньги на тот же адрес")
		}
		// Обновление баланса
		fromWallet.Balance -= amount
		toWallet.Balance += amount

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
			Amount:      amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetWalletBalance возвращает текущий баланс кошелька по адресу.
// Возвращает ошибку, если кошелек не найден.
func GetWalletBalance(address string) (float64, error) {
	var wallet database.Wallet
	if err := database.DB.Where("address = ?", address).First(&wallet).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// GetLastTransactions возвращает последние N транзакций,
// отсортированные по времени создания в порядке убывания.
func GetLastTransactions(count int) ([]database.Transaction, error) {
	var transactions []database.Transaction
	if err := database.DB.Order("timestamp desc").Limit(count).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
