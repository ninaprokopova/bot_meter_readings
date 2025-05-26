package mocks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

type MockMessageSender struct {
	mock.Mock
}

func (m *MockMessageSender) SendMessage(msg tgbotapi.MessageConfig) {
	m.Called(msg)
}
