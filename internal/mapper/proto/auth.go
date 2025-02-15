package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type authProtoMapper struct {
}

func NewAuthProtoMapper() *authProtoMapper {
	return &authProtoMapper{}
}

func (s *authProtoMapper) ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin {
	return &pb.ApiResponseLogin{
		Status:  status,
		Message: message,
		Data: &pb.TokenResponse{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister {
	return &pb.ApiResponseRegister{
		Status:  status,
		Message: message,
		Data: &pb.UserResponse{
			Id:        int32(response.ID),
			Firstname: response.FirstName,
			Lastname:  response.LastName,
			Email:     response.Email,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken {
	return &pb.ApiResponseRefreshToken{
		Status:  status,
		Message: message,
		Data: &pb.TokenResponse{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	}
}

func (s *authProtoMapper) ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe {
	return &pb.ApiResponseGetMe{
		Status:  status,
		Message: message,
		Data: &pb.UserResponse{
			Id:        int32(response.ID),
			Firstname: response.FirstName,
			Lastname:  response.LastName,
			Email:     response.Email,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		},
	}
}
