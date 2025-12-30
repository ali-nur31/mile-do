package v1

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type getUserResponse struct {
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUserByEmail godoc
// @Summary      get user info
// @Description  get user account by bearer token
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      302  {object}  getUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /users/me [get]
func (h *UserHandler) GetUserByEmail(c echo.Context) error {
	emailFromCtx := c.Get("email")

	if emailFromCtx == nil {
		slog.Error("email not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	email, ok := emailFromCtx.(string)
	if !ok {
		slog.Error("email in context is not a string", "value", emailFromCtx)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	if !strings.Contains(email, "@") {
		return c.JSON(http.StatusBadRequest, "email is invalid")
	}

	user, err := h.service.GetUser(c.Request().Context(), email)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	output := getUserResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusFound, output)
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
