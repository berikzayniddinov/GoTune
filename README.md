# 🎵 GoTune Backend

## 📋 Project Overview

GoTune — это backend-сервис для магазина музыкальных инструментов, написанный на Go. Система реализует REST API, авторизацию по JWT, хранение данных в NoSQL базе и легко масштабируется под будущие микросервисы.

Проект следует принципам Clean Architecture, Domain-Driven Design (DDD) и отделяет слои `entity`, `service`, `handler`, `repository`.



## 🚀 Tech Stack

| Технология     | Описание                               |
|----------------|----------------------------------------|
| Language       | Go (Golang)                            |
| Framework      | Chi Router / Fiber (или другой по выбору) |
| Database       | NoSQL (например, MongoDB или DynamoDB) |
| Auth           | JWT (JSON Web Tokens)                  |
| Dependency Mgmt| Go Modules                             |
| Architecture   | Clean Architecture / DDD               |

---

## 🗂️ Project Structure

| Директория     | Назначение                                                           |
|----------------|----------------------------------------------------------------------|
| `/cmd`         | Точка входа в приложение (`main.go`)                                 |
| `/internal/entity`     | Бизнес-сущности: `Instrument`, `User`, `Order`, и т.д.                  |
| `/internal/service`    | Прикладная логика: `CartService`, `AuthService` и другие                |
| `/internal/handler`    | HTTP-обработчики: реализация REST API                         |
| `/internal/repository` | Работа с NoSQL БД (MongoDB, DynamoDB): интерфейсы и реализации |
| `/internal/config`     | Конфигурация приложения: переменные среды, инициализация и secrets |
| `/internal/middleware` | Middleware: JWT, логгирование, авторизация и др.              |
| `/scripts`     | Миграции, генерация тестовых данных, сидеры и служебные скрипты     |

---

## 🔐 Authentication

- Используется **JWT токенизация** с доступом по `Authorization: Bearer <token>`
- Поддержка middleware для защиты приватных маршрутов
- Возможность добавления refresh-token механизма

---

## 🧪 Testing

- Юнит-тесты покрывают бизнес-логику
- Используются `httptest` и `testcontainers` (при необходимости)
- Можно запускать с `go test ./...`



## 💾 Database: NoSQL

- Поддержка MongoDB / DynamoDB (в зависимости от конфигурации)
- Структуры маппятся вручную (без ORM) для максимальной гибкости
- Возможность миграции через `mongo-migrate` или встроенные `scripts/`



## 🧠 Branch Naming Conventions

| Prefix      | Purpose                                     |
|-------------|---------------------------------------------|
| `feature/`  | Новая функциональность                      |
| `bugfix/`   | Исправление багов                           |
| `hotfix/`   | Критические фиксы в проде                   |
| `refactor/` | Рефакторинг без изменения поведения         |
| `test/`     | Добавление или изменение тестов             |
| `chore/`    | Обновление зависимостей, документации и т.п.|
| `release/`  | Подготовка и сборка релиза                  |

**Примеры**:
- `feature/add-cart-service`
- `bugfix/fix-auth-middleware`
- `refactor/split-user-handler`



## 🛠️ Build & Run

```bash
go mod tidy
go run ./cmd
```

---

## ✅ TODO

- [x] Entity слои: Instrument, User, Cart, Order
- [x] JWT Middleware
- [ ] Модуль оплаты
- [ ] Интеграция с внешним складом
- [ ] Swagger API документация

---

## 🤝 Contributing

Pull Requests приветствуются! Пожалуйста, соблюдайте соглашения по неймингам и архитектуре. Запускайте тесты перед коммитом.

---

## 📄 License

MIT License — используй, расширяй, вдохновляйся 🎧


---

Хочешь — сгенерирую тебе Swagger JSON + HTML доку под этот проект или docker-compose с Mongo + seed-данными.

```