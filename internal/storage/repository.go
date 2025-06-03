package storage

import "context"

type UserRepository interface {
	Subscribe(ctx context.Context, userID int64) error
	Unsubscribe(ctx context.Context, userID int64) error
	MarkAsSubmitted(ctx context.Context, userID int64) error
	ResetSubmissionStatus(ctx context.Context) error
	SaveMeterReadings(ctx context.Context, userID int64,
		coldWater, hotWater, electricityDay, electricityNight int) error
	GetShouldNotifyUsers(ctx context.Context) ([]int64, error)
	GetUserStatus(ctx context.Context, userID int64) (bool, bool, error)
}
