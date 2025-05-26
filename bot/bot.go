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

type MessageDeleter interface {
	DeleteMessage(msg tgbotapi.DeleteMessageConfig)
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

func (t *TelegramBotSender) DeleteMessage(msg tgbotapi.DeleteMessageConfig) {
	_, err := t.bot.Request(msg)
	if err != nil {
		log.Println("Ошибка при удалении сообщения:", err.Error())
	} else {
		log.Println("Сообщение удалено у пользователя", msg.ChatID)
	}
}

// Вопрос можно ли от поля api избавиться? сейчас олтправка есьб через sender
type Bot struct {
	api        *tgbotapi.BotAPI
	sender     MessageSender
	deleter    MessageDeleter
	userRepo   storage.UserRepository
	userStates map[int64]*UserState
}

func NewBot(cfg *config.Config, repo storage.UserRepository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	// Вот тут мне не нравится что два одинаковых экземпляра TelegramBotSender создается, не нравится
	return &Bot{
		api:        api,
		sender:     &TelegramBotSender{bot: api},
		deleter:    &TelegramBotSender{bot: api},
		userRepo:   repo,
		userStates: make(map[int64]*UserState),
	}, nil
}

func (b *Bot) Start() {
	go b.handleUpdates()
	go b.startReminder()
}
