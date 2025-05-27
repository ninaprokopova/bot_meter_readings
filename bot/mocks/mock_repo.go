package mocks

import (
	"context"

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
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) Unsubscribe(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) MarkAsSubmitted(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockUserRepo) ShouldNotify(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockUserRepo) ResetSubmissionStatus(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
func (m *MockUserRepo) SaveMeterReadings(ctx context.Context, userID int64,
	coldWater, hotWater, electricityDay, electricityNight int) error {
	args := m.Called(ctx, userID, coldWater, hotWater, electricityDay, electricityNight)
	return args.Error(0)
}
