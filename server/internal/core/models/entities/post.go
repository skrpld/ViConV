package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostId         uuid.UUID `json:"post_id"`
	UserId         uuid.UUID `json:"user_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	IdempotencyKey string    `json:"idempotency_key"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
