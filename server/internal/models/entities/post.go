package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostId         uuid.UUID
	UserId         uuid.UUID
	Title          string
	Content        string
	IdempotencyKey string
	Latitude       float64
	Longitude      float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
