package user

import "time"

type RefreshToken struct {
	Token     string
	UserID    int64
	ExpiresAt time.Time
	Revoked   bool
}
