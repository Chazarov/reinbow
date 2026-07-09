// Package repository содержит реализации слоя доступа к данным (DAO)
// для работы с постами, пользователями и refresh-токенами.
//
// DAO-объекты инкапсулируют SQL-запросы и преобразуют результаты
// в модели из пакета models.
package repository

import (
	"database/sql"
	"traineesheep/feedservice/internal/model"

	"github.com/lib/pq"
)

// FeedDAO обеспечивает доступ к данным постов и связям с изображениями.
type FeedDAO struct {
	db *sql.DB
}

// NewFeedDAO создаёт новый экземпляр FeedDAO с заданным подключением к БД.
func NewFeedDAO(db *sql.DB) *FeedDAO {
	return &FeedDAO{db: db}
}

// createPosts — вспомогательная функция, которая итерирует rows, сканирует
// поля поста и массив идентификаторов изображений (через pq.Array),
// после чего собирает срез models.Post.
func createPosts(rows *sql.Rows) []models.Post {
	var posts []models.Post
	var imageIDs []int64
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title,
			&p.Description, &p.CreatedAt, pq.Array(&imageIDs)); err != nil {
			continue
		}

		p.Images = make([]int, len(imageIDs))
		for i, id := range imageIDs {
			p.Images[i] = int(id)
		}
		posts = append(posts, p)
	}

	return posts
}

// LoadFeed загружает все посты из базы данных, отсортированные по дате
// создания (сначала новые). Каждый пост включает массив ID изображений.
// Возвращает срез постов и ошибку (в случае неудачи — пустой срез).
func (fd *FeedDAO) LoadFeed() ([]models.Post, error) {
	rows, postsError := fd.db.Query(`
        SELECT 
            p.id, 
            p.user_id, 
            COALESCE(u.username, '') as username,
            p.title, 
            p.description, 
            TO_CHAR(p.created_at, 'DD.MM.YYYY HH24:MI:SS') as created_at,
            ARRAY_AGG(i.image_id) FILTER (WHERE i.image_id IS NOT NULL) as image_ids
        FROM post p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN image_post i ON p.id = i.post_id
        GROUP BY p.id, u.username
        ORDER BY p.created_at DESC
    `)
	if postsError != nil {
		return make([]models.Post, 0), postsError
	}
	defer rows.Close()

	var posts []models.Post = createPosts(rows)

	return posts, nil
}

// LoadUserFeed загружает посты конкретного пользователя по его userID,
// отсортированные по дате создания (сначала новые). Каждый пост включает
// массив ID изображений. Возвращает срез постов и ошибку.
func (fd *FeedDAO) LoadUserFeed(userID int) ([]models.Post, error) {
	rows, postsError := fd.db.Query(`
        SELECT 
            p.id, 
            p.user_id, 
            COALESCE(u.username, '') as username,
            p.title, 
            p.description, 
			TO_CHAR(p.created_at, 'DD.MM.YYYY HH24:MI:SS') as created_at,
            ARRAY_AGG(i.image_id) FILTER (WHERE i.image_id IS NOT NULL) as image_ids
        FROM post p
        LEFT JOIN users u ON p.user_id = u.id
        LEFT JOIN image_post i ON p.id = i.post_id
        WHERE p.user_id = $1
        GROUP BY p.id, u.username
        ORDER BY p.created_at DESC
    `, userID)
	if postsError != nil {
		return make([]models.Post, 0), postsError
	}
	defer rows.Close()

	var posts []models.Post = createPosts(rows)

	return posts, nil
}

// CreatePost создаёт новый пост с указанными заголовком, описанием и
// списком ID изображений. Изображения должны быть предварительно сохранены
// во внешнем сервисе. Метод добавляет запись в таблицу post и создаёт
// соответствующие связи в image_post. Возвращает созданный пост и ошибку.
func (fd *FeedDAO) CreatePost(userID int, title string, description string, imageIDs []int) (models.Post, error) {
	if len(imageIDs) == 0 {
		// todo заглушка
		//return models.Post{}, fmt.Errorf("Изображения не выбраны")
	}

	var post models.Post
	postError := fd.db.QueryRow(
		"INSERT INTO post (user_id, title, description) VALUES ($1, $2, $3) RETURNING id, user_id, title, description, created_at",
		userID, title, description,
	).Scan(&post.ID, &post.UserID, &post.Title, &post.Description, &post.CreatedAt)

	if postError != nil {
		return models.Post{}, postError
	}

	for _, imgID := range imageIDs {
		_, imagePostError := fd.db.Exec(
			"INSERT INTO image_post (post_id, image_id) VALUES ($1, $2)",
			post.ID, imgID,
		)
		if imagePostError != nil {
			return models.Post{}, imagePostError
		}
	}

	return post, nil
}
