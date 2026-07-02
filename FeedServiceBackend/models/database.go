package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"traineesheep/feedservice/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB


func InitDB() {
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
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не удалось открыть БД: %v", err)
	}

	// Настройка пула соединений
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatalf("Нет соединения с БД: %v", err)
	}

	log.Println("Подключение к базе данных установлено")

	initTables()
}

func closeDB() {
	if DB != nil {
		DB.Close()
		log.Println("Соединение с базой данных закрыто")
	}
}

func initTables() {
	log.Println("Создаем таблицы")

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			tg_chat_id VARCHAR(255),
			email VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
	);`

	createPostTable := `
		CREATE TABLE IF NOT EXISTS post (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255),
			description TEXT,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
	);`

	createImagePostTable := `
		CREATE TABLE IF NOT EXISTS image_post (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL REFERENCES post(id) ON DELETE CASCADE,
			image_id INT NOT NULL
	);`

	createRefreshTokenTable := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(512) NOT NULL UNIQUE,
		expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
		created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
	);`


	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}
	log.Println("Таблица users готова")

	if _, err := DB.Exec(createPostTable); err != nil {
		log.Fatalf("Ошибка создания таблицы post: %v", err)
	}
	log.Println("Таблица post готова")

	if _, err := DB.Exec(createImagePostTable); err != nil {
		log.Fatalf("Ошибка создания таблицы image_post: %v", err)
	}
	log.Println("Таблица image_post готова")

	if _, err := DB.Exec(createRefreshTokenTable); err != nil {
		log.Fatalf("Ошибка создания таблицы refresh_tokens: %v", err)
	}
	log.Println("Таблица refresh_tokens готова")
}
