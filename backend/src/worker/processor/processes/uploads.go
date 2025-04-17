package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/services"
	"interrupted-export/src/utils"
	"os"
	"time"
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

	jsonDir := fmt.Sprintf("%s/jsons", dir)
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		return fmt.Errorf("failed to create json directory: %w", err)
	}

	exportFilePath := fmt.Sprintf("%s/uploads.json", jsonDir)
	file, err := os.Create(exportFilePath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(uploadJSON); err != nil {
		return fmt.Errorf("writing to export file: %w", err)
	}

	uploadDir := fmt.Sprintf("%s/uploads", dir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload export dir: %w", err)
	}

	for _, upload := range uploads {
		key := fmt.Sprintf("%s/%s", upload.Folder, upload.Filename)

		data, err := fetchFileWithRetry(ctx, key)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("Failed to fetch %s: %v", key, err))
			continue
		}

		filePath := fmt.Sprintf("%s/%s", uploadDir, upload.Filename)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			errorsList = append(errorsList, fmt.Sprintf("Failed to write %s: %v", filePath, err))
			continue
		}

		utils.Logger.WithFields(map[string]interface{}{
			"job_id":  job.ID,
			"user_id": user.ID,
			"file":    upload.Filename,
		}).Infof("File %s downloaded successfully", upload.Filename)
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

func fetchFileWithRetry(ctx context.Context, path string) ([]byte, error) {
	var data []byte
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		data, err = services.R2.GetObject(ctx, path)
		if err == nil {
			return data, nil
		}

		utils.Logger.WithError(err).WithField("path", path).
			Warnf("Retry %d: failed to fetch R2 object", attempt)

		time.Sleep(time.Second * time.Duration(attempt))
	}

	return nil, fmt.Errorf("failed to fetch object %s after 3 attempts: %w", path, err)
}
