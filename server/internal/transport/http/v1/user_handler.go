package v1

import (
	"log/slog"
	"net/http"

	"github.com/ali-nur31/mile-do/internal/service"
	"github.com/ali-nur31/mile-do/internal/transport/http/v1/dto"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser godoc
// @Summary      get user info
// @Description  get user account by bearer token
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.GetUserResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /users/me [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	userId, err := GetCurrentUserIdFromCtx(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), int64(userId))
	if err != nil {
		slog.Error("failed on getting user by id", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error", "error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ToGetUserResponse(user))
}
