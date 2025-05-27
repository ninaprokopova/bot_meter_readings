package bot

import (
	"context"
	"errors"
	"submit_meter_readings/bot/mocks"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetReading_GetMonthName(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Month
		expected string
	}{
		{"January", time.January, "январь"},
		{"February", time.February, "февраль"},
		{"March", time.March, "март"},
		{"April", time.April, "апрель"},
		{"May", time.May, "май"},
		{"June", time.June, "июнь"},
		{"July", time.July, "июль"},
		{"August", time.August, "август"},
		{"September", time.September, "сентябрь"},
		{"October", time.October, "октябрь"},
		{"November", time.November, "ноябрь"},
		{"December", time.December, "декабрь"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMonthName(tt.input)
			if got != tt.expected {
				t.Errorf("getMonthName(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestHandlerGetReadings_getReport(t *testing.T) {
	fixedTime := time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC)

	originalNow := timeNow
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = originalNow }()

	state := &UserState{
		Readings: MeterReadings{
			ColdWater:        1,
			HotWater:         2,
			ElectricityDay:   3,
			ElectricityNight: 4,
		},
	}
	expectedReport := "Показания счетчиков май 2024:\n" +
		"Холодная вода: 1\n" +
		"Горячая вода: 2\n" +
		"Электричество день (T1): 3\n" +
		"Электричество ночь (T2): 4"

	bot := &Bot{}
	t.Run("Test_getReport", func(t *testing.T) {
		got := bot.getReport(state)
		if got != expectedReport {
			t.Errorf("getReport() = \n%v\n, want \n%v", got, expectedReport)
		}
	})
}

func TestHandlerGetReadings_saveReadings(t *testing.T) {
	ctx := context.Background()

	testUserID := int64(123)
	testReadings := &UserState{
		Readings: MeterReadings{
			ColdWater:        1,
			HotWater:         2,
			ElectricityDay:   3,
			ElectricityNight: 4,
		},
	}

	mockRepo := new(mocks.MockUserRepo)
	mockRepo.On("SaveMeterReadings", ctx, testUserID,
		testReadings.Readings.ColdWater,
		testReadings.Readings.HotWater,
		testReadings.Readings.ElectricityDay,
		testReadings.Readings.ElectricityNight).Return(nil)

	testMsg := &tgbotapi.Message{
		From: &tgbotapi.User{ID: testUserID},
	}

	bot := &Bot{
		userRepo: mockRepo,
	}

	err := bot.saveReadings(ctx, testMsg, testReadings)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHandlerGetReadings_handleColdWaterInput(t *testing.T) {
	testChatID := int64(123)
	testValue := 10
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
	}

	state := &UserState{
		CurrentStep: "cold_water",
		Readings:    MeterReadings{},
	}

	mockSender := new(mocks.MockMessageSender)
	expectedMsg := tgbotapi.NewMessage(testChatID, "Введите показание счетчика горячей воды:")
	mockSender.On("SendMessage", expectedMsg).Return(nil)
	bot := &Bot{
		sender: mockSender,
	}

	bot.handleColdWaterInput(state, testValue, testMsg)

	assert.Equal(t, testValue, state.Readings.ColdWater, "ColdWater value should be updated")
	assert.Equal(t, "hot_water", state.CurrentStep, "CurrentStep should be updated to hot_water")

	mockSender.AssertExpectations(t)
}

func TestHandlerGetReadings_handleHotWaterInput(t *testing.T) {
	testChatID := int64(123)
	testValue := 15
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
	}

	state := &UserState{
		CurrentStep: "hot_water",
		Readings:    MeterReadings{},
	}

	mockSender := new(mocks.MockMessageSender)
	expectedMsg := tgbotapi.NewMessage(testChatID, "Введите показание счетчика дневного тарифа:")
	mockSender.On("SendMessage", expectedMsg).Return(nil)
	bot := &Bot{
		sender: mockSender,
	}

	bot.handleHotWaterInput(state, testValue, testMsg)

	assert.Equal(t, testValue, state.Readings.HotWater, "HotWater value should be updated")
	assert.Equal(t, "electricity_day", state.CurrentStep, "CurrentStep should be updated to electricity_day")

	mockSender.AssertExpectations(t)
}

func TestHandlerGetReadings_handleDayElectricityInput(t *testing.T) {
	testChatID := int64(123)
	testValue := 20
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
	}

	state := &UserState{
		CurrentStep: "electricity_day",
		Readings:    MeterReadings{},
	}

	mockSender := new(mocks.MockMessageSender)
	expectedMsg := tgbotapi.NewMessage(testChatID, "Введите показание счетчика ночного тарифа:")
	mockSender.On("SendMessage", expectedMsg).Return(nil)
	bot := &Bot{
		sender: mockSender,
	}

	bot.handleDayElectricityInput(state, testValue, testMsg)

	assert.Equal(t, testValue, state.Readings.ElectricityDay, "ElectricityNight value should be updated")
	assert.Equal(t, "electricity_night", state.CurrentStep, "CurrentStep should be updated to electricity_night")

	mockSender.AssertExpectations(t)
}

func TestHandlerGetReadings_handleNightElectricityInput(t *testing.T) {
	testChatID := int64(123)
	testValue := 20
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
	}

	state := &UserState{
		CurrentStep: "electricity_day",
		Readings:    MeterReadings{},
	}

	mockSender := new(mocks.MockMessageSender)
	expectedMsg := tgbotapi.NewMessage(testChatID, "Введите показание счетчика ночного тарифа:")
	mockSender.On("SendMessage", expectedMsg).Return(nil)
	bot := &Bot{
		sender: mockSender,
	}

	bot.handleDayElectricityInput(state, testValue, testMsg)

	assert.Equal(t, testValue, state.Readings.ElectricityDay, "ElectricityNight value should be updated")
	assert.Equal(t, "electricity_night", state.CurrentStep, "CurrentStep should be updated to electricity_night")

	mockSender.AssertExpectations(t)
}

func TestHandlerGetReadings_handleNightElectricityInputSuccess(t *testing.T) {
	ctx := context.Background()
	testChatID := int64(123)
	testUserID := int64(456)
	testValue := 30
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
		From: &tgbotapi.User{ID: testUserID},
	}

	testState := &UserState{
		CurrentStep: "electricity_night",
		Readings: MeterReadings{
			ColdWater:        5,
			HotWater:         6,
			ElectricityDay:   7,
			ElectricityNight: 0,
		},
	}

	reportState := &UserState{
		CurrentStep: "electricity_night",
		Readings: MeterReadings{
			ColdWater:        5,
			HotWater:         6,
			ElectricityDay:   7,
			ElectricityNight: 30,
		},
	}

	mockRepo := new(mocks.MockUserRepo)
	mockSender := new(mocks.MockMessageSender)

	bot := &Bot{
		sender:     mockSender,
		userRepo:   mockRepo,
		userStates: map[int64]*UserState{testUserID: testState},
	}

	mockRepo.On("SaveMeterReadings", ctx, testUserID,
		testState.Readings.ColdWater,
		testState.Readings.HotWater,
		testState.Readings.ElectricityDay,
		testValue).Return(nil)

	expectedReport := bot.getReport(reportState)
	expectedMsg := tgbotapi.NewMessage(testChatID, expectedReport)
	mockSender.On("SendMessage", expectedMsg).Return(nil)

	bot.handleNightElectricityInput(ctx, testState, testValue, testMsg)

	assert.Equal(t, testValue, testState.Readings.ElectricityNight)
	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	assert.NotContains(t, bot.userStates, testUserID, "User state should be deleted")
}

func TestHandlerGetReadings_handleNightElectricityInputError(t *testing.T) {
	ctx := context.Background()
	testChatID := int64(123)
	testUserID := int64(456)
	testValue := 30
	testMsg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: testChatID},
		From: &tgbotapi.User{ID: testUserID},
	}

	testState := &UserState{
		CurrentStep: "electricity_night",
		Readings: MeterReadings{
			ColdWater:        5,
			HotWater:         6,
			ElectricityDay:   7,
			ElectricityNight: 0,
		},
	}

	mockRepo := new(mocks.MockUserRepo)
	mockSender := new(mocks.MockMessageSender)

	bot := &Bot{
		sender:     mockSender,
		userRepo:   mockRepo,
		userStates: map[int64]*UserState{testUserID: testState},
	}

	expectedErr := errors.New("database error")
	mockRepo.On("SaveMeterReadings", ctx, testUserID,
		testState.Readings.ColdWater,
		testState.Readings.HotWater,
		testState.Readings.ElectricityDay,
		testValue).Return(expectedErr)

	msgSaveError := tgbotapi.NewMessage(testChatID, "Ошибка сохранения показаний")
	mockSender.On("SendMessage", msgSaveError).Return(nil)

	bot.handleNightElectricityInput(ctx, testState, testValue, testMsg)

	assert.Equal(t, testValue, testState.Readings.ElectricityNight)
	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	assert.NotContains(t, bot.userStates, testUserID, "User state should be deleted")
}
