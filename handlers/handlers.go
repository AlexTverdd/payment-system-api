// Package handlers содержит HTTP-обработчики для API платёжной системы.
// Он реализует конечные точки на базе фреймворка Gin.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment_system_api/business"
)

// SendRequest представляет тело запроса для POST /api/send.
type SendRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount"`
}

// SendHandler обрабатывает POST /api/send.
//
// Принимает JSON с From, To и Amount и выполняет транзакцию через бизнес-логику.
// Возвращает:
// - 200 OK при успешной транзакции
// - 402 Payment Required, если недостаточно средств
// - 404 Not Found, если кошелек не найден
// - 400 Bad Request, если тело запроса неверное
// - 500 Internal Server Error при других ошибках
func SendHandler(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ошибка": "Неверное тело запроса", "подробно": err.Error()})
		return
	}

	err := business.SendMoney(req.From, req.To, req.Amount)
	if err != nil {
		// таблица соответствий бизнес-ошибок к HTTP-кодам
		errorMap := map[string]int{
			"недостаточно средств":                     http.StatusPaymentRequired,
			"кошелек отправителя не найден":            http.StatusNotFound,
			"кошелек получаетеля не найден":            http.StatusNotFound,
			"сумма должна быть положительной":          http.StatusBadRequest,
			"нельзя отправлять деньги на тот же адрес": http.StatusBadRequest,
		}

		if code, ok := errorMap[err.Error()]; ok {
			c.JSON(code, gin.H{"ошибка": err.Error()})
			return
		}

		// всё остальное — внутренняя ошибка
		c.JSON(http.StatusInternalServerError, gin.H{
			"ошибка": "Транзакция неуспешна", "детали": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"сообщение": "Транзакция успешна"})
}

// GetBalanceHandler обрабатывает GET /api/wallet/{address}/balance.
//
// Возвращает баланс указанного кошелька.
// Если кошелек не найден — 404 Not Found.
// При внутренних ошибках — 500 Internal Server Error.
func GetBalanceHandler(c *gin.Context) {
	address := c.Param("address")

	balance, err := business.GetWalletBalance(address)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "Кошелек не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "Не удалось получить баланс", "детали": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"адрес": address, "баланс": balance})
}

// GetLastTransactionsHandler обрабатывает GET /api/transactions?count=N.
//
// Возвращает последние N транзакций.
// Если параметр count некорректный — 400 Bad Request.
// При внутренних ошибках — 500 Internal Server Error.
func GetLastTransactionsHandler(c *gin.Context) {
	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "Неверный 'count' параметр"})
		return
	}

	transactions, err := business.GetLastTransactions(count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "Не удалось получить транзакции"})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
