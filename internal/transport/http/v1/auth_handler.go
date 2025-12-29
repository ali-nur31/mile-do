package v1

import (
	"net/http"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
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

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(service service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var request registerUserInput

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	if request.Password != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, "passwords do not match")
	}

	input := domain.UserInput{
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

	input := domain.UserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	data, err := h.service.LoginUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, loginUserOutput{data.Token})
}
