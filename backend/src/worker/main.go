package main

import (
	"context"
	"database/sql"
	"interrupted-export/src/database"
	"interrupted-export/src/models"
	"interrupted-export/src/utils"
	"os"
	"time"
)

func workerLoop() {
	for {
		ctx := context.Background()

		tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to begin transaction")
			time.Sleep(2 * time.Second)
			continue
		}

		queue, err := models.FetchJob(ctx)
		if err == sql.ErrNoRows {
			_ = tx.Commit()
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			utils.Logger.WithError(err).Error("Failed to fetch job")
			_ = tx.Rollback()
			time.Sleep(2 * time.Second)
			continue
		}

		if err := queue.MarkProcessing(ctx); err != nil {
			utils.Logger.WithError(err).WithField("job_id", queue.ID).Error("Failed to mark job as processing")
			_ = tx.Rollback()
			continue
		}

		if err := tx.Commit(); err != nil {
			utils.Logger.WithError(err).Error("Failed to commit transaction")
			continue
		}

		utils.Logger.WithFields(map[string]interface{}{
			"job_id":  queue.ID,
			"user_id": queue.UserID,
		}).Info("Processing job")

		// We doin crazy shit here
		// models.MarkFailed(job.ID, "some error")
		// models.MarkDone(job.ID)
	}
}

func main() {
	utils.Logger.Info("Starting worker...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		utils.Logger.Fatal("DATABASE_URL is not set")
	}

	database.Connect(dsn)
	workerLoop()
}
