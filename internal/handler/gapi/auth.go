package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthHandleGrpc(auth service.AuthService) *authHandleGrpc {
	return &authHandleGrpc{authService: auth}
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.authService.Login(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseLogin{
		Status:  "success",
		Message: "Login successfully",
		Data: &pb.TokenResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
	}

	return pbRes, nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	res, err := s.authService.RefreshToken(ctx, req.RefreshToken)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseRefreshToken{
		Status:  "success",
		Message: "Refresh token successfully",
		Data: &pb.TokenResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
	}

	return pbRes, nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	res, err := s.authService.GetMe(ctx, req.AccessToken)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseGetMe{
		Status:  "success",
		Message: "Get me successfully",
		Data: &pb.UserResponse{
			Id:        res.UserID,
			Firstname: res.Firstname,
			Lastname:  res.Lastname,
			Email:     res.Email,
			CreatedAt: res.CreatedAt.Time.String(),
			UpdatedAt: res.UpdatedAt.Time.String(),
		},
	}

	return pbRes, nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	request := &requests.CreateUserRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.authService.Register(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseRegister{
		Status:  "success",
		Message: "Registration successfully",
		Data: &pb.UserResponse{
			Id:        res.UserID,
			Firstname: res.Firstname,
			Lastname:  res.Lastname,
			Email:     res.Email,
			CreatedAt: res.CreatedAt.Time.String(),
			UpdatedAt: res.UpdatedAt.Time.String(),
		},
	}

	return pbRes, nil
}
