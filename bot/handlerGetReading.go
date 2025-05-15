package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) startMeterReadingFlow(chatID, userID int64) {
	// Сохраняем состояние пользователя
	b.userStates[userID] = &UserState{
		CurrentStep: "cold_water",
	}

	msg := tgbotapi.NewMessage(chatID, "Введите показание счетчика холодной воды (целое число):")
	b.api.Send(msg)
}

func (b *Bot) handleMeterReadingInput(ctx context.Context, msg *tgbotapi.Message, state *UserState) error {
	value, err := strconv.Atoi(msg.Text)
	if err != nil {
		return b.sendTextReply(msg.Chat.ID, "Пожалуйста, введите целое число")
	}

	switch state.CurrentStep {
	case "cold_water":
		return b.handleColdWaterInput(state, value, msg)
	case "hot_water":
		return b.handleHotWaterInput(state, value, msg)
	case "electricity_day":
		return b.handleDayElectricityInput(state, value, msg)
	case "electricity_night":
		return b.handleNightElectricityInput(ctx, state, value, msg)
	}

	return nil
}

func (b *Bot) sendTextReply(chatID int64, text string) error {
	_, err := b.api.Send(tgbotapi.NewMessage(chatID, text))
	return err
}

func (b *Bot) handleColdWaterInput(state *UserState, value int, msg *tgbotapi.Message) error {
	state.Readings.ColdWater = value
	state.CurrentStep = "hot_water"
	return b.sendTextReply(msg.Chat.ID, "Введите показание счетчика горячей воды:")
}

func (b *Bot) handleHotWaterInput(state *UserState, value int, msg *tgbotapi.Message) error {
	state.Readings.HotWater = value
	state.CurrentStep = "electricity_day"
	return b.sendTextReply(msg.Chat.ID, "Введите показание счетчика дневного тарифа:")
}

func (b *Bot) handleDayElectricityInput(state *UserState, value int, msg *tgbotapi.Message) error {
	state.Readings.ElectricityDay = value
	state.CurrentStep = "electricity_night"
	return b.sendTextReply(msg.Chat.ID, "Введите показание счетчика ночного тарифа:")
}

func (b *Bot) handleNightElectricityInput(ctx context.Context, state *UserState, value int, msg *tgbotapi.Message) error {
	state.Readings.ElectricityNight = value

	// Сохраняем показания
	if err := b.saveReadings(ctx, msg, state); err != nil {
		b.sendTextReply(msg.Chat.ID, "Ошибка сохранения показаний")
		delete(b.userStates, msg.From.ID)
		return nil
	}

	// Формируем отчет
	report := b.getReport(state)
	b.sendTextReply(msg.Chat.ID, report)
	delete(b.userStates, msg.From.ID)
	return nil
}

func (b *Bot) saveReadings(ctx context.Context, msg *tgbotapi.Message, state *UserState) error {
	err := b.userRepo.SaveMeterReadings(ctx, msg.From.ID,
		state.Readings.ColdWater,
		state.Readings.HotWater,
		state.Readings.ElectricityDay,
		state.Readings.ElectricityNight)
	return err
}

func (*Bot) getReport(state *UserState) string {
	now := time.Now()
	monthName := getMonthName(now.Month())

	report := fmt.Sprintf(
		"Показания счетчиков %s %d:\n"+
			"Холодная вода: %d\n"+
			"Горячая вода: %d\n"+
			"Электричество день (T1): %d\n"+
			"Электричество ночь (T2): %d",
		monthName,
		now.Year(),
		state.Readings.ColdWater,
		state.Readings.HotWater,
		state.Readings.ElectricityDay,
		state.Readings.ElectricityNight)
	return report
}

func getMonthName(m time.Month) string {
	months := map[time.Month]string{
		time.January:   "январь",
		time.February:  "февраль",
		time.March:     "март",
		time.April:     "апрель",
		time.May:       "май",
		time.June:      "июнь",
		time.July:      "июль",
		time.August:    "август",
		time.September: "сентябрь",
		time.October:   "октябрь",
		time.November:  "ноябрь",
		time.December:  "декабрь",
	}
	return months[m]
}
