package routers

import (
	"net/http"

	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/core/service"
	"github.com/skrpld/NearBeee/internal/transport/rest/handlers"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
)

func NewAuthRouter(repo *repository.NearBeeeRepository, secret string) (*http.ServeMux, *service.AuthService) {
	srv := service.NewAuthService(repo, secret)
	controller := handlers.NewAuthController(srv)
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", web.Handle(controller.RegistrateUserHandler))
	router.HandleFunc("POST /auth/login", web.Handle(controller.LoginUserHandler))
	router.HandleFunc("POST /auth/refresh-token", web.Handle(controller.RefreshUserTokenHandler))

	return router, srv
}
