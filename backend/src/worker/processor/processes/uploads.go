package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/services"
	"interrupted-export/src/utils"
	"math/rand"
	"os"
	"strings"
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

func fetchFileWithRetry(parent context.Context, path string) ([]byte, error) {
	const maxAttempts = 3
	const perReqTimeout = 15 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(parent, perReqTimeout)
		data, err := services.R2.GetObject(attemptCtx, path)
		cancel()

		if err == nil {
			return data, nil
		}

		// Ghetto fix but fuck it we ball
		if strings.TrimSpace(err.Error()) == "operation error S3: GetObject, https response error StatusCode: 404, RequestID: , HostID: , NoSuchKey:" {
			utils.Logger.WithField("path", path).Debugf("file not found: %s", path)
			return nil, nil
		}

		lastErr = err

		utils.Logger.
			WithError(err).
			WithField("path", path).
			Warnf("attempt %d/%d failed", attempt, maxAttempts)

		backoff := time.Duration(1<<uint(attempt-1)) * time.Second
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}
		select {
		case <-parent.Done():
			return nil, parent.Err()
		case <-time.After(backoff + time.Duration(rand.Intn(500))*time.Millisecond):
		}
	}
	return nil, fmt.Errorf("failed to fetch %s after %d attempts: %w", path, maxAttempts, lastErr)
}
