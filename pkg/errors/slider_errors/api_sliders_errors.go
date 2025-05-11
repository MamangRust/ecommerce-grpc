package slider_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidName = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid slider name", http.StatusBadRequest)
	}
	ErrApiImageRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "slider image is required", http.StatusBadRequest)
	}

	ErrApiInvalidBody = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid request body", http.StatusBadRequest)
	}

	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid slider ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllSliders = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch sliders", http.StatusInternalServerError)
	}

	ErrApiFailedFindActiveSliders = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active sliders", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashedSliders = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed sliders", http.StatusInternalServerError)
	}

	ErrApiFailedCreateSlider = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create slider", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateSlider = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update slider", http.StatusInternalServerError)
	}

	ErrApiFailedTrashSlider = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash slider", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreSlider = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore slider", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteSliderPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete slider", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllSliders = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all sliders", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllPermanentSliders = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all sliders", http.StatusInternalServerError)
	}
)
