package merchantbusiness_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcFailedFindAllMerchantBusiness       = response.NewGrpcError("error", "Failed to fetch merchant businesses", int(codes.Internal))
	ErrGrpcFailedFindByIdMerchantBusiness      = response.NewGrpcError("error", "Failed to find merchant business by ID", int(codes.Internal))
	ErrGrpcFailedFindByActiveMerchantBusiness  = response.NewGrpcError("error", "Failed to fetch active merchant businesses", int(codes.Internal))
	ErrGrpcFailedFindByTrashedMerchantBusiness = response.NewGrpcError("error", "Failed to fetch trashed merchant businesses", int(codes.Internal))

	ErrGrpcFailedCreateMerchantBusiness   = response.NewGrpcError("error", "Failed to create merchant business", int(codes.Internal))
	ErrGrpcFailedUpdateMerchantBusiness   = response.NewGrpcError("error", "Failed to update merchant business", int(codes.Internal))
	ErrGrpcValidateCreateMerchantBusiness = response.NewGrpcError("error", "Validation failed: invalid create merchant business request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantBusiness = response.NewGrpcError("error", "Validation failed: invalid update merchant business request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashMerchantBusiness           = response.NewGrpcError("error", "Failed to trash merchant business", int(codes.Internal))
	ErrGrpcFailedRestoreMerchantBusiness         = response.NewGrpcError("error", "Failed to restore merchant business", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantBusinessPermanent = response.NewGrpcError("error", "Failed to permanently delete merchant business", int(codes.Internal))

	ErrGrpcFailedRestoreAllMerchantBusiness         = response.NewGrpcError("error", "Failed to restore all merchant businesses", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantBusinessPermanent = response.NewGrpcError("error", "Failed to permanently delete all merchant businesses", int(codes.Internal))

	ErrGrpcMerchantBusinessNotFound  = response.NewGrpcError("error", "Merchant business not found", int(codes.NotFound))
	ErrGrpcInvalidMerchantBusinessId = response.NewGrpcError("error", "Invalid merchant business ID", int(codes.InvalidArgument))
)
