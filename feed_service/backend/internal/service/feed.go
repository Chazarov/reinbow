// Package service содержит бизнес-логику приложения, реализованную
// в виде сервисов, которые используют соответствующие DAO-объекты.
//
// Сервисы инкапсулируют правила обработки данных и предоставляют
// удобный интерфейс для контроллеров.
package service

import (
	"mime/multipart"
	"traineesheep/feedservice/internal/client/image_service"
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/repository"
)

// FeedService предоставляет методы для работы с лентой постов:
// получение общей ленты, ленты конкретного пользователя и создание новых постов.
type FeedService struct {
	FeedDAO     *repository.FeedDAO
	ImageClient *client.ImageClient
}

// NewFeedService создаёт новый экземпляр FeedService с переданным DAO.
func NewFeedService(feedDAO *repository.FeedDAO,
	imageClient *client.ImageClient) *FeedService {
	return &FeedService{FeedDAO: feedDAO, ImageClient: imageClient}
}

// LoadFeed загружает все посты (общая лента), отсортированные по дате создания.
// Возвращает срез постов и ошибку.
func (feedService *FeedService) LoadFeed() ([]models.Post, error) {
	return feedService.FeedDAO.LoadFeed()
}

// LoadUserFeed загружает посты конкретного пользователя по его идентификатору.
// Возвращает срез постов и ошибку.
func (feedService *FeedService) LoadUserFeed(userID int) ([]models.Post, error) {
	return feedService.FeedDAO.LoadUserFeed(userID)
}

// CreatePost создаёт новый пост с указанным автором, заголовком, описанием
// и списком идентификаторов изображений. Изображения должны быть предварительно
// сохранены во внешнем сервисе. Возвращает созданный пост и ошибку.
func (feedService *FeedService) CreatePost(userID int, title string, description string, files []*multipart.FileHeader) (models.Post, error) {
	imageIDs, err := feedService.ImageClient.SaveFiles(files)
	if err != nil {
		// todo заглушка на случай, если сервис хранения не развернут
		// todo просто создаём пост без изображений
		// return models.Post{}, err
	}
	return feedService.FeedDAO.CreatePost(userID, title, description, imageIDs)
}
