package controllers

import (
	"context"
	"encoding/json"

	"github.com/skrpld/NearBeee/internal/models/dto"
	"github.com/skrpld/NearBeee/internal/models/entities"
	"github.com/skrpld/NearBeee/pkg/api/posts"
	"github.com/skrpld/NearBeee/pkg/consts"
	"github.com/skrpld/NearBeee/pkg/consts/errors"
)

type PostsService interface {
	CreatePost(rows *dto.CreatePostRequest) (*dto.CreatePostResponse, error)
	GetPostsByUserId(rows *dto.GetPostsByUserIdRequest) (*dto.GetPostsByUserIdResponse, error)
	GetPostsByLocation(rows *dto.GetPostsByLocation) (*dto.GetPostsByLocationResponse, error)
	GetPostById(rows *dto.GetPostByIdRequest) (*dto.GetPostByIdResponse, error)
	UpdatePostById(rows *dto.UpdatePostByIdRequest) (*dto.UpdatePostByIdResponse, error)
	DeletePostById(rows *dto.DeletePostByIdRequest) (*dto.DeletePostResponse, error)
}

type PostsController struct {
	posts.UnimplementedPostsServiceServer
	postsService PostsService
}

func NewPostsController(postsService PostsService) *PostsController {
	return &PostsController{postsService: postsService}
}

func getUserFromCtx(ctx context.Context) (*entities.User, error) {
	ctxUser := ctx.Value(consts.CtxUserKey)
	user, ok := ctxUser.(*entities.User)
	if !ok {
		return nil, errors.ErrNoPermissions
	}
	return user, nil
}

func (s *PostsController) CreatePost(ctx context.Context, req *posts.CreatePostRequest) (*posts.CreatePostResponse, error) {
	title := req.GetTitle()
	content := req.GetContent()
	idempotencyKey := req.GetIdempotencyKey()
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()

	user, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.postsService.CreatePost(&dto.CreatePostRequest{
		UserId:         user.UserId,
		Title:          title,
		Content:        content,
		IdempotencyKey: idempotencyKey,
		Latitude:       latitude,
		Longitude:      longitude})
	if err != nil {
		return nil, err
	}

	return &posts.CreatePostResponse{Message: response.Message}, nil
}

func (s *PostsController) GetPostById(ctx context.Context, req *posts.GetPostByIdRequest) (*posts.GetPostByIdResponse, error) {
	postId := req.GetPostId()
	user, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.postsService.GetPostById(&dto.GetPostByIdRequest{PostId: postId, UserId: user.UserId})
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &posts.GetPostByIdResponse{Message: string(data)}, nil
}
func (s *PostsController) GetUserPosts(ctx context.Context, req *posts.GetUserPostsRequest) (*posts.GetUserPostsResponse, error) {
	count := req.GetCount()
	user, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//TODO: реализовать количество (0 = все посты)

	response, err := s.postsService.GetPostsByUserId(&dto.GetPostsByUserIdRequest{UserId: user.UserId, Count: count})
	if err != nil {
		return nil, err
	}

	//временный костыль - а может и нет
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &posts.GetUserPostsResponse{Message: string(data)}, nil
}

func (s *PostsController) GetPostsByLocation(ctx context.Context, req *posts.GetPostsByLocationRequest) (*posts.GetPostsByLocationResponse, error) {
	count := req.GetCount()
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()
	radius := req.GetRadius()

	//TODO: реализовать количество (0 = все посты)

	response, err := s.postsService.GetPostsByLocation(&dto.GetPostsByLocation{
		Latitude: latitude, Longitude: longitude, Count: count, Radius: radius})
	if err != nil {
		return nil, err
	}

	//временный костыль - а может и нет
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &posts.GetPostsByLocationResponse{Message: string(data)}, nil
}

func (s *PostsController) UpdatePostById(ctx context.Context, req *posts.UpdatePostByIdRequest) (*posts.UpdatePostByIdResponse, error) {
	postId := req.GetPostId()
	title := req.GetTitle()
	content := req.GetContent()

	user, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.postsService.UpdatePostById(&dto.UpdatePostByIdRequest{
		PostId:  postId,
		UserId:  user.UserId,
		Title:   title,
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &posts.UpdatePostByIdResponse{Message: string(data)}, nil
}

func (s *PostsController) DeletePostById(ctx context.Context, req *posts.DeletePostByIdRequest) (*posts.DeletePostByIdResponse, error) {
	postId := req.GetPostId()
	user, err := getUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.postsService.DeletePostById(&dto.DeletePostByIdRequest{
		PostId: postId,
		UserId: user.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &posts.DeletePostByIdResponse{Message: response.Message}, nil
}
