package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/utils"
	"os"
)

func ProcessShorteners(job *models.Queue, user *models.User, dir string, ctx context.Context) error {
	shorteners, err := user.Shorteners(ctx)
	if err != nil {
		return fmt.Errorf("fetching shorteners: %w", err)
	}

	shortenerJSON, err := json.MarshalIndent(shorteners, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling shorteners: %w", err)
	}

	exportFilePath := fmt.Sprintf("%s/shorteners.json", dir)
	file, err := os.Create(exportFilePath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(shortenerJSON)
	if err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	utils.Logger.WithFields(map[string]interface{}{
		"job_id":  job.ID,
		"user_id": user.ID,
	}).Info("Shortener export job completed...")

	return nil
}
