package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUnsubscribe(ctx context.Context, userID int64, chatID int64, query *tgbotapi.CallbackQuery) {
	if err := b.userRepo.Unsubscribe(ctx, userID); err != nil {
		log.Printf("Unsubscribe error: %v", err)
		b.sendReply(chatID, "❌ Ошибка отписки.")
		return
	}
	b.sendReply(chatID, "🔕 Вы отписались от напоминаний.")
	b.api.Send(tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID))
}
