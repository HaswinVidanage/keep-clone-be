package entities

import "time"

type TokenClaim struct {
	UserID    int       `json:"id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expire_at"`
	IssuedAt  time.Time `json:"issued_at"`
}
