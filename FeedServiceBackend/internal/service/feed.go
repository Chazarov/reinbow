// Package service содержит бизнес-логику приложения, реализованную
// в виде сервисов, которые используют соответствующие DAO-объекты.
//
// Сервисы инкапсулируют правила обработки данных и предоставляют
// удобный интерфейс для контроллеров.
package service

import (
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/repository"
)

// FeedService предоставляет методы для работы с лентой постов:
// получение общей ленты, ленты конкретного пользователя и создание новых постов.
type FeedService struct {
	FeedDAO *repository.FeedDAO
}

// NewFeedService создаёт новый экземпляр FeedService с переданным DAO.
func NewFeedService(feedDAO *repository.FeedDAO) *FeedService {
	return &FeedService{FeedDAO: feedDAO}
}

// LoadFeed загружает все посты (общая лента), отсортированные по дате создания.
// Возвращает срез постов и ошибку.
func (fs *FeedService) LoadFeed() ([]models.Post, error) {
	return fs.FeedDAO.LoadFeed()
}

// LoadUserFeed загружает посты конкретного пользователя по его идентификатору.
// Возвращает срез постов и ошибку.
func (fs *FeedService) LoadUserFeed(userID int) ([]models.Post, error) {
	return fs.FeedDAO.LoadUserFeed(userID)
}

// CreatePost создаёт новый пост с указанным автором, заголовком, описанием
// и списком идентификаторов изображений. Изображения должны быть предварительно
// сохранены во внешнем сервисе. Возвращает созданный пост и ошибку.
func (fs *FeedService) CreatePost(userID int, title string, description string, imageIDs []int) (models.Post, error) {
	return fs.FeedDAO.CreatePost(userID, title, description, imageIDs)
}
