package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId                 uuid.UUID
	Email                  string
	PasswordHash           string
	RefreshToken           string
	RefreshTokenExpiryTime time.Time
}
