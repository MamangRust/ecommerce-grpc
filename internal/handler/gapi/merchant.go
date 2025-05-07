package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantService service.MerchantService
	mapping         protomapper.MerchantProtoMapper
}

func NewMerchantHandleGrpc(
	merchantService service.MerchantService,
	mapping protomapper.MerchantProtoMapper,
) *merchantHandleGrpc {
	return &merchantHandleGrpc{
		merchantService: merchantService,
		mapping:         mapping,
	}
}

func (s *merchantHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.merchantService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchant(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.merchantService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.merchantService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantService.CreateMerchant(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantService.UpdateMerchant(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	_, err := s.merchantService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
