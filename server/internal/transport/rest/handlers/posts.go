package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
	"github.com/skrpld/NearBeee/pkg/errors"
)

type PostsService interface {
	CreatePost(rows *dto.CreatePostRequest) (*dto.CreatePostResponse, error)
	GetPostsByUserId(rows *dto.GetPostsByUserIdRequest) (*dto.GetPostsByUserIdResponse, error)
	GetPostsByLocation(rows *dto.GetPostsByLocationRequest) (*dto.GetPostsByLocationResponse, error)
	GetPostByPostId(rows *dto.GetPostByPostIdRequest) (*dto.GetPostByPostIdResponse, error)
	UpdatePostById(rows *dto.UpdatePostByIdRequest) (*dto.UpdatePostByIdResponse, error)
	DeletePostById(rows *dto.DeletePostByIdRequest) (*dto.DeletePostResponse, error)
}
type PostsController struct {
	postsSrv PostsService
}

func NewPostsController(postsSrv PostsService) *PostsController {
	return &PostsController{postsSrv: postsSrv}
}

type formType string

const (
	userForm     formType = "user"
	locationForm formType = "location"
	postForm     formType = "post"
	nullForm     formType = ""
)

const (
	formValue = "type"
	pathValue = "post_id"
)

func (c *PostsController) CreatePostHandler(r *http.Request) (any, error) {
	var request dto.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.postsSrv.CreatePost(&request)
}

func (c *PostsController) GetPosts(r *http.Request) (any, error) {
	switch formType(r.FormValue(formValue)) {
	case userForm:
		return c.GetPostsByUserId(r)
	case locationForm:
		return c.GetPostsByLocation(r)
	case postForm, nullForm:
		return c.GetPostByPostId(r)
	default:
		return nil, errors.ErrInvalidPostType
	}
}

func (c *PostsController) GetPostsByUserId(r *http.Request) (any, error) {
	var request dto.GetPostsByUserIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.postsSrv.GetPostsByUserId(&request)
}

func (c *PostsController) GetPostsByLocation(r *http.Request) (any, error) {
	var request dto.GetPostsByLocationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	return c.postsSrv.GetPostsByLocation(&request)
}

func (c *PostsController) GetPostByPostId(r *http.Request) (any, error) {
	var request dto.GetPostByPostIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.PostId = r.PathValue(pathValue)

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}

	request.UserId = user.UserId

	return c.postsSrv.GetPostByPostId(&request)
}

func (c *PostsController) UpdatePostById(r *http.Request) (any, error) {
	var request dto.UpdatePostByIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.PostId = r.PathValue(pathValue)

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}
	request.UserId = user.UserId

	return c.postsSrv.UpdatePostById(&request)
}

func (c *PostsController) DeletePostById(r *http.Request) (any, error) {
	var request dto.DeletePostByIdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.PostId = r.PathValue(pathValue)

	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return nil, err
	}
	request.UserId = user.UserId

	return c.postsSrv.DeletePostById(&request)
}
