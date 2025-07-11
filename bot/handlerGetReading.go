package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) startMeterReadingFlow(chatID, userID, msgID int64) {
	b.userStates[userID] = &UserState{
		CurrentStep: "cold_water",
	}

	msg := tgbotapi.NewMessage(chatID, "Введите показание счетчика холодной воды (целое число):")
	b.sender.SendMessage(msg)
	msgForDelete := tgbotapi.NewDeleteMessage(chatID, int(msgID))
	b.deleter.DeleteMessage(msgForDelete)
}

func (b *Bot) handleMeterReadingInput(ctx context.Context, msg *tgbotapi.Message, state *UserState) {
	value, err := strconv.Atoi(msg.Text)
	if err != nil || value < 0 {
		msgWriteInteger := tgbotapi.NewMessage(msg.Chat.ID, "Пожалуйста, введите целое число")
		b.sender.SendMessage(msgWriteInteger)
		return
	}

	switch state.CurrentStep {
	case "cold_water":
		b.handleColdWaterInput(state, value, msg)
	case "hot_water":
		b.handleHotWaterInput(state, value, msg)
	case "electricity_day":
		b.handleDayElectricityInput(state, value, msg)
	case "electricity_night":
		b.handleNightElectricityInput(ctx, state, value, msg)
	}
}

// Можно на самом деле просто чиселкой chat.ID передавать, чтобы msg не протаскивать вниз
func (b *Bot) handleColdWaterInput(state *UserState, value int, msg *tgbotapi.Message) {
	state.Readings.ColdWater = value
	state.CurrentStep = "hot_water"
	msgHotWater := tgbotapi.NewMessage(msg.Chat.ID, "Введите показание счетчика горячей воды:")
	b.sender.SendMessage(msgHotWater)
}

func (b *Bot) handleHotWaterInput(state *UserState, value int, msg *tgbotapi.Message) {
	state.Readings.HotWater = value
	state.CurrentStep = "electricity_day"
	msgDayElectricity := tgbotapi.NewMessage(msg.Chat.ID, "Введите показание счетчика дневного тарифа:")
	b.sender.SendMessage(msgDayElectricity)
}

func (b *Bot) handleDayElectricityInput(state *UserState, value int, msg *tgbotapi.Message) {
	state.Readings.ElectricityDay = value
	state.CurrentStep = "electricity_night"
	msgNightElecticity := tgbotapi.NewMessage(msg.Chat.ID, "Введите показание счетчика ночного тарифа:")
	b.sender.SendMessage(msgNightElecticity)
}

func (b *Bot) handleNightElectricityInput(ctx context.Context, state *UserState, value int, msg *tgbotapi.Message) {
	state.Readings.ElectricityNight = value

	err := b.saveReadings(ctx, msg, state)
	if err != nil {
		msgSaveError := tgbotapi.NewMessage(msg.Chat.ID, "Ошибка сохранения показаний")
		b.sender.SendMessage(msgSaveError)
		delete(b.userStates, msg.From.ID)
		return
	}

	report := b.getReport(ctx, state, uint64(msg.From.ID))
	msgReport := tgbotapi.NewMessage(msg.Chat.ID, report)
	b.sender.SendMessage(msgReport)
	delete(b.userStates, msg.From.ID)
}

func (b *Bot) saveReadings(ctx context.Context, msg *tgbotapi.Message, state *UserState) error {
	err := b.userRepo.SaveMeterReadings(ctx, msg.From.ID,
		state.Readings.ColdWater,
		state.Readings.HotWater,
		state.Readings.ElectricityDay,
		state.Readings.ElectricityNight)
	return err
}

var timeNow = time.Now

func (b *Bot) getReport(ctx context.Context, state *UserState, userID uint64) string {
	now := timeNow()
	monthName := getMonthName(now.Month())
	template, _ := b.userRepo.GetTemplate(ctx, userID)
	report := fmt.Sprintf(
		"\nПоказания счетчиков %s %d:\n"+
			"Холодная вода: %d\n"+
			"Горячая вода: %d\n"+
			"Электричество день (T1): %d\n"+
			"Электричество ночь (T2): %d\n",
		monthName,
		now.Year(),
		state.Readings.ColdWater,
		state.Readings.HotWater,
		state.Readings.ElectricityDay,
		state.Readings.ElectricityNight)
	reportInTemplate := strings.Replace(template, "*показания*", report, 1)
	return reportInTemplate
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
