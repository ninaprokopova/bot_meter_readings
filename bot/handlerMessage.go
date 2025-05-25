package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка сообщений (показаний) и сообщений-команд, начинающихся с /
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		if state, ok := b.userStates[msg.From.ID]; ok {
			b.handleMeterReadingInput(ctx, msg, state)
			return
		}
	}

	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, msg)
	case "status":
		b.handleStatusCommand(ctx, msg)
	}
}

func (b *Bot) handleStartCommand(ctx context.Context, msg *tgbotapi.Message) {
	userID := msg.From.ID

	if err := b.userRepo.Subscribe(ctx, userID); err != nil {
		log.Printf("Subscribe error: %v", err)
		b.sendReply(msg.Chat.ID, "❌ Ошибка подписки. Попробуйте позже.")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "✅ Вы подписались на напоминания!\nОни будут приходить в 12:00 с 20 по 25 число каждого месяца. ")

	if _, err := b.api.Send(reply); err != nil {
		log.Printf("Send message error: %v", err)
	}
}

func (b *Bot) handleStatusCommand(ctx context.Context, msg *tgbotapi.Message) {
	userID := msg.From.ID

	shouldNotify, err := b.userRepo.ShouldNotify(ctx, userID)
	if err != nil {
		log.Printf("Status check error: %v", err)
		b.sendReply(msg.Chat.ID, "❌ Ошибка проверки статуса.")
		return
	}

	statusText := "🔔 Вы подписаны на напоминания"
	if !shouldNotify {
		statusText = "✅ Вы уже передали показания в этом месяце"
	}

	b.sendReply(msg.Chat.ID, statusText)
}

func (b *Bot) sendReply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Send reply error: %v", err)
	}
}
