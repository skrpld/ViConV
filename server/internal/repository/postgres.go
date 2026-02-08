package repository

import (
	"context"
	"database/sql"
	stderr "errors"
	"fmt"
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
		if ok && pgErr.Code == "23505" { // 23505 - unique_violation
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

func (r *ViconvRepository) CreatePost(userId uuid.UUID, title, content, idempotencyKey string, latitude, longitude float64) (*entities.Post, error) {
	var post entities.Post

	query := `INSERT INTO posts (user_id, title, content, idempotency_key, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	err := r.postgresDB.QueryRow(query, userId, title, content, idempotencyKey, latitude, longitude).
		Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.IdempotencyKey, &post.Latitude, &post.Longitude, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" { // 23505 - unique_violation
			return nil, errors.ErrIdempotencyKeyAlreadyExists
		}
		return nil, err
	}

	return &post, nil
}

func (r *ViconvRepository) GetPostsByUserId(userId uuid.UUID, count int64) ([]*entities.Post, error) {
	var posts []*entities.Post

	query := `SELECT * FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.postgresDB.Query(query, userId, count)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrExpiredToken
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var post entities.Post

		err = rows.Scan(&post.PostId, &post.UserId,
			&post.Title, &post.Content,
			&post.IdempotencyKey, &post.Latitude,
			&post.Longitude, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *ViconvRepository) GetPostsByLocation(latitude, longitude, radius float64, count int64) ([]*entities.Post, error) {
	var posts []*entities.Post

	query := `SELECT * FROM posts
		WHERE haversine_distance($1, $2, latitude, longitude) <= $3
		ORDER BY haversine_distance($1, $2, latitude, longitude) 
		LIMIT $4`

	rows, err := r.postgresDB.Query(query, latitude, longitude, radius, count)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error debug") //TODO: err no posts in radius + потестить когда ошибка срабатывает
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var post entities.Post

		err = rows.Scan(&post.PostId, &post.UserId,
			&post.Title, &post.Content,
			&post.IdempotencyKey, &post.Latitude,
			&post.Longitude, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *ViconvRepository) GetPostById(postId, userId uuid.UUID) (*entities.Post, error) {
	var post entities.Post

	query := `SELECT * FROM posts WHERE post_id = $1 AND user_id = $2`
	err := r.postgresDB.QueryRow(query, postId, userId).
		Scan(&post.PostId, &post.UserId,
			&post.Title, &post.Content,
			&post.IdempotencyKey, &post.Latitude,
			&post.Longitude, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidPostId
		}
		return nil, err
	}

	return &post, nil
}

func (r *ViconvRepository) UpdatePostById(post *entities.Post) (*entities.Post, error) {
	var newPost entities.Post

	query := `UPDATE posts SET title = $1, content = $2 WHERE post_id = $3 AND user_id = $4 RETURNING *`

	err := r.postgresDB.QueryRow(query, post.Title, post.Content, post.PostId, post.UserId).
		Scan(&newPost.PostId, &newPost.UserId, &newPost.Title,
			&newPost.Content, &newPost.IdempotencyKey,
			&newPost.Latitude, &newPost.Longitude,
			&newPost.CreatedAt, &newPost.UpdatedAt)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrInvalidPostId
		}
		return nil, err
	}
	
	return &newPost, nil
}

func (r *ViconvRepository) DeletePostById(postId, userId uuid.UUID) error {
	query := `DELETE FROM posts WHERE post_id = $1 AND user_id = $2`

	_, err := r.postgresDB.Exec(query, postId, userId)
	if err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return errors.ErrInvalidPostId
		}
		return err
	}

	return nil
}
