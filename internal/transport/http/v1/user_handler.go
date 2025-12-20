package v1

import (
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type createUserInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type getUserInput struct {
	Email string `json:"email"`
}

type GetUserOutput struct {
	Email     string           `json:"email"`
	Password  pgtype.Text      `json:"password"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (h *UserHandler) GetUserByEmail(c echo.Context) error {
	var input getUserInput

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if !strings.Contains(input.Email, "@") {
		return c.JSON(http.StatusBadRequest, "email is invalid")
	}

	user, err := h.service.GetUser(c.Request().Context(), input.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	output := GetUserOutput{
		Email:     user.Email,
		Password:  user.PasswordHash,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusFound, output)
}
