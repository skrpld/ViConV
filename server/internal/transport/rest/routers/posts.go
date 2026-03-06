package routers

import (
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/core/service"
	"github.com/skrpld/NearBeee/internal/transport/rest/handlers"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
)

func NewPostsRouter(repo *repository.NearBeeeRepository) *http.ServeMux {
	srv := service.NewPostsService(repo)
	controller := handlers.NewPostsController(srv)
	router := http.NewServeMux()

	router.HandleFunc("POST /posts/", web.Handle(controller.CreatePostHandler))
	router.HandleFunc("GET /posts/", web.Handle(controller.GetPosts))
	router.HandleFunc("GET /posts/{post_id}", web.Handle(controller.GetPosts))
	router.HandleFunc("PUT /posts/{post_id}", web.Handle(controller.UpdatePostById))
	router.HandleFunc("DELETE /posts/{post_id}", web.Handle(controller.DeletePostById))

	return router
}
