package v1

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type getUserOutput struct {
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

	output := getUserOutput{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusFound, output)
}
