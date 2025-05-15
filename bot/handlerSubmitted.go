package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSubmitted(ctx context.Context, userID int64, chatID int64, query *tgbotapi.CallbackQuery) {
	if err := b.userRepo.MarkAsSubmitted(ctx, userID); err != nil {
		log.Printf("Mark as submitted error: %v", err)
		b.sendReply(chatID, "❌ Ошибка сохранения. Попробуйте позже.")
		return
	}
	b.sendReply(chatID, "✅ Спасибо! Показания учтены.")
	b.api.Send(tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID))
}
