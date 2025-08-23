package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"payment_system_api/business"
)

// SendRequest - структура для тела запроса POST /api/send
type SendRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// SendHandler обрабатывает POST /api/send

func SendHandler(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ошибка": "Неверное тело запроса", "подробно": err.Error()})
		return
	}

	err := business.SendMoney(req.From, req.To, req.Amount)
	if err != nil {
		if err.Error() == "недостаточно средств" {
			c.JSON(http.StatusPaymentRequired, gin.H{"ошибка": "Недостаточно средств"})
			return
		}
		if err.Error() == "кошелек отправителя не найден" || err.Error() == "кошелек получаетеля не найден" {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "Транзакция неуспешна", "детали": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"сообщение": "Транзакция успешна"})
}
