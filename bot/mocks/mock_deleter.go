package mocks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

type MockMessageDeleter struct {
	mock.Mock
}

func (m *MockMessageDeleter) DeleteMessage(msg tgbotapi.DeleteMessageConfig) {
	m.Called(msg)
}
