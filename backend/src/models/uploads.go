package models

import (
	"context"
	"fmt"
	"interrupted-export/src/database"
	"time"
)

type Upload struct {
	ID        int
	UserID    int
	Filename  string
	Folder    string
	Filesize  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Uploads(ctx context.Context) ([]Upload, []string, error) {
	var uploads []Upload
	var errors []string

	query := `SELECT id, uploaded_by, filename, foldername, filesize, created_at, updated_at FROM uploads WHERE uploaded_by = ?`
	rows, err := database.DB.QueryContext(ctx, query, u.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	const layout = "2006-01-02 15:04:05"

	for rows.Next() {
		var u Upload
		var createdAtRaw, updatedAtRaw []byte

		if err := rows.Scan(
			&u.ID,
			&u.UserID,
			&u.Filename,
			&u.Folder,
			&u.Filesize,
			&createdAtRaw,
			&updatedAtRaw,
		); err != nil {
			errors = append(errors, fmt.Sprintf("upload_id=%v: failed to scan upload row: %v", u.ID, err))
			continue
		}

		u.CreatedAt, err = time.Parse(layout, string(createdAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("upload_id=%v: invalid created_at: %v", u.ID, err))
			continue
		}
		u.UpdatedAt, err = time.Parse(layout, string(updatedAtRaw))
		if err != nil {
			errors = append(errors, fmt.Sprintf("upload_id=%v: invalid updated_at: %v", u.ID, err))
			continue
		}

		uploads = append(uploads, u)
	}

	return uploads, errors, nil
}
