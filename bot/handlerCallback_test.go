package bot

import (
	"context"
	"errors"
	"submit_meter_readings/bot/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCallback_SubmittedSuccess(t *testing.T) {
	query := &tgbotapi.CallbackQuery{
		From: &tgbotapi.User{ID: 123},
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 456},
		},
		Data: "submitted",
	}

	mockSender := new(mocks.MockMessageSender)
	mockSender.On("SendMessage", mock.AnythingOfType("tgbotapi.MessageConfig")).Once()

	mockDeleter := new(mocks.MockMessageDeleter)
	mockDeleter.On("DeleteMessage", mock.AnythingOfType("tgbotapi.DeleteMessageConfig")).Once()

	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepo)
	mockRepo.On("MarkAsSubmitted", ctx, query.From.ID).Return(nil)

	testBot := &Bot{
		sender:     mockSender,
		deleter:    mockDeleter,
		userRepo:   mockRepo,
		userStates: make(map[int64]*UserState),
	}

	testBot.handleCallback(context.Background(), query)

	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestHandleCallback_SubmittedDBError(t *testing.T) {
	query := &tgbotapi.CallbackQuery{
		From: &tgbotapi.User{ID: 123},
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 456},
		},
		Data: "submitted",
	}

	mockSender := new(mocks.MockMessageSender)
	mockSender.On("SendMessage", mock.AnythingOfType("tgbotapi.MessageConfig")).Once()

	mockDeleter := new(mocks.MockMessageDeleter)
	mockDeleter.On("DeleteMessage", mock.AnythingOfType("tgbotapi.DeleteMessageConfig")).Once()

	mockRepo := new(mocks.MockUserRepo)
	expectedErr := errors.New("database connection failed")
	mockRepo.On("MarkAsSubmitted", mock.Anything, query.From.ID).Return(expectedErr)

	testBot := &Bot{
		sender:     mockSender,
		deleter:    mockDeleter,
		userRepo:   mockRepo,
		userStates: make(map[int64]*UserState),
	}

	testBot.handleCallback(context.Background(), query)

	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestHandleCallback_UnsubscribeSuccess(t *testing.T) {
	query := &tgbotapi.CallbackQuery{
		From: &tgbotapi.User{ID: 123},
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 456},
		},
		Data: "unsubscribe",
	}

	mockSender := new(mocks.MockMessageSender)
	mockSender.On("SendMessage", mock.AnythingOfType("tgbotapi.MessageConfig")).Once()

	mockDeleter := new(mocks.MockMessageDeleter)
	mockDeleter.On("DeleteMessage", mock.AnythingOfType("tgbotapi.DeleteMessageConfig")).Once()

	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepo)
	mockRepo.On("Unsubscribe", ctx, query.From.ID).Return(nil)

	testBot := &Bot{
		sender:     mockSender,
		deleter:    mockDeleter,
		userRepo:   mockRepo,
		userStates: make(map[int64]*UserState),
	}

	testBot.handleCallback(context.Background(), query)

	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestHandleCallback_UnsubscribeDBError(t *testing.T) {
	query := &tgbotapi.CallbackQuery{
		From: &tgbotapi.User{ID: 123},
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 456},
		},
		Data: "unsubscribe",
	}

	mockSender := new(mocks.MockMessageSender)
	mockSender.On("SendMessage", mock.AnythingOfType("tgbotapi.MessageConfig")).Once()

	mockDeleter := new(mocks.MockMessageDeleter)
	mockDeleter.On("DeleteMessage", mock.AnythingOfType("tgbotapi.DeleteMessageConfig")).Once()

	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepo)
	expectedErr := errors.New("database connection failed")
	mockRepo.On("Unsubscribe", ctx, query.From.ID).Return(expectedErr)

	testBot := &Bot{
		sender:     mockSender,
		deleter:    mockDeleter,
		userRepo:   mockRepo,
		userStates: make(map[int64]*UserState),
	}

	testBot.handleCallback(context.Background(), query)

	mockSender.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestStartMeterReadingFlow(t *testing.T) {
	mockSender := new(mocks.MockMessageSender)
	mockRepo := new(mocks.MockUserRepo)
	mockDeleter := new(mocks.MockMessageDeleter)

	bot := &Bot{
		sender:     mockSender,
		userRepo:   mockRepo,
		deleter:    mockDeleter,
		userStates: make(map[int64]*UserState),
	}

	query := &tgbotapi.CallbackQuery{
		From: &tgbotapi.User{ID: 123},
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 456},
		},
		Data: "generate_readings",
	}

	expectedMsg := tgbotapi.NewMessage(query.Message.Chat.ID, "Введите показание счетчика холодной воды (целое число):")
	mockSender.On("SendMessage", expectedMsg).Return()
	mockDeleter.On("DeleteMessage", mock.AnythingOfType("tgbotapi.DeleteMessageConfig")).Once()

	ctx := context.Background()
	bot.handleCallback(ctx, query)

	assert.NotNil(t, bot.userStates[123])
	assert.Equal(t, "cold_water", bot.userStates[123].CurrentStep)
	mockSender.AssertExpectations(t)
	mockDeleter.AssertExpectations(t)
}
