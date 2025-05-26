package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUnsubscribe(ctx context.Context, userID int64, chatID int64, query *tgbotapi.CallbackQuery) {
	if err := b.userRepo.Unsubscribe(ctx, userID); err != nil {
		log.Printf("Ошибка при попытке отписаться от бота: %v", err)
		msgSaveError := tgbotapi.NewMessage(chatID, "❌ Ошибка отписки. Попробуйте позже.")
		b.sender.SendMessage(msgSaveError)
		return
	}
	msgSubmit := tgbotapi.NewMessage(chatID, "🔕 Вы отписались от напоминаний.")
	msgForDelete := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
	b.sender.SendMessage(msgSubmit)
	b.deleter.DeleteMessage(msgForDelete)
}
