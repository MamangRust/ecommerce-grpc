package review_errors

import (
	"ecommerce/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiReviewInvalidProductId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid product ID", http.StatusBadRequest)
	}

	ErrApiReviewInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid merchant ID", http.StatusBadRequest)
	}

	ErrApiReviewNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Review not found", http.StatusNotFound)
	}

	ErrApiReviewInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid Review ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch reviews", http.StatusInternalServerError)
	}

	ErrApiFailedFindProductReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch reviews by product", http.StatusInternalServerError)
	}

	ErrApiFailedFindActiveReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active reviews", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashedReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed reviews", http.StatusInternalServerError)
	}

	ErrApiFailedFindMerchantReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch reviews by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedCreateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create review", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update review", http.StatusInternalServerError)
	}

	ErrApiBindCreateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to bind create review", http.StatusInternalServerError)
	}

	ErrApiBindUpdateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to bind update review", http.StatusInternalServerError)
	}

	ErrApiValidateCreateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to validate create review", http.StatusInternalServerError)
	}

	ErrApiValidateUpdateReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to validate update review", http.StatusInternalServerError)
	}

	ErrApiFailedTrashedReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash review", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreReview = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore review", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteReviewPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to delete review permanently", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllReviews = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all reviews", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllReviewsPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to delete all reviews permanently", http.StatusInternalServerError)
	}
)
