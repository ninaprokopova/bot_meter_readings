package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Callback = нажатие на кнопки
func (b *Bot) handleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) {
	userID := query.From.ID
	chatID := query.Message.Chat.ID

	switch query.Data {
	case "submitted":
		b.handleSubmitted(ctx, userID, chatID, query)
	case "unsubscribe":
		b.handleUnsubscribe(ctx, userID, chatID, query)
	case "generate_readings":
		b.startMeterReadingFlow(chatID, userID)
	}
}
