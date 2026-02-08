package dto

import (
	"viconv/internal/models/entities"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	UserId         uuid.UUID
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	IdempotencyKey string  `json:"idempotency_key"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}

type CreatePostResponse struct {
	Message string `json:"message"`
}

type GetPostsByUserIdRequest struct {
	UserId uuid.UUID
	Count  int64 `json:"count"`
}

type GetPostsByUserIdResponse struct {
	Posts []*entities.Post `json:"posts"`
}

type GetPostsByLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Count     int64   `json:"count"`
	Radius    float64 `json:"radius"`
}

type GetPostsByLocationResponse struct {
	Posts []*entities.Post `json:"posts"`
}

type GetPostByIdRequest struct {
	PostId string `json:"post_id"`
	UserId uuid.UUID
}

type GetPostByIdResponse struct {
	Post *entities.Post `json:"post"`
}

type UpdatePostByIdRequest struct {
	PostId  string `json:"post_id"`
	UserId  uuid.UUID
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostByIdResponse struct {
	Post *entities.Post `json:"post"`
}

type DeletePostByIdRequest struct {
	PostId string `json:"post_id"`
	UserId uuid.UUID
}

type DeletePostResponse struct {
	Message string `json:"message"`
}
