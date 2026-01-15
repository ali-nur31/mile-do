package dto

import "github.com/ali-nur31/mile-do/internal/domain"

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=32"`
}

type RegisterUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=72"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
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
