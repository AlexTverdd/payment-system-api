# payment-system-api

Test task for Infotecs — REST API на Go для платёжной системы.

## Содержание

- [Описание](#описание)  
- [Технологии](#технологии)  
- [Установка и запуск](#установка-и-запуск)  
- [API Endpoints](#api-endpoints)  
- [Примеры использования](#примеры-использования)  
- [Тестирование](#тестирование)  
- [Контакты](#контакты)

## Описание

Сервис на Go, реализующий:
- перевод средств между кошельками (`POST /api/send`),
- получение последних транзакций (`GET /api/transactions?count=N`),
- получение баланса конкретного кошелька (`GET /api/wallet/{address}/balance`).

Используется **Gin** для роутинга для базового сервера, **GORM** для работы с БД и **PostgreSQL** как хранилище данных.

## Технологии

- [Go](https://go.dev/)  
- [Gin](https://github.com/gin-gonic/gin) — HTTP веб-фреймворк  
- [net/http](https://pkg.go.dev/net/http) — стандартная библиотека  
- [GORM](https://gorm.io/) — ORM для работы с базой данных  
- [PostgreSQL](https://www.postgresql.org/)  
- [Docker](https://www.docker.com/) + `docker-compose`

## Установка и запуск

### Локально

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/AlexTverdd/payment-system-api.git
   cd payment-system-api
   
2. Установить зависимости:
    ```bash
    go mod tidy

3. Создать базу данных

4. Настроить переменную окружения .env пример

DATABASE_URL=postgres://postgres:postgres@localhost:5432/payment_system?sslmode=disable

5. Запустить сервис 
    ```bash
    go run main.go

6. Сервис будет доступен по адресу http://localhost:8080.

### Через Docker

1. Собрать и запустить контейнеры:
    ```bash
    docker-compose up --build

2. После запуска API будет доступен по адресу http://localhost:8080.
