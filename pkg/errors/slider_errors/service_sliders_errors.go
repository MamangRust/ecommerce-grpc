package slider_errors

import (
	"ecommerce/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllSliders            = errors.NewErrorResponse("failed to fetch sliders", http.StatusInternalServerError)
	ErrFailedFindActiveSliders         = errors.NewErrorResponse("failed to fetch active sliders", http.StatusInternalServerError)
	ErrFailedFindTrashedSliders        = errors.NewErrorResponse("failed to fetch trashed sliders", http.StatusInternalServerError)
	ErrFailedCreateSlider              = errors.NewErrorResponse("failed to create slider", http.StatusInternalServerError)
	ErrFailedUpdateSlider              = errors.NewErrorResponse("failed to update slider", http.StatusInternalServerError)
	ErrFailedTrashSlider               = errors.NewErrorResponse("failed to trash slider", http.StatusInternalServerError)
	ErrFailedRestoreSlider             = errors.NewErrorResponse("failed to restore slider", http.StatusInternalServerError)
	ErrFailedDeletePermanentSlider     = errors.NewErrorResponse("failed to permanently delete slider", http.StatusInternalServerError)
	ErrFailedFindSliderByID            = errors.NewErrorResponse("failed to fetch slider by ID", http.StatusInternalServerError)
	ErrFailedRestoreAllSliders         = errors.NewErrorResponse("failed to restore all sliders", http.StatusInternalServerError)
	ErrFailedDeleteAllPermanentSliders = errors.NewErrorResponse("failed to permanently delete all sliders", http.StatusInternalServerError)
)
