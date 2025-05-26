package bot

import (
	"log"
	"submit_meter_readings/config"
	"submit_meter_readings/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageSender interface {
	SendMessage(msg tgbotapi.MessageConfig)
}

type TelegramBotSender struct {
	bot *tgbotapi.BotAPI
}

func (t *TelegramBotSender) SendMessage(msg tgbotapi.MessageConfig) {
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err.Error())
	} else {
		log.Println("Сообщение отправлено пользователю", msg.ChatID, msg.Text)
	}
}

// Вопрос можно ли от поля api избавиться? сейчас олтправка есьб через sender
type Bot struct {
	api        *tgbotapi.BotAPI
	sender     MessageSender
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
		sender:     &TelegramBotSender{bot: api},
		userRepo:   repo,
		userStates: make(map[int64]*UserState),
	}, nil
}

func (b *Bot) Start() {
	go b.handleUpdates()
	go b.startReminder()
}
