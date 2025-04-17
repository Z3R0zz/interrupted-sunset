package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/utils"
	"os"
)

func ProcessUploads(job *models.Queue, user *models.User, dir string, ctx context.Context) error {
	uploads, errorsList, err := user.Uploads(ctx)
	if err != nil {
		return fmt.Errorf("fetching uploads: %w", err)
	}

	uploadJSON, err := json.MarshalIndent(uploads, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling uploads: %w", err)
	}

	exportFilePath := fmt.Sprintf("%s/uploads.json", dir)
	file, err := os.Create(exportFilePath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(uploadJSON); err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	utils.Logger.WithFields(map[string]interface{}{
		"job_id":  job.ID,
		"user_id": user.ID,
	}).Info("Uploads export job completed...")

	if len(errorsList) > 0 {
		errFilePath := fmt.Sprintf("%s/uploads_errors.txt", dir)
		ef, err := os.OpenFile(errFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening uploads_errors.txt: %w", err)
		}
		defer ef.Close()

		for _, e := range errorsList {
			_, _ = ef.WriteString(e + "\n")
		}

		utils.Logger.WithFields(map[string]interface{}{
			"job_id":  job.ID,
			"user_id": user.ID,
		}).Warnf("Uploads export completed with %d row-level errors", len(errorsList))
	}

	return nil
}
