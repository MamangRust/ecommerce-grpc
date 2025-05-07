package merchantaward_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantInvalidId = response.NewGrpcError("error", "Invalid merchant ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllMerchantAward     = response.NewGrpcError("error", "Failed to fetch all merchant awards", int(codes.Internal))
	ErrGrpcFailedFindMerchantAwardById    = response.NewGrpcError("error", "Failed to fetch merchant award by ID", int(codes.Internal))
	ErrGrpcFailedFindActiveMerchantAward  = response.NewGrpcError("error", "Failed to fetch active merchant awards", int(codes.Internal))
	ErrGrpcFailedFindTrashedMerchantAward = response.NewGrpcError("error", "Failed to fetch trashed merchant awards", int(codes.Internal))
	ErrGrpcFailedCreateMerchantAward      = response.NewGrpcError("error", "Failed to create merchant award", int(codes.Internal))
	ErrGrpcFailedUpdateMerchantAward      = response.NewGrpcError("error", "Failed to update merchant award", int(codes.Internal))
	ErrGrpcValidateCreateMerchantAward    = response.NewGrpcError("error", "Validation failed: invalid create merchant award request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantAward    = response.NewGrpcError("error", "Validation failed: invalid update merchant award request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashMerchantAward              = response.NewGrpcError("error", "Failed to trash merchant award", int(codes.Internal))
	ErrGrpcFailedRestoreMerchantAward            = response.NewGrpcError("error", "Failed to restore merchant award", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantAwardPermanent    = response.NewGrpcError("error", "Failed to permanently delete merchant award", int(codes.Internal))
	ErrGrpcFailedRestoreAllMerchantAward         = response.NewGrpcError("error", "Failed to restore all merchant awards", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantAwardPermanent = response.NewGrpcError("error", "Failed to permanently delete all merchant awards", int(codes.Internal))
)
