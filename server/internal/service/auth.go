package service

import (
	stderr "errors"
	"time"
	"viconv/internal/models/dto"
	"viconv/internal/models/entities"
	"viconv/pkg/consts/errors"
	"viconv/pkg/utils/hash"
	"viconv/pkg/utils/jwt"
	"viconv/pkg/utils/mail"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateRefreshTokenByUserId(userId bson.ObjectID, refreshToken string, refreshTokenExpiryTime time.Time) error
	GetUserById(id bson.ObjectID) (*entities.User, error)
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

	user, err := s.repo.GetUserByEmail(rows.Email)
	if err != nil {
		if !stderr.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	}
	if user != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	hashPassword, err := hash.HashString(rows.Password)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewRefreshToken(rows.Email, s.secret)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiryTime := time.Now().Add(time.Hour * 24 * 7)

	newUser, err := s.repo.CreateUser(rows.Email, hashPassword, refreshToken, refreshTokenExpiryTime)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.NewAccessToken(newUser.UserId.Hex(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.RegistrateUserResponse{
		Message:      newUser.UserId.Hex(),
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

	refreshToken, err := jwt.NewRefreshToken(rows.Email, s.secret)
	if err != nil {
		return nil, err
	}
	refreshTokenExpiryTime := time.Now().Add(time.Hour * 24 * 7)

	err = s.repo.UpdateRefreshTokenByUserId(user.UserId, refreshToken, refreshTokenExpiryTime)
	if err != nil {
		return nil, err
	}

	accessToken, err := jwt.NewAccessToken(user.UserId.Hex(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.LoginUserResponse{
		Message:      user.UserId.Hex(),
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

	accessToken, err := jwt.NewAccessToken(user.UserId.Hex(), s.secret)
	if err != nil {
		return nil, err
	}

	response := &dto.RefreshUserTokenResponse{
		Message:      user.UserId.Hex(),
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

	id, err := tokenClaims.GetSubject()
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	userId, err := bson.ObjectIDFromHex(id)
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
