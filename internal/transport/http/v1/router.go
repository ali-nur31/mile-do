package v1

import "github.com/labstack/echo/v4"

type Router struct {
	authHandler AuthHandler
}

func NewRouter(
	authHandler AuthHandler,
) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

func (r Router) InitRoutes(api *echo.Group) {
	auth := api.Group("/auth")

	{
		auth.POST("/register", r.authHandler.RegisterUser)
		auth.POST("/login", r.authHandler.LoginUser)
		auth.GET("/me", r.authHandler.GetUserByEmail)
	}
}
