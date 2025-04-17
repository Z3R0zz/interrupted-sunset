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

func (u *User) Shorteners(ctx context.Context) ([]Shortener, []string, error) {
	var shorteners []Shortener
	var errors []string

	query := `SELECT id, user_id, target_url, slug, created_at, updated_at FROM url_shorteners WHERE user_id = ?`
	rows, err := database.DB.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, nil, err
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
			errors = append(errors, fmt.Sprintf("shortener_id=%s: failed to scan row: %v", s.ID, err))
			continue
		}

		s.CreatedAt, err = time.Parse(layout, string(createdAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("shortener_id=%s: invalid created_at: %v", s.ID, err))
			continue
		}
		s.UpdatedAt, err = time.Parse(layout, string(updatedAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("shortener_id=%s: invalid updated_at: %v", s.ID, err))
			continue
		}

		shorteners = append(shorteners, s)
	}

	return shorteners, errors, nil
}
