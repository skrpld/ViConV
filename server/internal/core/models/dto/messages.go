package dto

import (
	"github.com/google/uuid"
	"github.com/skrpld/NearBeee/internal/core/models/entities"
)

type CreateMessageRequest struct {
	MessageId string    `json:"-"`
	PostId    string    `json:"post_id"`
	UserId    uuid.UUID `json:"-"`
	Content   string    `json:"content"`
}

type CreateMessageResponse struct {
	Message *entities.Message `json:"message"`
}

type GetMessageByMessageIdRequest struct {
	MessageId string `json:"-"`
}

type GetMessageByMessageIdResponse struct {
	Message *entities.Message `json:"message"`
}

type GetMessageByUserIdRequest struct {
	UserId uuid.UUID `json:"-"`
	Count  int64     `json:"count"`
}
type GetMessageByUserIdResponse struct {
	Messages []*entities.Message `json:"messages"`
}

type GetMessagesByPostIdRequest struct {
	PostId string `json:"post_id"`
	Count  int64  `json:"count"`
}
type GetMessagesByPostIdResponse struct {
	Messages []*entities.Message `json:"messages"`
}

type UpdateMessageByIdRequest struct {
	MessageId string    `json:"-"`
	UserId    uuid.UUID `json:"-"`
	Content   string    `json:"content"`
}
type UpdateMessageByIdResponse struct {
	Message *entities.Message `json:"message"`
}

type DeleteMessageByIdRequest struct {
	MessageId string    `json:"-"`
	UserId    uuid.UUID `json:"-"`
}

type DeleteMessageByIdResponse struct {
	Success bool `json:"success"`
}
