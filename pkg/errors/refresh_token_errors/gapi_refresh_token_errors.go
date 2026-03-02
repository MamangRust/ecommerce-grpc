package refreshtoken_errors

import (
	"ecommerce/pkg/errors"

	"google.golang.org/grpc/codes"
)

var ErrGrpcRefreshToken = errors.NewGrpcError("refresh token failed", int(codes.Unauthenticated))
