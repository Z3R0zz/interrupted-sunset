package models

import (
	"context"
	"database/sql"
	"interrupted-export/src/database"
	"time"
)

type Queue struct {
	ID           uint64
	UserID       uint
	Status       string
	AttemptCount uint8
	LastError    sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (q *Queue) ExistsInQueue(ctx context.Context) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM queue WHERE user_id = ? LIMIT 1)`
	err := database.DB.QueryRowContext(ctx, query, q.UserID).Scan(&exists)
	return exists, err
}

func (q *Queue) GetStatus(ctx context.Context) (string, error) {
	var status string
	query := `SELECT status FROM queue WHERE user_id = ? LIMIT 1`
	err := database.DB.QueryRowContext(ctx, query, q.UserID).Scan(&status)
	return status, err
}

func (q *Queue) Insert(ctx context.Context) error {
	query := `
		INSERT INTO queue (user_id)
		VALUES (?)
	`
	_, err := database.DB.ExecContext(ctx, query, q.UserID)
	return err
}

func GetQueueByUserID(ctx context.Context, userID uint) (*Queue, error) {
	var q Queue
	query := `SELECT id, user_id, status, attempt_count, last_error, created_at, updated_at FROM queues WHERE user_id = ? LIMIT 1`
	err := database.DB.QueryRowContext(ctx, query, userID).Scan(
		&q.ID, &q.UserID, &q.Status, &q.AttemptCount,
		&q.LastError, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
