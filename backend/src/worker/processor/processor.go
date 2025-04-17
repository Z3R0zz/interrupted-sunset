package processor

import (
	"context"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/worker/processor/processes"
	"os"
)

func ProcessExportJob(job *models.Queue, user *models.User) error {
	fmt.Println("Processing export job for user:", user.Username)
	ctx := context.Background()

	exportDir := fmt.Sprintf("tmp/export_user_%d", job.UserID)
	os.MkdirAll(exportDir, 0755)

	if err := processes.ProcessShorteners(job, user, exportDir, ctx); err != nil {
		return fmt.Errorf("processing shorteners: %w", err)
	}

	if err := processes.ProcessPastes(job, user, exportDir, ctx); err != nil {
		return fmt.Errorf("processing pastes: %w", err)
	}

	return nil
}
