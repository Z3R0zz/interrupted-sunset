package models

import (
	"context"
	"fmt"
	"interrupted-export/src/database"
	"time"
)

type Shortener struct {
	ID        string
	UserID    int
	TargetUrl string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Shorteners(ctx context.Context) ([]Shortener, error) {
	var shorteners []Shortener
	query := `SELECT id, user_id, target_url, slug, created_at, updated_at FROM url_shorteners WHERE user_id = ?`
	rows, err := database.DB.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	const layout = "2006-01-02 15:04:05"

	for rows.Next() {
		var s Shortener
		var createdAtRaw, updatedAtRaw []byte

		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.TargetUrl,
			&s.Slug,
			&createdAtRaw,
			&updatedAtRaw,
		); err != nil {
			return nil, err
		}

		s.CreatedAt, err = time.Parse(layout, string(createdAtRaw))
		if err != nil {
			return nil, fmt.Errorf("parsing created_at: %w", err)
		}

		s.UpdatedAt, err = time.Parse(layout, string(updatedAtRaw))
		if err != nil {
			return nil, fmt.Errorf("parsing updated_at: %w", err)
		}

		shorteners = append(shorteners, s)
	}

	return shorteners, nil
}
