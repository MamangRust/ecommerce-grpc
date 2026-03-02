package errors

import (
	"ecommerce/internal/pb"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr *AppError
	if !errors.As(err, &apiErr) {
		return status.Error(codes.Internal, "internal error")
	}

	grpcCode := httpToGrpcCode(apiErr.Code)

	st := status.New(grpcCode, apiErr.Message)

	detail := &pb.ErrorResponse{
		Status:  http.StatusText(apiErr.Code),
		Message: apiErr.Message,
		Code:    int32(apiErr.Code),
	}

	stWithDetails, err := st.WithDetails(detail)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

func httpToGrpcCode(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}

func NewGrpcError(message string, httpCode int) error {
	grpcCode := httpToGrpcCode(httpCode)

	st := status.New(grpcCode, message)

	detail := &pb.ErrorResponse{
		Status:  http.StatusText(httpCode),
		Message: message,
		Code:    int32(httpCode),
	}

	stWithDetails, err := st.WithDetails(detail)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}
