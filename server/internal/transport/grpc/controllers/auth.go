package controllers

import (
	"context"
	"viconv/internal/models/dto"
	"viconv/pkg/api/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	RegistrateUser(rows *dto.RegistrateUserRequest) (*dto.RegistrateUserResponse, error)
	LoginUser(rows *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	RefreshUserToken(rows *dto.RefreshUserTokenRequest) (*dto.RefreshUserTokenResponse, error)
}

type AuthController struct {
	auth.UnimplementedAuthServiceServer
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) RegistrateUser(ctx context.Context, req *auth.RegistrateUserRequest) (*auth.RegistrateUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	response, err := c.authService.RegistrateUser(&dto.RegistrateUserRequest{Email: email, Password: password})
	if err != nil {
		return nil, err
	}

	err = grpc.SetHeader(ctx, metadata.Pairs("Authorization", "Bearer "+response.AccessToken))
	if err != nil {
		return nil, err
	}

	return &auth.RegistrateUserResponse{Message: response.Message}, nil
}

func (c *AuthController) LoginUser(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	response, err := c.authService.LoginUser(&dto.LoginUserRequest{Email: email, Password: password})
	if err != nil {
		return nil, err
	}

	err = grpc.SetHeader(ctx, metadata.Pairs("Authorization", "Bearer "+response.AccessToken))
	if err != nil {
		return nil, err
	}
	return &auth.LoginUserResponse{Message: response.Message, RefreshToken: response.RefreshToken}, nil
}

func (c *AuthController) RefreshUserToken(ctx context.Context, req *auth.RefreshUserTokenRequest) (*auth.RefreshUserTokenResponse, error) {
	refreshToken := req.GetRefreshToken()

	response, err := c.authService.RefreshUserToken(&dto.RefreshUserTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return nil, err
	}

	err = grpc.SetHeader(ctx, metadata.Pairs("Authorization", "Bearer "+response.AccessToken))
	if err != nil {
		return nil, err
	}
	return &auth.RefreshUserTokenResponse{Message: response.Message}, nil
}
