package models

import (
	"database/sql"
	"time"
)

type Queue struct {
	ID           uint64
	UserID       uint64
	Status       string
	AttemptCount uint8
	LastError    sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
