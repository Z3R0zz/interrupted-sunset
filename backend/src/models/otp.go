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

	err = o.DeleteOld(ctx)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to delete old OTP record")
		return err
	}

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

func (o *OTP) DeleteOld(ctx context.Context) error {
	_, err := database.DB.ExecContext(ctx, `
		DELETE FROM otp
		WHERE user_id = ?
	`, o.UserID)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to delete old OTP record")
		return err
	}
	return nil
}

func (o *OTP) Verify(ctx context.Context) error {
	valid, err := o.Valid(ctx)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to verify OTP")
		return fmt.Errorf("internal server error")
	}

	if !valid {
		return fmt.Errorf("invalid OTP code")
	}

	if err := o.Get(ctx); err != nil {
		utils.Logger.WithError(err).Error("Failed to get OTP record")
		return fmt.Errorf("internal server error")
	}

	_, err = database.DB.ExecContext(ctx, `
		UPDATE users
		SET email = ?, email_verified_at = ?
		WHERE id = ?
	`, o.Email, time.Now().UTC(), o.UserID)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to update user email_verified_at")
		return fmt.Errorf("internal server error")
	}

	_, err = database.DB.ExecContext(ctx, `
		DELETE FROM otp
		WHERE user_id = ? AND code = ?
	`, o.UserID, o.Code)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to delete OTP record")
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (o *OTP) Get(ctx context.Context) error {
	var createdAtStr, updatedAtStr, expiresAtStr string

	err := database.DB.QueryRowContext(ctx, `
		SELECT id, user_id, email, code, created_at, updated_at, expires_at
		FROM otp
		WHERE user_id = ? AND code = ?
	`, o.UserID, o.Code).Scan(&o.ID, &o.UserID, &o.Email, &o.Code,
		&createdAtStr, &updatedAtStr, &expiresAtStr)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to get OTP record")
		return err
	}

	layout := "2006-01-02 15:04:05"
	o.CreatedAt, err = time.Parse(layout, createdAtStr)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to parse created_at")
		return err
	}
	o.UpdatedAt, err = time.Parse(layout, updatedAtStr)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to parse updated_at")
		return err
	}
	o.ExpiresAt, err = time.Parse(layout, expiresAtStr)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to parse expires_at")
		return err
	}

	return nil
}

func (o *OTP) Valid(ctx context.Context) (bool, error) {
	var count int
	err := database.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM otp
		WHERE user_id = ? AND code = ? AND expires_at > ?
	`, o.UserID, o.Code, time.Now().UTC()).Scan(&count)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to verify OTP")
		return false, err
	}

	return count > 0, nil
}

func (o *OTP) GenerateCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%06d", n)
}
