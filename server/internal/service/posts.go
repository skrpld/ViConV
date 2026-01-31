package service

import (
	"viconv/internal/models/dto"
	"viconv/internal/models/entities"
	"viconv/pkg/consts"
	"viconv/pkg/consts/errors"

	"github.com/google/uuid"
)

type PostsRepository interface {
	CreatePost(userId uuid.UUID, title, content, idempotencyKey string, latitude, longitude float64) (*entities.Post, error)
	GetPostsByUserId(userId uuid.UUID, count int64) ([]*entities.Post, error)
	GetPostsByLocation(latitude, longitude, radius float64, count int64) ([]*entities.Post, error)
	GetPostById(postId, userId uuid.UUID) (*entities.Post, error)
	UpdatePostById(post *entities.Post) (*entities.Post, error)
	DeletePostById(postId, userId uuid.UUID) error
}

type PostsService struct {
	repo PostsRepository
}

func NewPostsService(repo PostsRepository) *PostsService {
	return &PostsService{repo: repo}
}

func (s *PostsService) CreatePost(rows *dto.CreatePostRequest) (*dto.CreatePostResponse, error) {
	post, err := s.repo.CreatePost(rows.UserId, rows.Title, rows.Content, rows.IdempotencyKey, rows.Latitude, rows.Longitude)
	if err != nil {
		return nil, err
	}

	response := dto.CreatePostResponse{
		Message: post.PostId.String(),
	}

	return &response, nil
}

func (s *PostsService) GetPostsByUserId(rows *dto.GetPostsByUserIdRequest) (*dto.GetPostsByUserIdResponse, error) {
	posts, err := s.repo.GetPostsByUserId(rows.UserId, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetPostsByUserIdResponse{
		Posts: posts,
	}

	return &response, nil
}

func (s *PostsService) GetPostsByLocation(rows *dto.GetPostsByLocation) (*dto.GetPostsByLocationResponse, error) {
	posts, err := s.repo.GetPostsByLocation(rows.Latitude, rows.Longitude, rows.Radius, rows.Count)
	if err != nil {
		return nil, err
	}

	response := dto.GetPostsByLocationResponse{
		Posts: posts,
	}

	return &response, nil
}

func (s *PostsService) GetPostById(rows *dto.GetPostByIdRequest) (*dto.GetPostByIdResponse, error) {
	postId, err := uuid.Parse(rows.PostId)
	if err != nil {
		return nil, errors.ErrInvalidPostId
	}

	post, err := s.repo.GetPostById(postId, rows.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.GetPostByIdResponse{
		Post: post,
	}

	return &response, nil
}

// UpdatePostById - работает сейчас, что в запросе отправляются только данные, которые изменяются
// Если что-то не изменяется, то поле = "no_change"
// Можно переделать, что клиент изначально берет пост, а потом отправляет все поля,
// Если они не изменились то они просто обновятся на те же значения
func (s *PostsService) UpdatePostById(rows *dto.UpdatePostByIdRequest) (*dto.UpdatePostByIdResponse, error) {
	postId, err := uuid.Parse(rows.PostId)
	if err != nil {
		return nil, errors.ErrInvalidPostId
	}

	post, err := s.repo.GetPostById(postId, rows.UserId)
	if err != nil {
		return nil, err
	}

	if rows.Title != consts.NoChangeKey {
		post.Title = rows.Title
	}
	if rows.Content != consts.NoChangeKey {
		post.Content = rows.Content
	}

	post, err = s.repo.UpdatePostById(post)
	if err != nil {
		return nil, err
	}

	response := dto.UpdatePostByIdResponse{
		Post: post,
	}

	return &response, nil
}

func (s *PostsService) DeletePostById(rows *dto.DeletePostByIdRequest) (*dto.DeletePostResponse, error) {
	postId, err := uuid.Parse(rows.PostId)
	if err != nil {
		return nil, errors.ErrInvalidPostId
	}

	err = s.repo.DeletePostById(postId, rows.UserId)
	if err != nil {
		return nil, err
	}

	response := dto.DeletePostResponse{
		Message: rows.PostId,
	}

	return &response, nil
}
