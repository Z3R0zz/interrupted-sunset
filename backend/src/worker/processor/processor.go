package processor

import (
	"context"
	"fmt"
	"interrupted-export/src/models"
	"os"
)

func ProcessExportJob(job *models.Queue, user *models.User) error {
	fmt.Println("Processing export job for user:", user.Username)
	ctx := context.Background()

	exportDir := fmt.Sprintf("/tmp/export_user_%d", job.UserID)
	os.MkdirAll(exportDir, 0755)

	shorteners, err := user.Shorteners(ctx)
	if err != nil {
		return fmt.Errorf("fetching shorteners: %w", err)
	}

	fmt.Println("Found: ", shorteners)

	return nil
}
