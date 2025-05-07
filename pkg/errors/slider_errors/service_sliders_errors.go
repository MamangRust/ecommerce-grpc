package slider_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllSliders            = response.NewErrorResponse("failed to fetch sliders", http.StatusInternalServerError)
	ErrFailedFindActiveSliders         = response.NewErrorResponse("failed to fetch active sliders", http.StatusInternalServerError)
	ErrFailedFindTrashedSliders        = response.NewErrorResponse("failed to fetch trashed sliders", http.StatusInternalServerError)
	ErrFailedCreateSlider              = response.NewErrorResponse("failed to create slider", http.StatusInternalServerError)
	ErrFailedUpdateSlider              = response.NewErrorResponse("failed to update slider", http.StatusInternalServerError)
	ErrFailedTrashSlider               = response.NewErrorResponse("failed to trash slider", http.StatusInternalServerError)
	ErrFailedRestoreSlider             = response.NewErrorResponse("failed to restore slider", http.StatusInternalServerError)
	ErrFailedDeletePermanentSlider     = response.NewErrorResponse("failed to permanently delete slider", http.StatusInternalServerError)
	ErrFailedRestoreAllSliders         = response.NewErrorResponse("failed to restore all sliders", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentSliders = response.NewErrorResponse("failed to permanently delete all sliders", http.StatusInternalServerError)
)
