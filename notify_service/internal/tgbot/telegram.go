package tgbot

// Пакет tgbot реализует функционал отправки уведомлений в Телеграм с помощью бота

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Функция HandleSendMessage используется для отправки сообщений в Telegram бота
func SendMessage(b *bot.Bot, Ctx context.Context, chat_id int64, content string) error {
	if _, err := b.SendMessage(Ctx, &bot.SendMessageParams{
		ChatID:    chat_id,              // ID чата куда надо отправить сообщение
		Text:      content,              // Текст сообщения
		ParseMode: models.ParseModeHTML, // Параметр для использования HTML при форматировании текста
	}); err != nil {
		return err
	}
	return nil
}

func Handler(Ctx context.Context, b *bot.Bot, update *models.Update) {

}
