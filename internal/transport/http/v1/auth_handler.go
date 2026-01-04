package v1

import (
	"net/http"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/v1/dto"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(service service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// RegisterUser godoc
// @Summary      register new user
// @Description  create new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.RegisterUserRequest true "Account Info"
// @Success      201  {object}  dto.AuthUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var request dto.RegisterUserRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	if request.Password != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": "passwords do not match"})
	}

	output, err := h.service.CreateUser(c.Request().Context(), domain.UserInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ToAuthUserResponse(output))
}

// LoginUser godoc
// @Summary      login user
// @Description  login to existing user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.LoginUserRequest true "Account Info"
// @Success      202  {object}  dto.AuthUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/login [post]
func (h *AuthHandler) LoginUser(c echo.Context) error {
	var request dto.LoginUserRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	output, err := h.service.LoginUser(c.Request().Context(), domain.UserInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, dto.ToAuthUserResponse(output))
}

// RefreshAccessToken godoc
// @Summary      refresh access token
// @Description  refresh access token by refresh_token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body dto.RefreshAccessTokenRequest true "Refresh token"
// @Success      200  {object}  dto.AuthUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c echo.Context) error {
	var request dto.RefreshAccessTokenRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": err.Error()})
	}

	output, err := h.service.RefreshTokens(c.Request().Context(), request.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToAuthUserResponse(output))
}
