package storage

import "context"

type UserRepository interface {
	Subscribe(ctx context.Context, userID int64) error
	Unsubscribe(ctx context.Context, userID int64) error
	MarkAsSubmitted(ctx context.Context, userID int64) error
	ShouldNotify(ctx context.Context, userID int64) (bool, error)
	ResetSubmissionStatus(ctx context.Context) error
	GetSubscribedUsers(ctx context.Context) ([]int64, error)
	SaveMeterReadings(ctx context.Context, userID int64,
		coldWater, hotWater, electricityDay, electricityNight int) error
}
