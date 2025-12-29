package v1

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type registerUserInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type registerUserOutput struct {
	Token string `json:"token"`
}

type loginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserOutput struct {
	Token string `json:"token"`
}

type getUserOutput struct {
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(service service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) GetUserByEmail(c echo.Context) error {
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

func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var request registerUserInput

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	if request.Password != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, "passwords do not match")
	}

	input := service.UserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	data, err := h.service.CreateUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, registerUserOutput{data.Token})
}

func (h *AuthHandler) LoginUser(c echo.Context) error {
	var request loginUserInput

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	input := service.UserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	data, err := h.service.LoginUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, loginUserOutput{data.Token})
}
