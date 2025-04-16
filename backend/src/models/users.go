package models

import "time"

type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username" validate:"max=45,min=3,required"`
	Email           string    `json:"email" validate:"email,max=255,required"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	Password        string    `json:"password" validate:"max=255,min=8,required"`
}
