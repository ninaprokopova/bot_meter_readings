package storage

import (
	"context"
	"database/sql"
	"errors"
)

func (s *PostgresStorage) Subscribe(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (user_id, is_subscribed, has_submitted)
		VALUES ($1, TRUE, FALSE)
		ON CONFLICT (user_id) DO UPDATE
		SET is_subscribed = TRUE, last_reminder_date = NULL
	`, userID)
	return err
}

func (s *PostgresStorage) Unsubscribe(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE users 
		SET is_subscribed = FALSE
		WHERE user_id = $1
	`, userID)
	return err
}

func (s *PostgresStorage) MarkAsSubmitted(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE users 
		SET has_submitted = TRUE,
		    last_reminder_date = NOW()
		WHERE user_id = $1
	`, userID)
	return err
}

func (s *PostgresStorage) ShouldNotify(ctx context.Context, userID int64) (bool, error) {
	var (
		isSubscribed bool
		hasSubmitted bool
	)

	err := s.db.QueryRowContext(ctx, `
		SELECT is_subscribed, has_submitted
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&isSubscribed, &hasSubmitted)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return isSubscribed && !hasSubmitted, nil
}

func (s *PostgresStorage) ResetSubmissionStatus(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE users
		SET has_submitted = FALSE
		WHERE is_subscribed = TRUE
	`)
	return err
}

func (s *PostgresStorage) GetSubscribedUsers(ctx context.Context) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id
		FROM users
		WHERE is_subscribed = TRUE
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, nil
}
