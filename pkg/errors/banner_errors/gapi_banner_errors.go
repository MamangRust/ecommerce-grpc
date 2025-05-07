package banner_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcBannerNotFound  = response.NewGrpcError("error", "Banner not found", int(codes.NotFound))
	ErrGrpcBannerInvalidId = response.NewGrpcError("error", "Invalid Banner ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAll     = response.NewGrpcError("error", "Failed to fetch banners", int(codes.Internal))
	ErrGrpcFailedFindActive  = response.NewGrpcError("error", "Failed to fetch active banners", int(codes.Internal))
	ErrGrpcFailedFindTrashed = response.NewGrpcError("error", "Failed to fetch trashed banners", int(codes.Internal))

	ErrGrpcFailedCreateBanner   = response.NewGrpcError("error", "Failed to create banner", int(codes.Internal))
	ErrGrpcFailedUpdateBanner   = response.NewGrpcError("error", "Failed to update banner", int(codes.Internal))
	ErrGrpcValidateCreateBanner = response.NewGrpcError("error", "Validation failed: invalid create banner request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateBanner = response.NewGrpcError("error", "Validation failed: invalid update banner request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedBanner   = response.NewGrpcError("error", "Failed to move banner to trash", int(codes.Internal))
	ErrGrpcFailedRestoreBanner   = response.NewGrpcError("error", "Failed to restore banner", int(codes.Internal))
	ErrGrpcFailedDeletePermanent = response.NewGrpcError("error", "Failed to delete banner permanently", int(codes.Internal))

	ErrGrpcFailedRestoreAllBanners = response.NewGrpcError("error", "Failed to restore all banners", int(codes.Internal))
	ErrGrpcFailedDeleteAllBanners  = response.NewGrpcError("error", "Failed to delete all banners permanently", int(codes.Internal))
)
