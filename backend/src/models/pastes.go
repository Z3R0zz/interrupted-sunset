package models

import (
	"context"
	"fmt"
	"interrupted-export/src/database"
	"interrupted-export/src/services"
	"interrupted-export/src/utils"
	"time"
)

type Paste struct {
	ID           int
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

func (u *User) Pastes(ctx context.Context) ([]Paste, []string, error) {
	var pastes []Paste
	var errors []string

	query := `SELECT id, uploaded_by, slug, foldername, title, lang, content, pgp_signature, pgp_key, filesize, created_at, updated_at FROM pastes WHERE uploaded_by = ?`
	rows, err := database.DB.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, nil, err
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
			errors = append(errors, fmt.Sprintf("paste_id=%v: failed to scan paste row: %v", p.ID, err))
			continue
		}

		p.CreatedAt, err = time.Parse(layout, string(createdAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("paste_id=%v: invalid created_at: %v", p.ID, err))
			continue
		}
		p.UpdatedAt, err = time.Parse(layout, string(updatedAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("paste_id=%v: invalid updated_at: %v", p.ID, err))
			continue
		}

		if p.Content == nil {
			path := fmt.Sprintf("%s/%s", p.Folder, p.Slug)
			var data []byte
			success := false

			for i := 1; i <= 3; i++ {
				data, err = services.R2.GetObject(ctx, path)
				if err == nil {
					success = true
					break
				}
				utils.Logger.WithError(err).WithField("paste_id", p.ID).Warnf("Retry %d: failed to fetch paste content from R2", i)
				time.Sleep(500 * time.Millisecond)
			}

			if !success {
				errors = append(errors, fmt.Sprintf("paste_id=%v: failed to fetch content after 3 attempts", p.ID))
				continue
			}

			content := string(data)
			p.Content = &content
		}

		pastes = append(pastes, p)
	}

	return pastes, errors, nil
}
