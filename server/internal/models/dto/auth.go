package dto

import "viconv/internal/models/entities"

type RegistrateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegistrateUserResponse struct {
	Message      string `json:"message"`
	RefreshToken string
	AccessToken  string
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string
}

type RefreshUserTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshUserTokenResponse struct {
	Message      string `json:"message"`
	RefreshToken string
	AccessToken  string
}

type AuthorizeUserRequest struct {
	AccessToken string `json:"access_token"`
}

type AuthorizeUserResponse struct {
	User *entities.User
}
