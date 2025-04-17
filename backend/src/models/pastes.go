package models

import (
	"context"
	"fmt"
	"interrupted-export/src/database"
	"interrupted-export/src/services"
	"time"
)

type Paste struct {
	ID           string
	UserID       int
	Slug         string
	Folder       string
	Title        string
	Lang         string
	Content      *string
	PGPSignature *string
	PGPKey       *string
	Size         int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) Pastes(ctx context.Context) ([]Paste, error) {
	var pastes []Paste
	query := `SELECT id, uploaded_by, slug, foldername, title, lang, content, pgp_signature, pgp_key, filesize, created_at, updated_at FROM pastes WHERE user_id = ?`
	rows, err := database.DB.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	const layout = "2006-01-02 15:04:05"

	for rows.Next() {
		var p Paste
		var createdAtRaw, updatedAtRaw []byte

		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Slug,
			&p.Folder,
			&p.Title,
			&p.Lang,
			&p.Content,
			&p.PGPSignature,
			&p.PGPKey,
			&p.Size,
			&createdAtRaw,
			&updatedAtRaw,
		); err != nil {
			return nil, err
		}

		p.CreatedAt, err = time.Parse(layout, string(createdAtRaw))
		if err != nil {
			return nil, fmt.Errorf("parsing created_at: %w", err)
		}

		p.UpdatedAt, err = time.Parse(layout, string(updatedAtRaw))
		if err != nil {
			return nil, fmt.Errorf("parsing updated_at: %w", err)
		}

		if p.Content == nil {
			path := fmt.Sprintf("%s/%s", p.Folder, p.Slug)
			data, err := services.R2.GetObject(ctx, path)
			if err != nil {
				return nil, fmt.Errorf("failed to get object from R2: %w", err)
			}

			content := string(data)
			p.Content = &content
		}

		pastes = append(pastes, p)
	}

	return pastes, nil
}
