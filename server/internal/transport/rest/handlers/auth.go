package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/models/dto"
)

type AuthService interface {
	RegistrateUser(rows *dto.RegistrateUserRequest) (*dto.RegistrateUserResponse, error)
	LoginUser(rows *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	RefreshUserToken(rows *dto.RefreshUserTokenRequest) (*dto.RefreshUserTokenResponse, error)
}

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) RegistrateUserHandler(r *http.Request) (any, error) {
	var request dto.RegistrateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.authService.RegistrateUser(&request)
}

func (c *AuthController) LoginUserHandler(r *http.Request) (any, error) {
	var request dto.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.authService.LoginUser(&request)
}

func (c *AuthController) RefreshUserTokenHandler(r *http.Request) (any, error) {
	var request dto.RefreshUserTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.authService.RefreshUserToken(&request)
}
