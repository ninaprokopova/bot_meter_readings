package bot

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetShouldNotifyUsers(ctx context.Context) ([]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockUserRepo) Subscribe(ctx context.Context, userID int64) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserRepo) Unsubscribe(ctx context.Context, userID int64) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserRepo) MarkAsSubmitted(ctx context.Context, userID int64) error {
	args := m.Called(ctx)
	return args.Error(0)
}
func (m *MockUserRepo) ShouldNotify(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockUserRepo) ResetSubmissionStatus(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
func (m *MockUserRepo) SaveMeterReadings(ctx context.Context, userID int64,
	coldWater, hotWater, electricityDay, electricityNight int) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockMessageSender struct {
	mock.Mock
}

func (m *MockMessageSender) SendMessage(msg tgbotapi.MessageConfig) {
	m.Called(msg)
}

func TestCheckAndSendReminders_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepo)
	userRepo.On("GetShouldNotifyUsers", ctx).Return([]int64{1, 2}, nil)

	mockSender := new(MockMessageSender)
	mockSender.On("SendMessage", mock.AnythingOfType("tgbotapi.MessageConfig")).Twice()

	bot := &Bot{
		userRepo: userRepo,
		sender:   mockSender,
	}

	bot.checkAndSendReminders(ctx)

	userRepo.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestCheckAndSendReminders_DBError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepo)
	expectedErr := errors.New("database error")
	userRepo.On("GetShouldNotifyUsers", ctx).Return([]int64{1, 2}, expectedErr)

	bot := &Bot{
		userRepo: userRepo,
	}

	bot.checkAndSendReminders(ctx)

	userRepo.AssertExpectations(t)
}

func TestGetRemainMessage(t *testing.T) {
	userId := int64(123)
	msg := getRemindMessage(userId)

	if msg.ChatID != userId {
		t.Errorf("Expected ChatId %v, got %v", userId, msg.ChatID)
	}

	markup, ok := msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	if !ok {
		t.Errorf("Expected InlineKeyboardMarkup, got %T", msg.ReplyMarkup)
	}

	if len(markup.InlineKeyboard) != 1 {
		t.Errorf("Expected 1 row of buttons, got %v", len(markup.InlineKeyboard))
	}

	if len(markup.InlineKeyboard[0]) != 3 {
		t.Errorf("Expected 3 botton, got %v", len(markup.InlineKeyboard[0]))
	}

	if *markup.InlineKeyboard[0][0].CallbackData != "submitted" {
		t.Errorf("Expected Callback 'submitted', got %v", markup.InlineKeyboard[0][0].Text)
	}

	if *markup.InlineKeyboard[0][1].CallbackData != "generate_readings" {
		t.Errorf("Expected Callback 'generate_readings', got %v", markup.InlineKeyboard[0][0].Text)
	}

	if *markup.InlineKeyboard[0][2].CallbackData != "unsubscribe" {
		t.Errorf("Expected Callback 'unsubscribe', got %v", markup.InlineKeyboard[0][0].Text)
	}
}
