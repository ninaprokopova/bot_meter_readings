package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdates() {
	updates := b.api.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		switch {
		case update.Message != nil:
			b.handleMessage(ctx, update.Message)
		case update.CallbackQuery != nil:
			b.handleCallback(ctx, update.CallbackQuery)
		}

		cancel()
	}
}
