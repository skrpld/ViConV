package dto

import (
	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
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

type GetPostsByLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Count     int64   `json:"count"`
	Radius    float64 `json:"radius"`
}

type GetPostsByLocationResponse struct {
	Posts []*entities.Post `json:"posts"`
}

type GetPostByPostIdRequest struct {
	PostId string
	UserId uuid.UUID
}

type GetPostByPostIdResponse struct {
	Post *entities.Post `json:"post"`
}

type UpdatePostByIdRequest struct {
	PostId  string
	UserId  uuid.UUID
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostByIdResponse struct {
	Post *entities.Post `json:"post"`
}

type DeletePostByIdRequest struct {
	PostId string
	UserId uuid.UUID
}

type DeletePostResponse struct {
	Message string `json:"message"`
}
