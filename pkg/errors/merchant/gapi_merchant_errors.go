package merchant_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantId          = response.NewGrpcError("error", "invalid merchant ID", int(codes.InvalidArgument))
	ErrGrpcFailedFindAllMerchants     = response.NewGrpcError("error", "Failed to fetch all merchants", int(codes.Internal))
	ErrGrpcFailedFindMerchantById     = response.NewGrpcError("error", "Failed to fetch merchant by ID", int(codes.Internal))
	ErrGrpcFailedFindActiveMerchants  = response.NewGrpcError("error", "Failed to fetch active merchants", int(codes.Internal))
	ErrGrpcFailedFindTrashedMerchants = response.NewGrpcError("error", "Failed to fetch trashed merchants", int(codes.Internal))

	ErrGrpcFailedCreateMerchant   = response.NewGrpcError("error", "Failed to create merchant", int(codes.Internal))
	ErrGrpcFailedUpdateMerchant   = response.NewGrpcError("error", "Failed to update merchant", int(codes.Internal))
	ErrGrpcValidateCreateMerchant = response.NewGrpcError("error", "Validation failed: invalid create merchant request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = response.NewGrpcError("error", "Validation failed: invalid update merchant request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashMerchant               = response.NewGrpcError("error", "Failed to trash merchant", int(codes.Internal))
	ErrGrpcFailedRestoreMerchant             = response.NewGrpcError("error", "Failed to restore merchant", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantPermanent     = response.NewGrpcError("error", "Failed to permanently delete merchant", int(codes.Internal))
	ErrGrpcFailedRestoreAllMerchants         = response.NewGrpcError("error", "Failed to restore all merchants", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantsPermanent = response.NewGrpcError("error", "Failed to permanently delete all merchants", int(codes.Internal))
)
