package merchantpolicy_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidMerchantPolicyID = response.NewGrpcError("error", "Invalid merchant policy ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllMerchantPolicies     = response.NewGrpcError("error", "Failed to fetch all merchant policies", int(codes.Internal))
	ErrGrpcFailedFindActiveMerchantPolicies  = response.NewGrpcError("error", "Failed to fetch active merchant policies", int(codes.Internal))
	ErrGrpcFailedFindTrashedMerchantPolicies = response.NewGrpcError("error", "Failed to fetch trashed merchant policies", int(codes.Internal))
	ErrGrpcFailedFindMerchantPolicyById      = response.NewGrpcError("error", "Failed to find merchant policy by ID", int(codes.Internal))

	ErrGrpcFailedCreateMerchantPolicy   = response.NewGrpcError("error", "Failed to create merchant policy", int(codes.Internal))
	ErrGrpcFailedUpdateMerchantPolicy   = response.NewGrpcError("error", "Failed to update merchant policy", int(codes.Internal))
	ErrGrpcValidateCreateMerchantPolicy = response.NewGrpcError("error", "Validation failed: invalid create merchant policy request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchantPolicy = response.NewGrpcError("error", "Validation failed: invalid update merchant policy request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedMerchantPolicy              = response.NewGrpcError("error", "Failed to trash merchant policy", int(codes.Internal))
	ErrGrpcFailedRestoreMerchantPolicy              = response.NewGrpcError("error", "Failed to restore merchant policy", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantPolicyPermanent      = response.NewGrpcError("error", "Failed to permanently delete merchant policy", int(codes.Internal))
	ErrGrpcFailedRestoreAllMerchantPolicies         = response.NewGrpcError("error", "Failed to restore all merchant policies", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantPoliciesPermanent = response.NewGrpcError("error", "Failed to permanently delete all merchant policies", int(codes.Internal))
)
