package v1

import "github.com/labstack/echo/v4"

type Router struct {
	userHandler UserHandler
}

func NewRouter(
	userHandler UserHandler,
) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

func (r Router) InitRoutes(api *echo.Group) {
	api.GET("/me", r.userHandler.GetUserByEmail)
}
