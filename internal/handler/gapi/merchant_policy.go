package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	merchantpolicy_errors "ecommerce/pkg/errors/merchant_policy_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPolicyHandleGrpc struct {
	pb.UnimplementedMerchantPoliciesServiceServer
	merchantPolicyService service.MerchantPoliciesService
	mapping               protomapper.MerchantPolicyProtoMapper
	mappingMerchant       protomapper.MerchantProtoMapper
}

func NewMerchantPolicyHandleGrpc(
	merchantPolicyService service.MerchantPoliciesService,
	mapping protomapper.MerchantPolicyProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) *merchantPolicyHandleGrpc {
	return &merchantPolicyHandleGrpc{
		merchantPolicyService: merchantPolicyService,
		mapping:               mapping,
		mappingMerchant:       mappingMerchant,
	}
}

func (s *merchantPolicyHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPolicies, error) {
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

	merchant, totalRecords, err := s.merchantPolicyService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantPolicy(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	merchant, err := s.merchantPolicyService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantPolicyHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
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

	merchant, totalRecords, err := s.merchantPolicyService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantPolicyDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

	return so, nil
}

func (s *merchantPolicyHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
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

	users, totalRecords, err := s.merchantPolicyService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantPolicyDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantPolicyHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	req := &requests.CreateMerchantPolicyRequest{
		MerchantID:  int(request.GetMerchantId()),
		PolicyType:  request.GetPolicyType(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantpolicy_errors.ErrGrpcValidateCreateMerchantPolicy
	}

	merchant, err := s.merchantPolicyService.CreateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully created merchant policy", merchant)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(request.GetMerchantPolicyId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	req := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &id,
		PolicyType:       request.GetPolicyType(),
		Title:            request.GetTitle(),
		Description:      request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantpolicy_errors.ErrGrpcValidateUpdateMerchantPolicy
	}

	merchant, err := s.merchantPolicyService.UpdateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully updated merchant policy", merchant)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	merchant, err := s.merchantPolicyService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantPolicyDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantPolicyHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	merchant, err := s.merchantPolicyService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantPolicyDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantPolicyHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	_, err := s.merchantPolicyService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantPolicyHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantPolicyService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantPolicyHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantPolicyService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
