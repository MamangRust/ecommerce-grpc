package refreshtoken_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var ErrGrpcRefreshToken = response.NewGrpcError("error", "refresh token failed", int(codes.Unauthenticated))
