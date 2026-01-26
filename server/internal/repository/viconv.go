package repository

import (
	"context"
	"time"
	"viconv/internal/database/mongodb"
	"viconv/internal/models/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ViconvRepository struct {
	ctx *context.Context
	db  *mongodb.DB
}

func NewViconvRepository(ctx *context.Context, db *mongodb.DB) *ViconvRepository {
	return &ViconvRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *ViconvRepository) CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error) {
	result, err := r.db.Collection("users").InsertOne(*r.ctx, bson.D{
		{"email", email},
		{"password_hash", passwordHash},
		{"refresh_token", refreshToken},
		{"refresh_token_expiry_time", refreshTokenExpiryTime}, //TODO: refactor
	})
	if err != nil {
		return nil, err
	}

	return &entities.User{
		UserId:                 (result.InsertedID).(bson.ObjectID),
		Email:                  email,
		PasswordHash:           passwordHash,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiryTime,
	}, nil
}

func (r *ViconvRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Collection("users").FindOne(*r.ctx, bson.D{
		{"email", email}, //TODO: refactor
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ViconvRepository) GetUserById(id bson.ObjectID) (*entities.User, error) {
	var user entities.User
	err := r.db.Collection("users").FindOne(*r.ctx, bson.D{
		{"_id", id}, //TODO: refactor
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ViconvRepository) UpdateRefreshTokenByUserId(userId bson.ObjectID, refreshToken string, refreshTokenExpiryTime time.Time) error {
	_, err := r.db.Collection("users").UpdateOne(*r.ctx, bson.D{{"_id", userId}},
		bson.D{
			{
				"$set", bson.D{
					{"refresh_token", refreshToken},
					{"refresh_token_expiry_time", refreshTokenExpiryTime},
				},
			},
		}) //TODO: refactor
	if err != nil {
		return err
	}
	return nil
}
