package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type authResponseMapper struct {
}

func NewAuthResponseMapper() *authResponseMapper {
	return &authResponseMapper{}
}

func (s *authResponseMapper) ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin {
	return &response.ApiResponseLogin{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		},
	}
}

func (s *authResponseMapper) ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister {
	return &response.ApiResponseRegister{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		},
	}
}

func (s *authResponseMapper) ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken {
	return &response.ApiResponseRefreshToken{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		},
	}
}

func (s *authResponseMapper) ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe {
	return &response.ApiResponseGetMe{
		Status:  res.Status,
		Message: res.Message,
		Data: &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		},
	}
}
