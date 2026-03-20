package dto

import (
	"github.com/skrpld/NearBeee/internal/core/models/entities"
)

type RegistrateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegistrateUserResponse struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"-"`
	AccessToken  string `json:"-"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponse struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"-"`
}

type RefreshUserTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshUserTokenResponse struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"-"`
	AccessToken  string `json:"-"`
}

type AuthorizeUserRequest struct {
	AccessToken string
}

type AuthorizeUserResponse struct {
	User *entities.User
}
