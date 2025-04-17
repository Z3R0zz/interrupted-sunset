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
	pastes, err := user.Pastes(ctx)
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

	_, err = file.Write(pasteJSON)
	if err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	utils.Logger.WithFields(map[string]interface{}{
		"job_id":  job.ID,
		"user_id": user.ID,
	}).Info("Paste export job completed...")

	return nil
}
