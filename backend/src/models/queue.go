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

type Job struct {
	ID     int
	UserID int
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

func FetchJob(ctx context.Context) (*Queue, error) {
	row := database.DB.QueryRowContext(ctx, `
        SELECT id, user_id, status, attempt_count, last_error, created_at, updated_at FROM queue
        WHERE status = 'waiting'
        ORDER BY created_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `)

	var q Queue
	err := row.Scan(
		&q.ID, &q.UserID, &q.Status, &q.AttemptCount,
		&q.LastError, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Queue) MarkProcessing(ctx context.Context) error {
	_, err := database.DB.ExecContext(ctx, `UPDATE queue SET status = 'processing', updated_at = NOW() WHERE id = ?`, q.ID)
	return err
}

func (q *Queue) MarkDone() {
	database.DB.Exec(`UPDATE queue SET status = 'done', updated_at = NOW() WHERE id = ?`, q.ID)
}

func (q *Queue) MarkFailed(errStr string) {
	database.DB.Exec(`
        UPDATE queue
        SET status = 'failed', attempt_count = attempt_count + 1, last_error = ?, updated_at = NOW()
        WHERE id = ?
    `, errStr, q.ID)
}
