package database

import (
	"database/sql"
	"log"
)

// Migrate выполняет создание необходимых таблиц в базе данных, если они
// ещё не существуют. При ошибке создания любой из таблиц программа
// завершается с фатальной ошибкой.
//
// Создаваемые таблицы:
//   - users (id, username, password, email, created_at)
//   - post (id, user_id, title, description, created_at)
//   - image_post (id, post_id, image_id)
//
// Таблица refresh_tokens закомментирована, так как в текущей итерации
// хранение токенов реализовано в памяти.
func Migrate(db *sql.DB) {
	log.Println("Создаем таблицы")

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
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

	// todo в текущей итерации приложения это лишнее
	/*createRefreshTokenTable := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(512) NOT NULL UNIQUE,
		expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
		created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
	);`*/

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}
	log.Println("Таблица users готова")

	if _, err := db.Exec(createPostTable); err != nil {
		log.Fatalf("Ошибка создания таблицы post: %v", err)
	}
	log.Println("Таблица post готова")

	if _, err := db.Exec(createImagePostTable); err != nil {
		log.Fatalf("Ошибка создания таблицы image_post: %v", err)
	}
	log.Println("Таблица image_post готова")

	/*if _, err := db.Exec(createRefreshTokenTable); err != nil {
		log.Fatalf("Ошибка создания таблицы refresh_tokens: %v", err)
	}
	log.Println("Таблица refresh_tokens готова")*/
}
