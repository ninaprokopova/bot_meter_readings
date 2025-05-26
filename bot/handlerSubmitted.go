package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSubmitted(ctx context.Context, userID int64, chatID int64, query *tgbotapi.CallbackQuery) {
	err := b.userRepo.MarkAsSubmitted(ctx, userID)
	if err != nil {
		log.Printf("Mark as submitted error: %v", err)
		msgSaveError := tgbotapi.NewMessage(chatID, "❌ Ошибка сохранения. Попробуйте позже.")
		b.sender.SendMessage(msgSaveError)
		return
	}
	msgSubmit := tgbotapi.NewMessage(chatID, "✅ Вы передали показания. До следующего месяца")
	b.sender.SendMessage(msgSubmit)
	msgForDelete := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
	b.deleter.DeleteMessage(msgForDelete)
}
