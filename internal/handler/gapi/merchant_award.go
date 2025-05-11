package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardHandleGrpc struct {
	pb.UnimplementedMerchantAwardServiceServer
	merchantAwardService service.MerchantAwardService
	mapping              protomapper.MerchantAwardProtoMapper
	mappingMerchant      protomapper.MerchantProtoMapper
}

func NewMerchantAwardHandleGrpc(
	merchantAwardService service.MerchantAwardService,
	mapping protomapper.MerchantAwardProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) *merchantAwardHandleGrpc {
	return &merchantAwardHandleGrpc{
		merchantAwardService: merchantAwardService,
		mapping:              mapping,
		mappingMerchant:      mappingMerchant,
	}
}

func (s *merchantAwardHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAward, error) {
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

	merchant, totalRecords, err := s.merchantAwardService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAward(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantAwardHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
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

	merchant, totalRecords, err := s.merchantAwardService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

	return so, nil
}

func (s *merchantAwardHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
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

	users, totalRecords, err := s.merchantAwardService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantAwardHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     int(request.GetMerchantId()),
		Title:          request.GetTitle(),
		Description:    request.GetDescription(),
		IssuedBy:       request.GetIssuedBy(),
		IssueDate:      request.GetIssueDate(),
		ExpiryDate:     request.GetExpiryDate(),
		CertificateUrl: request.GetCertificateUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateCreateMerchantAward
	}

	merchant, err := s.merchantAwardService.CreateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully created merchant award", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetMerchantCertificationId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	req := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &id,
		Title:                   request.GetTitle(),
		Description:             request.GetDescription(),
		IssuedBy:                request.GetIssuedBy(),
		IssueDate:               request.GetIssueDate(),
		ExpiryDate:              request.GetExpiryDate(),
		CertificateUrl:          request.GetCertificateUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateUpdateMerchantAward
	}

	merchant, err := s.merchantAwardService.UpdateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully updated merchant award", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	_, err := s.merchantAwardService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
