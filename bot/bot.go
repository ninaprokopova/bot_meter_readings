package bot

import (
	"submit_meter_readings/config"
	"submit_meter_readings/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	userRepo   storage.UserRepository
	userStates map[int64]*UserState
}

func NewBot(cfg *config.Config, repo storage.UserRepository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:        api,
		userRepo:   repo,
		userStates: make(map[int64]*UserState),
	}, nil
}

func (b *Bot) Start() {
	go b.handleUpdates()
	go b.startReminder()
}
