package middlewares

import (
	"context"
	"ecommerce/pkg/logger"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func SecurityLoggingInterceptor(l logger.LoggerInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		st, _ := status.FromError(err)

		var userID string
		if val := ctx.Value("userID"); val != nil {
			userID = val.(string)
		}

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("status_code", st.Code().String()),
			zap.Duration("duration", duration),
			zap.String("user_id", userID),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			l.Warn("Security event: Request failed", fields...)
		} else {
			// Don't log everything at Info level to avoid noise, but maybe for sensitive methods?
			// For now, log all at Info to ensure visibility during implementation.
			l.Info("gRPC request", fields...)
		}

		return resp, err
	}
}
