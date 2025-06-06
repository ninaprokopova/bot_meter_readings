package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) startReminder() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastResetMonth := time.Now().Month()

	var wasRemindToday = false

	for now := range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		lastResetMonth = b.resetStatus(ctx, now, lastResetMonth)
		wasRemindToday = b.sendReminds(ctx, now, wasRemindToday)

		cancel()
	}
}

func (b *Bot) resetStatus(ctx context.Context, now time.Time, lastResetMonth time.Month) time.Month {
	isTimeForResetStatus := now.Month() != lastResetMonth
	if isTimeForResetStatus {
		lastResetMonth = now.Month()
		err := b.userRepo.ResetSubmissionStatus(ctx)
		if err != nil {
			log.Printf("Reset status error: %v", err)
		} else {
			log.Printf("Reset submission status for new month %v", lastResetMonth)
		}
	}
	return lastResetMonth
}

func (b *Bot) sendReminds(ctx context.Context, now time.Time, wasRemindToday bool) bool {
	// firstDay := 20
	// lastDay := 25
	// startHour := 12
	// lastHour := 15

	firstDay := 6
	lastDay := 25
	startHour := 12
	lastHour := 15

	if now.Hour() < startHour && now.Hour() > lastHour {
		wasRemindToday = false
	}

	isTimeToSendReminders := now.Day() >= firstDay && now.Day() <= lastDay && now.Hour() >= startHour && now.Hour() <= lastHour
	if isTimeToSendReminders && !wasRemindToday {
		b.checkAndSendReminders(ctx)
		wasRemindToday = true
	}
	return wasRemindToday
}

func (b *Bot) checkAndSendReminders(ctx context.Context) {
	users, err := b.userRepo.GetShouldNotifyUsers(ctx)

	if err != nil {
		log.Printf("Get subscribed users error: %v", err)
		return
	}

	for _, userID := range users {
		remindMessage := getRemindMessage(userID)
		b.sender.SendMessage(remindMessage)
	}
}

var (
	remindText    = "⏰ Пора передать показания счетчиков!"
	remindButtons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Передал показания", "submitted"),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести показания", "generate_readings"),
			tgbotapi.NewInlineKeyboardButtonData("🔕 Отписаться", "unsubscribe"),
		),
	)
)

func getRemindMessage(userID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(userID, remindText)
	msg.ReplyMarkup = remindButtons
	return msg
}
