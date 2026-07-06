package tgbot

// Пакет tgbot реализует функционал отправки уведомлений в Телеграм с помощью бота

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Функция HandleSendMessage используется для отправки сообщений в Telegram бота
func HandleSendMessage(b *bot.Bot, ctx context.Context, chat_id int64, content string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chat_id,              // ID чата куда надо отправить сообщение
		Text:      content,              // Текст сообщения
		ParseMode: models.ParseModeHTML, // Параметр для использования HTML при форматировании текста
	}); err != nil {
		log.Println("Ошибка при отправке сообщения в Telegram бота:", err.Error())
	}
}

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {

}
