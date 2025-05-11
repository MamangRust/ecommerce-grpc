package reviewdetail_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidBody = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "invalid_body", "Please provide a valid request body", http.StatusBadRequest)
	}

	ErrApiInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "invalid_id", "Please provide a valid ID", http.StatusBadRequest)
	}

	ErrApiInvalidReviewId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "invalid_review_id", "Please provide a valid review ID", http.StatusBadRequest)
	}

	ErrApiReviewDetailTypeRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Type is required", http.StatusBadRequest)
	}

	ErrApiReviewDetailCaptionRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "validation_error", "Caption is required", http.StatusBadRequest)
	}

	ErrApiReviewDetailFileRequired = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "file_required", "A file must be uploaded for review detail", http.StatusBadRequest)
	}

	ErrApiReviewDetailNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Review detail not found", http.StatusNotFound)
	}

	ErrApiReviewDetailInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid review detail ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllReviewDetails = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch review details", http.StatusInternalServerError)
	}

	ErrApiFailedFindActiveReviewDetails = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active review details", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashedReviewDetails = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed review details", http.StatusInternalServerError)
	}

	ErrApiFailedCreateReviewDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create review detail", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateReviewDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update review detail", http.StatusInternalServerError)
	}

	ErrApiFailedTrashedReviewDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash review detail", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreReviewDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore review detail", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteReviewDetailPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to delete review detail permanently", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllReviewDetail = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all review details", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllReviewDetailPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to delete all review details permanently", http.StatusInternalServerError)
	}
)
