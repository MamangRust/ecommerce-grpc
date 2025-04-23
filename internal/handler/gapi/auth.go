package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors_custom"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	mapping     protomapper.AuthProtoMapper
}

func NewAuthHandleGrpc(auth service.AuthService, mapping protomapper.AuthProtoMapper) *authHandleGrpc {
	return &authHandleGrpc{authService: auth, mapping: mapping}
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
			Status:  "invalid_request",
			Message: "request cannot be nil",
			Code:    int32(codes.InvalidArgument),
		}))
	}

	request := &requests.AuthRequest{
		Email:    strings.TrimSpace(req.Email),
		Password: req.Password,
	}

	res, errRes := s.authService.Login(request)

	if errRes != nil {
		if errRes.Code == http.StatusUnauthorized {
			return nil, status.Error(
				codes.Unauthenticated,
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  "invalid_credentials",
					Message: "invalid email or password",
					Code:    int32(codes.Unauthenticated),
				}),
			)
		}
		return nil, status.Errorf(
			codes.Code(errRes.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  errRes.Status,
				Message: errRes.Message,
				Code:    int32(errRes.Code),
			}),
		)
	}

	if res == nil || res.AccessToken == "" {
		return nil, status.Error(
			codes.Internal,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "token_generation_failed",
				Message: "authentication service returned empty token",
				Code:    int32(codes.Internal),
			}),
		)
	}

	return s.mapping.ToProtoResponseLogin("success", "Login successful", res), nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "refresh token is required",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.authService.RefreshToken(req.RefreshToken)

	if err != nil {
		switch {
		case err.Code == http.StatusUnauthorized:
			return nil, status.Error(
				codes.Unauthenticated,
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  "invalid_refresh_token",
					Message: "invalid or expired refresh token",
					Code:    int32(codes.Unauthenticated),
				}),
			)
		case err.Code == http.StatusInternalServerError:
			return nil, status.Error(
				codes.Internal,
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  "refresh_failed",
					Message: err.Message,
					Code:    int32(codes.Internal),
				}),
			)
		default:
			return nil, status.Errorf(
				codes.Code(err.Code),
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  err.Status,
					Message: err.Message,
					Code:    int32(err.Code),
				}),
			)
		}
	}

	if res == nil || res.AccessToken == "" || res.RefreshToken == "" {
		return nil, status.Error(
			codes.Internal,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "token_generation_failed",
				Message: "failed to generate valid tokens",
				Code:    int32(codes.Internal),
			}),
		)
	}

	return s.mapping.ToProtoResponseRefreshToken("success", "tokens refreshed successfully", res), nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	if req == nil || req.AccessToken == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "access token is required",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.authService.GetMe(req.AccessToken)

	if err != nil {
		switch {
		case err.Code == http.StatusUnauthorized:
			return nil, status.Error(
				codes.Unauthenticated,
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  "invalid_token",
					Message: "invalid or expired access token",
					Code:    int32(codes.Unauthenticated),
				}),
			)
		case err.Code == http.StatusInternalServerError:
			return nil, status.Error(
				codes.Internal,
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  "server_error",
					Message: err.Message,
					Code:    int32(codes.Internal),
				}),
			)
		default:
			return nil, status.Errorf(
				codes.Code(err.Code),
				errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
					Status:  err.Status,
					Message: err.Message,
					Code:    int32(err.Code),
				}),
			)
		}
	}

	if res == nil {
		return nil, status.Error(
			codes.Internal,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "user_not_found",
				Message: "failed to retrieve user data",
				Code:    int32(codes.Internal),
			}),
		)
	}

	return s.mapping.ToProtoResponseGetMe("success", "user data retrieved successfully", res), nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	if req == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "registration request cannot be nil",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	request := &requests.CreateUserRequest{
		FirstName:       strings.TrimSpace(req.Firstname),
		LastName:        strings.TrimSpace(req.Lastname),
		Email:           strings.TrimSpace(req.Email),
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	if err := request.Validate(); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: err.Error(),
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.authService.Register(request)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	if res == nil {
		return nil, status.Error(
			codes.Internal,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "registration_failed",
				Message: "failed to complete registration",
				Code:    int32(codes.Internal),
			}),
		)
	}

	return s.mapping.ToProtoResponseRegister("success", "Registration successful", res), nil
}
