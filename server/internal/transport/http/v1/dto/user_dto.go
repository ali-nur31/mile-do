package dto

import "github.com/ali-nur31/mile-do/internal/domain"

type GetUserResponse struct {
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func ToGetUserResponse(output *domain.UserOutput) GetUserResponse {
	return GetUserResponse{
		Email:     output.Email,
		CreatedAt: output.CreatedAt.String(),
	}
}
