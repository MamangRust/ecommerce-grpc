package middlewares

import (
	"context"
	"ecommerce/pkg/auth"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	tokenManager auth.TokenManager
	accessible   map[string]bool
}

func NewAuthInterceptor(tokenManager auth.TokenManager, accessible map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{
		tokenManager: tokenManager,
		accessible:   accessible,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if interceptor.accessible[info.FullMethod] {
			return handler(ctx, req)
		}

		userID, err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, "userID", userID)
		return handler(newCtx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	if !strings.HasPrefix(accessToken, "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	token := strings.TrimPrefix(accessToken, "Bearer ")
	userID, err := interceptor.tokenManager.ValidateToken(token)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return userID, nil
}
