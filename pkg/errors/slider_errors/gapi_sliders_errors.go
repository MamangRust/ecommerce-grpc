package slider_errors

import (
	"ecommerce/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidID = response.NewGrpcError("error", "invalid ID", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllSliders     = response.NewGrpcError("error", "Failed to fetch all sliders", int(codes.Internal))
	ErrGrpcFailedFindActiveSliders  = response.NewGrpcError("error", "Failed to fetch active sliders", int(codes.Internal))
	ErrGrpcFailedFindTrashedSliders = response.NewGrpcError("error", "Failed to fetch trashed sliders", int(codes.Internal))
	ErrGrpcFailedCreateSlider       = response.NewGrpcError("error", "Failed to create slider", int(codes.Internal))
	ErrGrpcFailedUpdateSlider       = response.NewGrpcError("error", "Failed to update slider", int(codes.Internal))

	ErrGrpcValidateCreateSlider = response.NewGrpcError("error", "validation failed: invalid create slider request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSlider = response.NewGrpcError("error", "validation failed: invalid update slider request", int(codes.InvalidArgument))

	ErrGrpcFailedTrashSlider               = response.NewGrpcError("error", "Failed to trash slider", int(codes.Internal))
	ErrGrpcFailedRestoreSlider             = response.NewGrpcError("error", "Failed to restore slider", int(codes.Internal))
	ErrGrpcFailedDeletePermanentSlider     = response.NewGrpcError("error", "Failed to permanently delete slider", int(codes.Internal))
	ErrGrpcFailedRestoreAllSliders         = response.NewGrpcError("error", "Failed to restore all sliders", int(codes.Internal))
	ErrGrpcFailedDeleteAllPermanentSliders = response.NewGrpcError("error", "Failed to permanently delete all sliders", int(codes.Internal))
)
