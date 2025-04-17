package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/utils"
	"os"
)

func ProcessPastes(job *models.Queue, user *models.User, dir string, ctx context.Context) error {
	pastes, errorsList, err := user.Pastes(ctx)
	if err != nil {
		return fmt.Errorf("fetching pastes: %w", err)
	}

	pasteJSON, err := json.MarshalIndent(pastes, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling pastes: %w", err)
	}

	exportFilePath := fmt.Sprintf("%s/pastes.json", dir)
	file, err := os.Create(exportFilePath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(pasteJSON); err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	utils.Logger.WithFields(map[string]interface{}{
		"job_id":  job.ID,
		"user_id": user.ID,
	}).Info("Paste export job completed...")

	if len(errorsList) > 0 {
		errFilePath := fmt.Sprintf("%s/paste_errors.txt", dir)
		ef, err := os.Create(errFilePath)
		if err != nil {
			return fmt.Errorf("creating paste_errors.txt: %w", err)
		}
		defer ef.Close()

		for _, e := range errorsList {
			_, _ = ef.WriteString(e + "\n")
		}

		utils.Logger.WithFields(map[string]interface{}{
			"job_id":  job.ID,
			"user_id": user.ID,
		}).Warnf("Export completed with %d paste errors", len(errorsList))
	}

	return nil
}
