package controller

import (
	"fmt"
	"strconv"
	"traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

// Feed обрабатывает GET-запросы на получение ленты постов.
//
// Если параметр запроса user_id отсутствует, возвращает общую ленту,
// отсортированную по дате создания (сначала новые).
// Если передан user_id, загружает ленту конкретного пользователя.
// Некорректный user_id приводит к ответу 400 Bad Request.
//
// Ответ:
//   - 200: { success: true, data: []Post }
//   - 400: { success: false, err_message: "Некорректный user_id" }
//   - 500: { success: false, err_message: "Ошибка получения постов" }
func (ctrl *Controller) Feed(c fiber.Ctx) error {
	logger := c.Locals(utils.LoggerKey).(*zerolog.Logger)
	userID := c.Query("user_id")

	var posts []models.Post
	var postsError error

	if userID == "" {
		posts, postsError = ctrl.FeedService.LoadFeed()
	} else {
		userIDInt, userIDError := strconv.Atoi(userID)
		if userIDError != nil {
			logger.Error().
				Str("user_id", userID).
				Str("path", c.Path()).
				Msg("Ошибка парсинга user_id: : " + userIDError.Error())
			return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
				Data:       nil,
				Success:    false,
				ErrMessage: "Некорректный user_id",
			})
		}
		posts, postsError = ctrl.FeedService.LoadUserFeed(userIDInt)
	}

	if postsError != nil {
		logger.Error().
			Str("path", c.Path()).
			Msg("Ошибка выборки данных из БД: " + postsError.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       posts,
			Success:    false,
			ErrMessage: "Ошибка получения постов",
		})
	}

	logger.Info().
		Str("user_id", userID).
		Str("path", c.Path()).
		Msg(fmt.Sprintf("Лента загружена, постов %d", len(posts)))

	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Data:       posts,
		Success:    true,
		ErrMessage: "",
	})
}
