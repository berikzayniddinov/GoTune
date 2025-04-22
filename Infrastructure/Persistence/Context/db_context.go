// Infrastructure/Persistence/Context/db_context.go
package Context

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Константы для подключения к БД
const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "password" // Замените на ваш пароль
	DbName   = "dbhabit"
	Schema   = "gotune"
)

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB() (*sql.DB, error) {
	// Строка подключения к PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d instrument=%s password=%s dbname=%s schema=%s sslmode=disable",
		Host, Port, User, Password, DbName, Schema)

	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")

	return db, nil
}

func InitUserTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			address VARCHAR(255),
			phone VARCHAR(50)
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table initialized")
	return nil
}
