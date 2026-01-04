package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/labstack/echo/v4"
)

type getUserResponse struct {
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUser godoc
// @Summary      get user info
// @Description  get user account by bearer token
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  getUserResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /users/me [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), int64(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	output := getUserResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	return c.JSON(http.StatusOK, output)
}

// LogoutUser godoc
// @Summary      logout user
// @Description  logout from user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string "successful log out"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /users/ [delete]
func (h *UserHandler) LogoutUser(c echo.Context) error {
	userId, err := getCurrentUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	err = h.service.LogoutUser(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "successful log out"})
}

func getCurrentUserIDFromToken(c echo.Context) (int32, error) {
	switch t := c.Get("userId").(type) {
	case int64:
		return int32(t), nil
	default:
		slog.Error("userId in context is not an integer", "value", t)
		return -1, fmt.Errorf("failed to convert userId from string to integer")
	}
}
