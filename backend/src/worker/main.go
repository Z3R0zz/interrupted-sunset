package main

import (
	"context"
	"database/sql"
	"errors"
	"interrupted-export/src/config"
	"interrupted-export/src/database"
	"interrupted-export/src/models"
	"interrupted-export/src/services"
	"interrupted-export/src/utils"
	"interrupted-export/src/worker/processor"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		if errors.Is(err, sql.ErrNoRows) {
			_ = tx.Commit()

			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			utils.Logger.WithError(err).Error("Failed to fetch job")
			_ = tx.Rollback()
			time.Sleep(2 * time.Second)
			continue
		}

		user := models.User{
			ID: queue.UserID,
		}

		if err := user.Get(ctx); err != nil {
			utils.Logger.WithError(err).WithField("user_id", queue.UserID).Error("Failed to get user")
			_ = tx.Rollback()
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

		if err := processor.ProcessExportJob(queue, &user); err != nil {
			utils.Logger.WithError(err).WithField("job_id", queue.ID).Error("Failed to process job")
			queue.MarkFailed(err.Error())
			continue
		}

		queue.MarkDone()
	}
}

func main() {
	utils.Logger.Info("Starting worker...")

	config.Load()

	database.Connect(os.Getenv("DATABASE_URL"))

	if err := services.ConnectR2(); err != nil {
		utils.Logger.WithError(err).Error("Failed to connect to R2")
		os.Exit(1)
	}

	workerLoop()
}
