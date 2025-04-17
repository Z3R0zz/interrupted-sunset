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
	shorteners, errorsList, err := user.Shorteners(ctx)
	if err != nil {
		return fmt.Errorf("fetching shorteners: %w", err)
	}

	shortenerJSON, err := json.MarshalIndent(shorteners, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling shorteners: %w", err)
	}

	jsonDir := fmt.Sprintf("%s/jsons", dir)
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		return fmt.Errorf("failed to create json directory: %w", err)
	}

	exportFilePath := fmt.Sprintf("%s/shorteners.json", jsonDir)
	file, err := os.Create(exportFilePath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(shortenerJSON); err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	utils.Logger.WithFields(map[string]interface{}{
		"job_id":  job.ID,
		"user_id": user.ID,
	}).Info("Shortener export job completed...")

	if len(errorsList) > 0 {
		errFilePath := fmt.Sprintf("%s/shortener_errors.txt", dir)
		ef, err := os.OpenFile(errFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening shortener_errors.txt: %w", err)
		}
		defer ef.Close()

		for _, e := range errorsList {
			_, _ = ef.WriteString(e + "\n")
		}

		utils.Logger.WithFields(map[string]interface{}{
			"job_id":  job.ID,
			"user_id": user.ID,
		}).Warnf("Shortener export completed with %d row-level errors", len(errorsList))
	}

	return nil
}
