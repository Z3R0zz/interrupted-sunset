package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"interrupted-export/src/database"
	"interrupted-export/src/mail"
	"interrupted-export/src/utils"
	"math/big"
	"time"
)

type OTP struct {
	ID        uint64
	UserID    uint
	Email     string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func (o *OTP) Create(ctx context.Context) error {
	o.Code = o.GenerateCode()

	tx, err := database.DB.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to begin transaction")
		return err
	}

	defer func() {
		if err != nil {
			utils.Logger.WithError(err).Error("Rolling back transaction")
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO otp (user_id, email, code, expires_at)
		VALUES (?, ?, ?, ?)
	`, o.UserID, o.Email, o.Code, o.ExpiresAt)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to insert OTP record")
		return err
	}

	sender, err := mail.NewEmailSender()
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to initialize email sender")
		return err
	}

	emailBody := fmt.Sprintf("Your OTP code is: %s", o.Code)
	err = sender.SendEmail(o.Email, "Your OTP Code", []byte(emailBody))
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to send OTP email")
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

func (o *OTP) GenerateCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%06d", n)
}
