package repository

import (
	"context"
	"database/sql"
	stderr "errors"
	"time"
	"viconv/internal/database/mongodb"
	"viconv/internal/database/postgres"
	"viconv/internal/models/entities"
	"viconv/pkg/consts/errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type ViconvRepository struct {
	ctx        *context.Context
	postgresDB *postgres.PostgresDB
	mongoDB    *mongodb.MongoDB
}

func NewViconvRepository(ctx *context.Context, postgresDB *postgres.PostgresDB, mongoDB *mongodb.MongoDB) *ViconvRepository {
	return &ViconvRepository{
		ctx:        ctx,
		postgresDB: postgresDB,
		mongoDB:    mongoDB,
	}
}

func (r *ViconvRepository) CreateUser(email, passwordHash, refreshToken string, refreshTokenExpiryTime time.Time) (*entities.User, error) {
	var user entities.User

	query := `INSERT INTO users (email, password_hash, refresh_token, refresh_token_expiry_time) VALUES ($1, $2, $3, $4) RETURNING *`

	err := r.postgresDB.QueryRow(query, email, passwordHash, refreshToken, refreshTokenExpiryTime).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			return nil, errors.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *ViconvRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	query := `SELECT * FROM users WHERE email = $1`

	err := r.postgresDB.QueryRow(query, email).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidEmail
		}
		return nil, err
	}

	return &user, nil
}

func (r *ViconvRepository) GetUserById(userId uuid.UUID) (*entities.User, error) {
	var user entities.User

	query := `SELECT * FROM users WHERE user_id = $1`

	err := r.postgresDB.QueryRow(query, userId).
		Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.RefreshTokenExpiryTime)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidEmail
		}
		return nil, err
	}

	return &user, nil
}

func (r *ViconvRepository) UpdateRefreshTokenByUserId(userId uuid.UUID, refreshToken string, refreshTokenExpiryTime time.Time) error {
	query := `UPDATE users SET refresh_token = $1, refresh_token_expiry_time = $2 WHERE user_id = $3`

	_, err := r.postgresDB.Exec(query, refreshToken, refreshTokenExpiryTime, userId)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return errors.ErrInvalidToken
		}
		return err
	}
	return nil
}
