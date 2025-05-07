package merchantdetail_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantDetailId = response.NewGrpcError("error", "invalid merchant detail ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllMerchantDetail     = response.NewGrpcError("error", "Failed to fetch all merchant details", int(codes.Internal))
	ErrGrpcFailedFindActiveMerchantDetail  = response.NewGrpcError("error", "Failed to fetch active merchant details", int(codes.Internal))
	ErrGrpcFailedFindTrashedMerchantDetail = response.NewGrpcError("error", "Failed to fetch trashed merchant details", int(codes.Internal))
	ErrGrpcFailedFindMerchantDetailById    = response.NewGrpcError("error", "Failed to fetch merchant detail by ID", int(codes.Internal))
	ErrGrpcFailedCreateMerchantDetail      = response.NewGrpcError("error", "Failed to create merchant detail", int(codes.Internal))
	ErrGrpcFailedUpdateMerchantDetail      = response.NewGrpcError("error", "Failed to update merchant detail", int(codes.Internal))
	ErrGrpcValidateCreateMerchantDetail    = response.NewGrpcError("error", "Validation failed: invalid create merchant detail request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantDetail    = response.NewGrpcError("error", "Validation failed: invalid update merchant detail request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedMerchantDetail            = response.NewGrpcError("error", "Failed to trash merchant detail", int(codes.Internal))
	ErrGrpcFailedRestoreMerchantDetail            = response.NewGrpcError("error", "Failed to restore merchant detail", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantDetailPermanent    = response.NewGrpcError("error", "Failed to permanently delete merchant detail", int(codes.Internal))
	ErrGrpcFailedRestoreAllMerchantDetail         = response.NewGrpcError("error", "Failed to restore all merchant details", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantDetailPermanent = response.NewGrpcError("error", "Failed to permanently delete all merchant details", int(codes.Internal))
)
