package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"interrupted-export/src/database"
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

	_, err := database.DB.ExecContext(ctx, `
		INSERT INTO otp (user_id, email, code, expires_at)
		VALUES (?, ?, ?, ?)
	`, o.UserID, o.Email, o.Code, o.ExpiresAt)

	if err != nil {
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
