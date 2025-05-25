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

	var lastResetMonth time.Month

	var wasRemindToday = false

	for now := range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// вот это будто можно в отдельную функцию вынести
		if now.Day() == 1 && now.Hour() == 0 && now.Minute() == 0 && now.Month() != lastResetMonth {
			if err := b.userRepo.ResetSubmissionStatus(ctx); err != nil {
				log.Printf("Reset status error: %v", err)
			} else {
				lastResetMonth = now.Month()
				log.Println("Reset submission status for new month")
			}
		}

		firstDay := 20
		lastDay := 25
		startHour := 12
		lastHour := 20

		if now.Hour() < startHour && now.Hour() > lastHour {
			wasRemindToday = false
		}
		isTimeToSendReminders := now.Day() >= firstDay && now.Day() <= lastDay && now.Hour() >= startHour && now.Hour() <= lastHour

		if isTimeToSendReminders && !wasRemindToday {
			b.checkAndSendReminders(ctx, sendReminder)
			wasRemindToday = true
		}

		cancel()
	}
}

func (b *Bot) checkAndSendReminders(ctx context.Context, sendReminder func(int64, *Bot)) {
	users, err := b.userRepo.GetShouldNotifyUsers(ctx)

	if err != nil {
		log.Printf("Get subscribed users error: %v", err)
		return
	}

	for _, userID := range users {
		sendReminder(userID, b)
	}
}

func sendReminder(userID int64, b *Bot) {
	msg := getRemindMessage(userID)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Send reminder error to %d: %v", userID, err)
	} else {
		log.Printf("Reminder sent to %d", userID)
	}
}

func getRemindMessage(userID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(userID, "⏰ Пора передать показания счетчиков!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Передал показания", "submitted"),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести показания", "generate_readings"),
			tgbotapi.NewInlineKeyboardButtonData("🔕 Отписаться", "unsubscribe"),
		),
	)
	return msg
}
