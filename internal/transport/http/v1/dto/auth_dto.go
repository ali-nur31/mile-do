package dto

import "github.com/ali-nur31/mile-do/internal/domain"

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ToAuthUserResponse(output *domain.AuthOutput) AuthUserResponse {
	return AuthUserResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}
