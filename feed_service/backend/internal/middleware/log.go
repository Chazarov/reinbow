package middleware

import (
	"traineesheep/feedservice/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// RequestIDMiddleware внедряет в каждый запрос уникальный trace_id
// и кладёт в контекст логгер, уже содержащий этот идентификатор.
func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		if c.Path() == "/health" {
			return c.Next()
		}
		// Попробуем взять trace_id из заголовка X-Request-ID
		traceID := c.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Создаём дочерний логгер, автоматически добавляющий поле trace_id
		logger := log.With().Str("trace_id", traceID).Logger()

		// Кладём логгер в контекст запроса
		c.Locals(utils.LoggerKey, &logger)

		// Устанавливаем заголовок, чтобы пробросить trace_id клиенту/в другие сервисы
		c.Set("X-Request-ID", traceID)

		// Можно сразу залогировать факт начала обработки запроса
		logger.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("Входящий запрос")

		// Выполняем следующий обработчик
		return c.Next()
	}
}
