package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) startReminder() {
	// Первая проверка при запуске
	b.checkAndSendReminders(context.Background())

	ticker := time.NewTicker(1 * time.Minute) // Проверка каждую минуту
	defer ticker.Stop()

	var lastResetMonth time.Month

	for now := range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Сброс статусов 1-го числа в 00:00
		if now.Day() == 1 && now.Hour() == 0 && now.Minute() == 0 && now.Month() != lastResetMonth {
			if err := b.userRepo.ResetSubmissionStatus(ctx); err != nil {
				log.Printf("Reset status error: %v", err)
			} else {
				lastResetMonth = now.Month()
				log.Println("Reset submission status for new month")
			}
		}

		// Отправка напоминаний в 12:00 с 20 по 25 число
		if now.Day() >= 20 && now.Day() <= 24 && now.Hour() == 12 && now.Minute() == 00 {
			b.checkAndSendReminders(ctx)
		}

		cancel()
	}
}

func (b *Bot) checkAndSendReminders(ctx context.Context) {
	users, err := b.userRepo.GetSubscribedUsers(ctx)
	if err != nil {
		log.Printf("Get subscribed users error: %v", err)
		return
	}

	for _, userID := range users {
		shouldNotify, err := b.userRepo.ShouldNotify(ctx, userID)
		if err != nil {
			log.Printf("Check notification status error for user %d: %v", userID, err)
			continue
		}

		if !shouldNotify {
			continue
		}

		b.sendReminder(userID)
	}
}

func (b *Bot) sendReminder(userID int64) {
	msg := tgbotapi.NewMessage(userID, "⏰ Пора передать показания счетчиков!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Передал показания", "submitted"),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести показания", "generate_readings"),
			tgbotapi.NewInlineKeyboardButtonData("🔕 Отписаться", "unsubscribe"),
		),
	)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Send reminder error to %d: %v", userID, err)
	} else {
		log.Printf("Reminder sent to %d", userID)
	}
}
