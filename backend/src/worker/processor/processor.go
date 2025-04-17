package processor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"interrupted-export/src/models"
	"interrupted-export/src/utils"
	"interrupted-export/src/worker/processor/processes"
	"io"
	"os"
	"path/filepath"
)

func ProcessExportJob(job *models.Queue, user *models.User) error {
	utils.Logger.Infof("Processing export job for user: %s (user_id=%d)", user.Username, user.ID)
	ctx := context.Background()

	exportDir := fmt.Sprintf("tmp/export_user_%d", job.UserID)
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("failed to create export dir: %w", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(exportDir); err != nil {
			utils.Logger.WithError(err).WithField("path", exportDir).
				Warn("Failed to clean up temporary export directory")
		}
	}

	defer cleanup()

	if err := processes.ProcessShorteners(job, user, exportDir, ctx); err != nil {
		return fmt.Errorf("processing shorteners: %w", err)
	}

	if err := processes.ProcessPastes(job, user, exportDir, ctx); err != nil {
		return fmt.Errorf("processing pastes: %w", err)
	}

	if err := processes.ProcessUploads(job, user, exportDir, ctx); err != nil {
		return fmt.Errorf("processing uploads: %w", err)
	}

	archivePath := fmt.Sprintf("tmp/export_user_%d.tar.gz", job.UserID)
	output, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("creating tar file: %w", err)
	}
	defer output.Close()

	if err := createArchiveFromDir(exportDir, output); err != nil {
		_ = os.Remove(archivePath)
		return fmt.Errorf("creating archive: %w", err)
	}

	if err := processes.ProcessArchiveFile(archivePath, user); err != nil {
		_ = os.Remove(archivePath)
		return fmt.Errorf("processing archive file: %w", err)
	}

	return nil
}

func createArchiveFromDir(basePath string, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(basePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(basePath, file)
		if err != nil {
			return err
		}

		return addToArchive(tw, file, relPath)
	})
}

func addToArchive(tw *tar.Writer, fullPath, archivePath string) error {
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	header.Name = archivePath
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}
