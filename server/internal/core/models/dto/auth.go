package dto

import (
	"github.com/skrpld/NearBeee/internal/core/models/entities"
)

type RegistrateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegistrateUserResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"-"`
	AccessToken  string `json:"-"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"-"`
}

type RefreshUserTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshUserTokenResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"-"`
	AccessToken  string `json:"-"`
}

type AuthorizeUserRequest struct {
	AccessToken string `json:"access_token"`
}

type AuthorizeUserResponse struct {
	User *entities.User
}
