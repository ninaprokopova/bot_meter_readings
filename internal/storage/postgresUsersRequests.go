package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
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

func (s *PostgresStorage) GetUserStatus(ctx context.Context, userID int64) (isSubscribed, hasSubmitted bool, err error) {
	err = s.db.QueryRowContext(ctx, `
		SELECT is_subscribed, has_submitted
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&isSubscribed, &hasSubmitted)

	if errors.Is(err, sql.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, nil
	}

	return isSubscribed, hasSubmitted, nil
}

func (s *PostgresStorage) ResetSubmissionStatus(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE users
		SET has_submitted = FALSE
	`)
	return err
}

func (s *PostgresStorage) GetShouldNotifyUsers(ctx context.Context) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id
		FROM users
		WHERE is_subscribed = true AND has_submitted = false
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

func (s *PostgresStorage) ChangeTemplate(ctx context.Context, userID uint64, newTemplate string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO template (user_id, template)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET template = $2
	`, userID, newTemplate)

	return err
}

func (s *PostgresStorage) GetTemplate(ctx context.Context, userID uint64) (template string, err error) {
	err = s.db.QueryRowContext(ctx, `
		SELECT template
		FROM template
		WHERE user_id = $1
	`, userID).Scan(&template)
	if err != nil {
		log.Println(err)
		template = "*показания*"
		return template, err
	}
	return template, err
}
