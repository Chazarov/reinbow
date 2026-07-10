// Package database отвечает за подключение к базе данных PostgreSQL,
// настройку пула соединений и миграции.
//
// Подключение конфигурируется через переменные окружения (DB_HOST, DB_PORT,
// DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE). При невозможности установить
// соединение приложение завершается с фатальной ошибкой.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"traineesheep/feedservice/internal/utils"

	_ "github.com/lib/pq"
)

// New создаёт и возвращает указатель на sql.DB с настроенным подключением
// к PostgreSQL. Параметры подключения берутся из переменных окружения:
//   - DB_HOST (по умолчанию "localhost")
//   - DB_PORT (по умолчанию "5432")
//   - DB_USER (по умолчанию "postgres")
//   - DB_PASSWORD (по умолчанию "123")
//   - DB_NAME (по умолчанию "postgres")
//   - DB_SSLMODE (по умолчанию "disable")
//
// Настраивает пул соединений: максимум 25 открытых соединений, 5 бездействующих,
// время жизни соединения — 5 минут. При ошибке открытия или ping соединения
// программа завершается с log.Fatalf.
func New() *sql.DB {
	var db *sql.DB

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		utils.GetEnv("DB_HOST", "localhost"),
		utils.GetEnv("DB_PORT", "5432"),
		utils.GetEnv("DB_USER", "postgres"),
		utils.GetEnv("DB_PASSWORD", "123"),
		utils.GetEnv("DB_NAME", "postgres"),
		utils.GetEnv("DB_SSLMODE", "disable"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не удалось открыть БД: %v", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("Нет соединения с БД: %v", err)
	}

	log.Println("Подключение к базе данных установлено")

	return db
}
