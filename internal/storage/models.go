package storage

import "time"

type User struct {
	ID               int64
	IsSubscribed     bool
	HasSubmitted     bool
	LastReminderDate *time.Time
}
