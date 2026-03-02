package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	review_errors "ecommerce/pkg/errors/review"
	"encoding/json"
	"log"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type reviewHandleGrpc struct {
	pb.UnimplementedReviewServiceServer
	reviewService service.ReviewService
}

func NewReviewHandleGrpc(
	reviewService service.ReviewService,

) *reviewHandleGrpc {
	return &reviewHandleGrpc{
		reviewService: reviewService,
	}
}

func (s *reviewHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReview, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewService.FindAllReview(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviews := make([]*pb.ReviewResponse, len(reviews))
	for i, review := range reviews {
		protoReviews[i] = &pb.ReviewResponse{
			Id:        int32(review.ReviewID),
			UserId:    int32(review.UserID),
			ProductId: int32(review.ProductID),
			Name:      review.Name,
			Comment:   review.Comment,
			Rating:    int32(review.Rating),
			CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReview{
		Status:     "success",
		Message:    "Successfully fetched reviews",
		Data:       protoReviews,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewHandleGrpc) FindByProduct(ctx context.Context, request *pb.FindAllReviewProductRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	product_id := int(request.GetProductId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReviewByProduct{
		ProductID: product_id,
		Page:      page,
		PageSize:  pageSize,
		Search:    search,
	}

	reviews, totalRecords, err := s.reviewService.FindByProduct(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviews := make([]*pb.ReviewsDetailResponse, len(reviews))
	for i, review := range reviews {
		protoReview := &pb.ReviewsDetailResponse{
			Id:        int32(review.ReviewID),
			UserId:    int32(review.UserID),
			ProductId: int32(review.ProductID),
			Name:      review.Name,
			Comment:   review.Comment,
			Rating:    int32(review.Rating),
			CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if review.DeletedAt.Valid {
			deletedAt := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
			protoReview.DeletedAt = deletedAt
		}

		if review.ReviewDetails != nil {
			var details []struct {
				DetailID  int    `json:"detail_id"`
				Type      string `json:"type"`
				URL       string `json:"url"`
				Caption   string `json:"caption"`
				CreatedAt string `json:"created_at"`
			}

			detailsBytes, err := json.Marshal(review.ReviewDetails)
			if err != nil {
				log.Printf("Error marshaling review details: %v", err)
			} else {
				err = json.Unmarshal(detailsBytes, &details)
				if err != nil {
					log.Printf("Error unmarshaling review details: %v", err)
				} else if len(details) > 0 {
					firstDetail := details[0]
					protoReview.ReviewDetail = &pb.ReviewDetailResponse{
						Id:        int32(firstDetail.DetailID),
						Type:      firstDetail.Type,
						Url:       firstDetail.URL,
						Caption:   firstDetail.Caption,
						CreatedAt: firstDetail.CreatedAt,
					}
				}
			}
		}

		protoReviews[i] = protoReview
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetail{
		Status:     "success",
		Message:    "Successfully fetched product reviews",
		Data:       protoReviews,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllReviewMerchantRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	merchant_id := int(request.GetMerchantId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReviewByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	reviews, totalRecords, err := s.reviewService.FindByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviews := make([]*pb.ReviewsDetailResponse, len(reviews))
	for i, review := range reviews {
		protoReview := &pb.ReviewsDetailResponse{
			Id:        int32(review.ReviewID),
			UserId:    int32(review.UserID),
			ProductId: int32(review.ProductID),
			Name:      review.Name,
			Comment:   review.Comment,
			Rating:    int32(review.Rating),
			CreatedAt: review.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
			UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		}

		if review.DeletedAt.Valid {
			deletedAt := review.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
			protoReview.DeletedAt = deletedAt
		}

		if review.ReviewDetails != nil {
			var details []struct {
				DetailID  int    `json:"detail_id"`
				Type      string `json:"type"`
				URL       string `json:"url"`
				Caption   string `json:"caption"`
				CreatedAt string `json:"created_at"`
			}

			detailsBytes, err := json.Marshal(review.ReviewDetails)
			if err != nil {
				log.Printf("Error marshaling review details: %v", err)
			} else {
				err = json.Unmarshal(detailsBytes, &details)
				if err != nil {
					log.Printf("Error unmarshaling review details: %v", err)
				} else if len(details) > 0 {
					firstDetail := details[0]
					protoReview.ReviewDetail = &pb.ReviewDetailResponse{
						Id:        int32(firstDetail.DetailID),
						Type:      firstDetail.Type,
						Url:       firstDetail.URL,
						Caption:   firstDetail.Caption,
						CreatedAt: firstDetail.CreatedAt,
					}
				}
			}
		}

		protoReviews[i] = protoReview
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetail{
		Status:     "success",
		Message:    "Successfully fetched merchant reviews",
		Data:       protoReviews,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviews := make([]*pb.ReviewResponseDeleteAt, len(reviews))
	for i, review := range reviews {
		var deletedAt string
		if review.DeletedAt.Valid {
			deletedAt = review.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviews[i] = &pb.ReviewResponseDeleteAt{
			Id:        int32(review.ReviewID),
			UserId:    int32(review.UserID),
			ProductId: int32(review.ProductID),
			Name:      review.Name,
			Comment:   review.Comment,
			Rating:    int32(review.Rating),
			CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active reviews",
		Data:       protoReviews,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviews := make([]*pb.ReviewResponseDeleteAt, len(reviews))
	for i, review := range reviews {
		var deletedAt string
		if review.DeletedAt.Valid {
			deletedAt = review.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviews[i] = &pb.ReviewResponseDeleteAt{
			Id:        int32(review.ReviewID),
			UserId:    int32(review.UserID),
			ProductId: int32(review.ProductID),
			Name:      review.Name,
			Comment:   review.Comment,
			Rating:    int32(review.Rating),
			CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed reviews",
		Data:       protoReviews,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewHandleGrpc) Create(ctx context.Context, request *pb.CreateReviewRequest) (*pb.ApiResponseReview, error) {
	req := &requests.CreateReviewRequest{
		UserID:    int(request.GetUserId()),
		ProductID: int(request.GetProductId()),
		Rating:    int(request.GetRating()),
		Comment:   request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		return nil, review_errors.ErrGrpcValidateCreateReview
	}

	review, err := s.reviewService.CreateReview(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReview := &pb.ReviewResponse{
		Id:        int32(review.ReviewID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReview{
		Status:  "success",
		Message: "Successfully created review",
		Data:    protoReview,
	}, nil
}

func (s *reviewHandleGrpc) Update(ctx context.Context, request *pb.UpdateReviewRequest) (*pb.ApiResponseReview, error) {
	id := int(request.GetReviewId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateReviewRequest{
		ReviewID: &id,
		Name:     request.GetName(),
		Rating:   int(request.GetRating()),
		Comment:  request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		return nil, review_errors.ErrGrpcValidateUpdateReview
	}

	review, err := s.reviewService.UpdateReview(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReview := &pb.ReviewResponse{
		Id:        int32(review.ReviewID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReview{
		Status:  "success",
		Message: "Successfully updated review",
		Data:    protoReview,
	}, nil
}

func (s *reviewHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	review, err := s.reviewService.TrashReview(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if review.DeletedAt.Valid {
		deletedAt = review.DeletedAt.Time.Format("2006-01-02")
	}

	protoReview := &pb.ReviewResponseDeleteAt{
		Id:        int32(review.ReviewID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDeleteAt{
		Status:  "success",
		Message: "Successfully trashed review",
		Data:    protoReview,
	}, nil
}

func (s *reviewHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	review, err := s.reviewService.RestoreReview(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if review.DeletedAt.Valid {
		deletedAt = review.DeletedAt.Time.Format("2006-01-02")
	}

	protoReview := &pb.ReviewResponseDeleteAt{
		Id:        int32(review.ReviewID),
		UserId:    int32(review.UserID),
		ProductId: int32(review.ProductID),
		Name:      review.Name,
		Comment:   review.Comment,
		Rating:    int32(review.Rating),
		CreatedAt: review.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: review.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDeleteAt{
		Status:  "success",
		Message: "Successfully restored review",
		Data:    protoReview,
	}, nil
}

func (s *reviewHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	_, err := s.reviewService.DeleteReviewPermanently(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewDelete{
		Status:  "success",
		Message: "Successfully deleted review permanently",
	}, nil
}

func (s *reviewHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewService.RestoreAllReview(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully restored all reviews",
	}, nil
}

func (s *reviewHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewService.DeleteAllPermanentReview(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully deleted all reviews permanently",
	}, nil
}
