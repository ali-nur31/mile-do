package v1

import (
	"net/http"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/labstack/echo/v4"
)

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type registerUserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
// @Param        input body registerUserRequest true "Account Info"
// @Success      201  {object}  registerUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}   map[string]string "Internal Server Error"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var request registerUserRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if request.Password != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "passwords do not match"})
	}

	input := domain.UserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	data, err := h.service.CreateUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, authUserResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	})
}

// LoginUser godoc
// @Summary      login user
// @Description  login to existing user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body loginUserRequest true "Account Info"
// @Success      202  {object}  loginUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/login [post]
func (h *AuthHandler) LoginUser(c echo.Context) error {
	var request loginUserRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	input := domain.UserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	data, err := h.service.LoginUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, authUserResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	})
}

// RefreshAccessToken godoc
// @Summary      refresh access token
// @Description  refresh access token by refresh_token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body refreshAccessTokenRequest true "Refresh token"
// @Success      202  {object}  authUserResponse
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c echo.Context) error {
	var request refreshAccessTokenRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	output, err := h.service.RefreshTokens(c.Request().Context(), request.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, authUserResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
