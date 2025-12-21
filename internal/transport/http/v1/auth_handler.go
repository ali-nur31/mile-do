package v1

import (
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

type RegisterUserOutput struct {
	Token string `json:"token"`
}

type loginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserOutput struct {
	Token string `json:"token"`
}

type getUserInput struct {
	Email string `json:"email"`
}

type GetUserOutput struct {
	Email     string           `json:"email"`
	Password  pgtype.Text      `json:"password"`
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
	var input getUserInput

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
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

	return c.JSON(http.StatusCreated, RegisterUserOutput{data.Token})
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
