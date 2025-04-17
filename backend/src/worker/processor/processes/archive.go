package processes

import (
	"context"
	"fmt"
	"interrupted-export/src/mail"
	"interrupted-export/src/models"
	"interrupted-export/src/services"
	"interrupted-export/src/utils"
	"os"
	"path/filepath"
	"time"
)

func ProcessArchiveFile(filePath string, user *models.User) error {
	const maxSize = 20 * 1024 * 1024 // 20MB

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("checking archive size: %w", err)
	}

	if info.Size() <= maxSize {
		sender, err := mail.NewEmailSender()
		if err != nil {
			return fmt.Errorf("initializing email sender: %w", err)
		}

		if err := sender.SendArchive(user.Email, filePath); err != nil {
			utils.Logger.WithError(err).WithField("user_id", user.ID).Error("Failed to send export archive")
			return fmt.Errorf("failed to send archive: %w", err)
		}
	} else {
		url, err := uploadToR2AndPresign(context.Background(), filePath)
		if err != nil {
			return fmt.Errorf("uploading large archive: %w", err)
		}

		sender, err := mail.NewEmailSender()
		if err != nil {
			return fmt.Errorf("initializing email sender: %w", err)
		}

		body := fmt.Sprintf(`Your archive was too large to send via email.

			You can securely download it from the link below (valid for 24 hours):

			%s

			Thanks for using interrupted.me! 🖤
		`, url)

		if err := sender.SendEmail(user.Email, "Your interrupted.me export is ready", []byte(body)); err != nil {
			return fmt.Errorf("failed to send archive link: %w", err)
		}

	}

	if err := os.Remove(filePath); err != nil {
		utils.Logger.WithError(err).WithField("file_path", filePath).
			Error("Failed to remove archive file after sending")
		return fmt.Errorf("failed to remove archive file: %w", err)
	}

	return nil
}

func uploadToR2AndPresign(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	key := filepath.Base(filePath)

	if err := services.R2.UploadFile(ctx, key, file); err != nil {
		return "", fmt.Errorf("uploading to R2: %w", err)
	}

	url, err := services.R2.GeneratePresignedURL(ctx, key, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("generating R2 presigned URL: %w", err)
	}

	return url, nil
}
