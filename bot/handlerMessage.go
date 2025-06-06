package bot

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка сообщений (показаний) и сообщений-команд, начинающихся с /
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	if !msg.IsCommand() {
		state, ok := b.userStates[userID]
		if ok {
			if state.CurrentStep == "template" {
				b.saveNewTemplate(ctx, msg, userID, chatID)
			} else {
				b.handleMeterReadingInput(ctx, msg, state)
			}
		}
		return
	}

	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, chatID, userID)
	case "status":
		b.handleStatusCommand(ctx, chatID, userID)
	case "template":
		b.handleTemplateCommand(chatID, userID)
	}
}

// Мысль! А что если во время передачи показаний выбрать \start, \status, \template??
// Надо потестить
func (b *Bot) saveNewTemplate(ctx context.Context, msg *tgbotapi.Message, userID int64, chatID int64) {
	newTemplate := msg.Text
	var msgText string
	if strings.Contains(newTemplate, "*показания*") {
		err := b.userRepo.ChangeTemplate(ctx, uint64(userID), newTemplate)
		if err != nil {
			msgText = "Ошибка изменения шаблона, попробуйте позднее"
		} else {
			msgText = "Шаблон изменен :)"
		}
	} else {
		msgText = "Шаблон не изменен: в новом шаблоне нет подстроки *показания* \n Попробуйте еще раз: /template"
	}
	msgToUser := tgbotapi.NewMessage(chatID, msgText)
	b.sender.SendMessage(msgToUser)
	delete(b.userStates, userID)
}

func (b *Bot) handleTemplateCommand(chatID, userID int64) {
	b.userStates[userID] = &UserState{CurrentStep: "template"}
	msgError := tgbotapi.NewMessage(chatID, MessageToChangeTemplate)
	b.sender.SendMessage(msgError)
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
	isSubscribed, hasSubmitted, err := b.userRepo.GetUserStatus(ctx, userID)
	if err != nil {
		log.Printf("Status check error: %v", err)
		msgError := tgbotapi.NewMessage(chatID, "❌ Ошибка проверки статуса.")
		b.sender.SendMessage(msgError)
		return
	}

	var statusText string
	if !isSubscribed {
		statusText = "🔕 Вы не подписаны на напоминания"
	} else {
		if hasSubmitted {
			statusText = "✅ Вы уже передали показания в этом месяце"
		} else {
			statusText = "🔔 Вы подписаны на напоминания"
		}
	}
	msgStatus := tgbotapi.NewMessage(chatID, statusText)
	b.sender.SendMessage(msgStatus)
}
