package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ali-nur31/mile-do/internal/domain"
	"github.com/ali-nur31/mile-do/internal/transport/http/v1/dto"
	"github.com/ali-nur31/mile-do/pkg/validator"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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

	if validateErrors := validator.ValidateStruct(request); validateErrors != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "validation failed", "details": validateErrors})
	}

	if request.Password != request.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request", "error": "passwords do not match"})
	}

	output, err := h.authService.RegisterUser(c.Request().Context(), domain.AuthInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		slog.Error("failed on register", "error", err)
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

	if validateErrors := validator.ValidateStruct(request); validateErrors != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "validation failed", "details": validateErrors})
	}

	output, err := h.authService.LoginUser(c.Request().Context(), domain.AuthInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		slog.Error("failed on login", "error", err)
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

	if validateErrors := validator.ValidateStruct(request); validateErrors != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "validation failed", "details": validateErrors})
	}

	output, err := h.authService.RefreshTokens(c.Request().Context(), request.RefreshToken)
	if err != nil {
		slog.Error("failed on refreshing token", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToAuthUserResponse(output))
}

// LogoutUser godoc
// @Summary      logout user
// @Description  logout from user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string "successful log out"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /auth/logout [delete]
func (h *AuthHandler) LogoutUser(c echo.Context) error {
	claims, err := GetCurrentClaimsFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	accessToken := c.Get("accessToken")

	err = h.authService.LogoutUser(c.Request().Context(), int32(claims.ID), fmt.Sprint(accessToken), claims.ExpiresAt.Time)
	if err != nil {
		slog.Error("failed on logout", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "successful log out"})
}

func GetCurrentClaimsFromCtx(c echo.Context) (*domain.Claims, error) {
	switch t := c.Get("claims").(type) {
	case *domain.Claims:
		return t, nil
	default:
		slog.Error("userId in context is not an integer", "value", t)
		return nil, fmt.Errorf("failed to convert userId to integer")
	}
}
