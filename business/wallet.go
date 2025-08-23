package business

import (
	"errors"

	"payment_system_api/database"

	"gorm.io/gorm"
)

// SendMoney выолняет транзакцию перевода средств
func SendMoney(fromAddress, toAddress string, amount float64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var fromWallet, toWallet database.Wallet

		// Блокировка кошелька для предотвращения гонок данных
		if err := tx.Where("address = ?", fromAddress).First(&fromWallet).Error; err != nil {
			return errors.New("кошелек отправителя не найден")
		}
		if err := tx.Where("address = ?", toAddress).First(&toWallet).Error; err != nil {
			return errors.New("кошелек получаетеля не найден")
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

		//Запись транзакции
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

// GetWalletBalance возвращает баланс кошелька по адресу
func GetWalletBalance(address string) (float64, error) {
	var wallet database.Wallet
	if err := database.DB.Where("address = ?", address).First(&wallet).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
