package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка сообщений (показаний) и сообщений-команд, начинающихся с /
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	if !msg.IsCommand() {
		state, ok := b.userStates[userID]
		if ok {
			b.handleMeterReadingInput(ctx, msg, state)
		}
		return
	}

	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, chatID, userID)
	case "status":
		b.handleStatusCommand(ctx, chatID, userID)
	}
}

func (b *Bot) handleStartCommand(ctx context.Context, chatID, userID int64) {
	err := b.userRepo.Subscribe(ctx, userID)
	if err != nil {
		log.Printf("Subscribe error: %v", err)
		msgError := tgbotapi.NewMessage(chatID, "❌ Ошибка подписки. Попробуйте позже.")
		b.sender.SendMessage(msgError)
		return
	}

	reply := tgbotapi.NewMessage(chatID, "✅ Вы подписались на напоминания!\nОни будут приходить c 12:00 по 15:00 с 20 по 25 число каждого месяца.")
	b.sender.SendMessage(reply)
}

func (b *Bot) handleStatusCommand(ctx context.Context, chatID, userID int64) {

	shouldNotify, err := b.userRepo.ShouldNotify(ctx, userID)
	if err != nil {
		log.Printf("Status check error: %v", err)
		msgError := tgbotapi.NewMessage(chatID, "❌ Ошибка проверки статуса.")
		b.sender.SendMessage(msgError)
		return
	}

	statusText := "🔔 Вы подписаны на напоминания"
	if !shouldNotify {
		statusText = "✅ Вы уже передали показания в этом месяце"
	}
	msgStatus := tgbotapi.NewMessage(chatID, statusText)
	b.sender.SendMessage(msgStatus)
}
