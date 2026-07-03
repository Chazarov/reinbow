package controller

import (
	"log"
	"strconv"
	models "traineesheep/feedservice/internal/model"
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) Feed(c *fiber.Ctx) error {
	userID := c.Query("user_id")

	var posts []models.Post
	var postsError error

	if userID == "" {
		posts, postsError = ctrl.FeedService.LoadFeed()
	} else {
		userIDInt, userIDError := strconv.Atoi(userID)
		if userIDError != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
				Data:       nil,
				Success:    false,
				ErrMessage: "Некорректный user_id",
			})
		}
		posts, postsError = ctrl.FeedService.LoadUserFeed(userIDInt)
	}

	if postsError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Data:       posts,
			Success:    false,
			ErrMessage: "Ошибка получения постов",
		})
	}

	if userID == "" {
		log.Printf("GET /feed: Лента загружена, постов: %d", len(posts))
	} else {
		log.Printf("GET /feed/%s: Лента пользователя загружена, постов: %d", userID, len(posts))
	}
	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Data:       posts,
		Success:    true,
		ErrMessage: "",
	})
}
