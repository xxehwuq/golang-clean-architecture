package v1

import (
	"github.com/xxehwuq/go-clean-architecture/internal/delivery/fasthttp"
	"github.com/xxehwuq/go-clean-architecture/internal/usecase"
)

func NewRouter(server *fasthttp.Server, userUsecase usecase.UserUsecase) {
	api := server.Group("/api")
	{
		initUsersRoutes(api, userUsecase)
	}
}
