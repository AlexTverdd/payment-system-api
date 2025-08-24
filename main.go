// Package main запускает API сервера платёжной системы.
// Он настраивает маршруты, подключается к базе данных, выполняет миграции
// и запускает HTTP-сервер на Gin.
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"payment_system_api/config"
	"payment_system_api/database"
	"payment_system_api/handlers"
)

// main инициализирует сервер приложения.
// Она выполняет следующие шаги:
// 1. Загружает конфигурацию.
// 2. Подключается к базе данных.
// 3. Выполняет миграции и начальную настройку данных.
// 4. Настраивает маршруты API.
// 5. Запускает HTTP-сервер на порту 8080.
func main() {

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	database.ConnectDB(cfg.DatabaseURL)

	// Выполнение миграции и создание начальных данных
	database.Migrate()
	database.InitialSetup()

	//Настройка Gin
	router := gin.Default()

	//Группировка маршрутов
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/send", handlers.SendHandler)
		apiRoutes.GET("/wallet/:address/balance", handlers.GetBalanceHandler)
		apiRoutes.GET("/transactions", handlers.GetLastTransactionsHandler)
	}
	log.Println("Старт сервера на порту 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
