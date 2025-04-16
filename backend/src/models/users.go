package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"interrupted-export/src/config"
	"interrupted-export/src/database"
	"interrupted-export/src/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              uint
	Username        string
	Email           string
	EmailVerifiedAt *time.Time
	Password        string
}

func (u *User) AttemptLogin(ctx context.Context) (string, error) {
	var hashedPassword string

	row := database.DB.QueryRowContext(ctx, "SELECT id, username, email, password FROM users WHERE username = ?", u.Username)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("invalid credentials")
		}

		utils.Logger.WithError(err).WithFields(logrus.Fields{
			"username": u.Username,
			"email":    u.Email,
		}).Error("error querying user")

		return "", fmt.Errorf("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	expiration := time.Now().Add(30 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

func (u *User) Get(ctx context.Context) error {
	var verifiedAtRaw []byte

	row := database.DB.QueryRowContext(ctx, "SELECT id, username, email, email_verified_at FROM users WHERE id = ?", u.ID)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &verifiedAtRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}

		utils.Logger.WithError(err).WithFields(logrus.Fields{
			"user_id": u.ID,
		}).Error("error querying user")

		return fmt.Errorf("internal server error")
	}

	if len(verifiedAtRaw) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(verifiedAtRaw))
		if err != nil {
			utils.Logger.WithError(err).Error("failed to parse email_verified_at")
		} else {
			u.EmailVerifiedAt = &parsedTime
		}
	}

	return nil
}
