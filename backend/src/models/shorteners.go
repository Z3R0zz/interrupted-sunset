package models

import (
	"context"
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

	for rows.Next() {
		var s Shortener
		if err := rows.Scan(&s.ID, &s.UserID, &s.TargetUrl, &s.Slug, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		shorteners = append(shorteners, s)
	}

	return shorteners, nil
}
