package bot

import (
	"context"
	"errors"
	"submit_meter_readings/bot/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

func TestHandleMessage_NonCommand(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo:   userRepo,
		sender:     sender,
		userStates: make(map[int64]*UserState),
	}

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		From: &tgbotapi.User{ID: 456},
		Text: "просто текст",
	}

	bot.userStates[456] = &UserState{}
	sender.On("SendMessage", mock.Anything).Return()

	bot.handleMessage(ctx, msg)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StartCommand_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
	}

	userRepo.On("Subscribe", ctx, int64(456)).Return(nil)

	expectedMsg := tgbotapi.NewMessage(123, "✅ Вы подписались на напоминания!\nОни будут приходить c 12:00 по 15:00 с 20 по 25 число каждого месяца.")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StartCommand_DBError(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
	}

	expectedErr := errors.New("database error")
	userRepo.On("Subscribe", ctx, int64(456)).Return(expectedErr)

	expectedMsg := tgbotapi.NewMessage(123, "❌ Ошибка подписки. Попробуйте позже.")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StatusCommand_SubscribeNotSubmitted(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/status",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 7}},
	}

	userRepo.On("GetUserStatus", mock.Anything, int64(456)).Return(true, false, nil)

	expectedMsg := tgbotapi.NewMessage(123, "🔔 Вы подписаны на напоминания")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StatusCommand_SubscribeSubmitted(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/status",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 7}},
	}

	userRepo.On("GetUserStatus", mock.Anything, int64(456)).Return(true, true, nil)

	expectedMsg := tgbotapi.NewMessage(123, "✅ Вы уже передали показания в этом месяце")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StatusCommand_NotSubscribe(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/status",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 7}},
	}

	userRepo.On("GetUserStatus", mock.Anything, int64(456)).Return(false, false, nil)

	expectedMsg := tgbotapi.NewMessage(123, "🔕 Вы не подписаны на напоминания")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestHandleMessage_StatusCommand_DBError(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mocks.MockUserRepo)
	sender := new(mocks.MockMessageSender)

	bot := &Bot{
		userRepo: userRepo,
		sender:   sender,
	}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 123},
		From:     &tgbotapi.User{ID: 456},
		Text:     "/status",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 7}},
	}

	expectedErr := errors.New("database error")
	userRepo.On("GetUserStatus", ctx, int64(456)).Return(true, true, expectedErr)

	expectedMsg := tgbotapi.NewMessage(123, "❌ Ошибка проверки статуса.")
	sender.On("SendMessage", expectedMsg).Return()

	bot.handleMessage(ctx, msg)

	userRepo.AssertExpectations(t)
	sender.AssertExpectations(t)
}
