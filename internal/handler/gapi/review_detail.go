package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	reviewdetail_errors "ecommerce/pkg/errors/review_detail"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type reviewDetailHandleGrpc struct {
	pb.UnimplementedReviewDetailServiceServer
	reviewDetailService service.ReviewDetailService
}

func NewReviewDetailHandleGrpc(
	reviewDetailService service.ReviewDetailService,
) *reviewDetailHandleGrpc {
	return &reviewDetailHandleGrpc{
		reviewDetailService: reviewDetailService,
	}
}

func (s *reviewDetailHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetails, error) {
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

	reviewDetails, totalRecords, err := s.reviewDetailService.FindAllReviews(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponse, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		protoReviewDetails[i] = &pb.ReviewDetailsResponse{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetails{
		Status:     "success",
		Message:    "Successfully fetched review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.reviewDetailService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetail := &pb.ReviewDetailsResponse{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully fetched review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
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

	reviewDetails, totalRecords, err := s.reviewDetailService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		var deletedAt string
		if reviewDetail.DeletedAt.Valid {
			deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviewDetails[i] = &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
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

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewDetailHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
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

	reviewDetails, totalRecords, err := s.reviewDetailService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		var deletedAt string
		if reviewDetail.DeletedAt.Valid {
			deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviewDetails[i] = &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
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

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewDetailHandleGrpc) Create(ctx context.Context, request *pb.CreateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	req := &requests.CreateReviewDetailRequest{
		ReviewID: int(request.GetReviewId()),
		Type:     request.GetType(),
		Url:      request.GetUrl(),
		Caption:  request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		return nil, reviewdetail_errors.ErrGrpcValidateCreateReviewDetail
	}

	reviewDetail, err := s.reviewDetailService.CreateReviewDetail(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetail := &pb.ReviewDetailsResponse{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully created review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailHandleGrpc) Update(ctx context.Context, request *pb.UpdateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetReviewDetailId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &id,
		Type:           request.GetType(),
		Url:            request.GetUrl(),
		Caption:        request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		return nil, reviewdetail_errors.ErrGrpcValidateUpdateReviewDetail
	}

	reviewDetail, err := s.reviewDetailService.UpdateReviewDetail(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetail := &pb.ReviewDetailsResponse{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully updated review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.reviewDetailService.TrashedReviewDetail(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if reviewDetail.DeletedAt.Valid {
		deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
	}

	protoReviewDetail := &pb.ReviewDetailsResponseDeleteAt{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully trashed review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.reviewDetailService.RestoreReviewDetail(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if reviewDetail.DeletedAt.Valid {
		deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
	}

	protoReviewDetail := &pb.ReviewDetailsResponseDeleteAt{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully restored review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	_, err := s.reviewDetailService.DeleteReviewDetailPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewDelete{
		Status:  "success",
		Message: "Successfully deleted review detail permanently",
	}, nil
}

func (s *reviewDetailHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailService.RestoreAllReviewDetail(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully restored all review details",
	}, nil
}

func (s *reviewDetailHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailService.DeleteAllReviewDetailPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully deleted all review details permanently",
	}, nil
}
