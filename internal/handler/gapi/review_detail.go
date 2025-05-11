package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	reviewdetail_errors "ecommerce/pkg/errors/review_detail"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailHandleGrpc struct {
	pb.UnimplementedReviewDetailServiceServer
	reviewDetailService service.ReviewDetailService
	mapping             protomapper.ReviewDetailProtoMapper
	mappingReview       protomapper.ReviewProtoMapper
}

func NewReviewDetailHandleGrpc(
	reviewDetailService service.ReviewDetailService,
	mapping protomapper.ReviewDetailProtoMapper,
	mappingReview protomapper.ReviewProtoMapper,
) *reviewDetailHandleGrpc {
	return &reviewDetailHandleGrpc{
		reviewDetailService: reviewDetailService,
		mapping:             mapping,
		mappingReview:       mappingReview,
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

	Review, totalRecords, err := s.reviewDetailService.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewDetail(paginationMeta, "success", "Successfully fetched Review", Review)
	return so, nil
}

func (s *reviewDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	Review, err := s.reviewDetailService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully fetched Review", Review)

	return so, nil

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

	Review, totalRecords, err := s.reviewDetailService.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewDetailDeleteAt(paginationMeta, "success", "Successfully fetched active Review", Review)

	return so, nil
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

	users, totalRecords, err := s.reviewDetailService.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewDetailDeleteAt(paginationMeta, "success", "Successfully fetched trashed Review", users)

	return so, nil
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

	review, err := s.reviewDetailService.CreateReviewDetail(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully created Review Detail", review)
	return so, nil
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

	review, err := s.reviewDetailService.UpdateReviewDetail(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully updated Review Detail", review)
	return so, nil
}

func (s *reviewDetailHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	Review, err := s.reviewDetailService.TrashedReviewDetail(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDetailDeleteAt("success", "Successfully trashed Review", Review)

	return so, nil
}

func (s *reviewDetailHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	Review, err := s.reviewDetailService.RestoreReviewDetail(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDetailDeleteAt("success", "Successfully restored Review", Review)

	return so, nil
}

func (s *reviewDetailHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	_, err := s.reviewDetailService.DeleteReviewDetailPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingReview.ToProtoResponseReviewDelete("success", "Successfully deleted Review permanently")

	return so, nil
}

func (s *reviewDetailHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailService.RestoreAllReviewDetail()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingReview.ToProtoResponseReviewAll("success", "Successfully restore all Review")

	return so, nil
}

func (s *reviewDetailHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailService.DeleteAllReviewDetailPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingReview.ToProtoResponseReviewAll("success", "Successfully delete Review permanen")

	return so, nil
}
