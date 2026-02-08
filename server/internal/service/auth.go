package service

import (
	"time"
	"viconv/internal/models/dto"
	"viconv/internal/models/entities"
	"viconv/pkg/consts/errors"
	"viconv/pkg/utils/hash"
	"viconv/pkg/utils/jwt"
	"viconv/pkg/utils/mail"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateRefreshTokenByUserId(userId uuid.UUID, refreshToken string, refreshTokenExpiryTime time.Time) error
	GetUserById(userId uuid.UUID) (*entities.User, error)
}

type AuthService struct {
	repo   AuthRepository
	secret string
}

func NewAuthService(repo AuthRepository, secret string) *AuthService {
	return &AuthService{repo, secret}
}

func (s *AuthService) RegistrateUser(rows *dto.RegistrateUserRequest) (*dto.RegistrateUserResponse, error) {
	if !mail.IsEmailValid(rows.Email) {
		return nil, errors.ErrInvalidEmail
	}

	hashPassword, err := hash.HashString(rows.Password)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshTokenExpiryDuration, err := jwt.NewRefreshToken(rows.Email, s.secret)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiryTime := time.Now().Add(refreshTokenExpiryDuration)

	newUser, err := s.repo.CreateUser(rows.Email, hashPassword, refreshToken, refreshTokenExpiryTime)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.NewAccessToken(newUser.UserId.String(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.RegistrateUserResponse{
		Message:      newUser.UserId.String(),
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return response, nil
}

func (s *AuthService) LoginUser(rows *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	if !mail.IsEmailValid(rows.Email) {
		return nil, errors.ErrInvalidEmail
	}

	user, err := s.repo.GetUserByEmail(rows.Email)
	if err != nil {
		return nil, errors.ErrInvalidEmail
	}

	err = hash.CompareHashAndPassword(user.PasswordHash, rows.Password)
	if err != nil {
		return nil, errors.ErrInvalidPassword
	}

	refreshToken, refreshTokenExpiryDuration, err := jwt.NewRefreshToken(rows.Email, s.secret)
	if err != nil {
		return nil, err
	}
	refreshTokenExpiryTime := time.Now().Add(refreshTokenExpiryDuration)

	err = s.repo.UpdateRefreshTokenByUserId(user.UserId, refreshToken, refreshTokenExpiryTime)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.NewAccessToken(user.UserId.String(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.LoginUserResponse{
		Message:      user.UserId.String(),
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return response, nil
}

func (s *AuthService) RefreshUserToken(rows *dto.RefreshUserTokenRequest) (*dto.RefreshUserTokenResponse, error) {
	tokenClaims, err := jwt.ValidateToken(rows.RefreshToken, s.secret)
	if err != nil {
		return nil, err
	}

	email, err := tokenClaims.GetSubject()
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if rows.RefreshToken != user.RefreshToken {
		return nil, errors.ErrInvalidToken
	}

	accessToken, err := jwt.NewAccessToken(user.UserId.String(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.RefreshUserTokenResponse{
		Message:      user.UserId.String(),
		RefreshToken: accessToken,
		AccessToken:  accessToken,
	}

	return response, nil
}

func (s *AuthService) AuthorizeUser(rows *dto.AuthorizeUserRequest) (*dto.AuthorizeUserResponse, error) {
	tokenClaims, err := jwt.ValidateToken(rows.AccessToken, s.secret)
	if err != nil {
		return nil, err
	}

	sub, err := tokenClaims.GetSubject()
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	response := &dto.AuthorizeUserResponse{
		User: user,
	}

	return response, nil
}
