package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	merchantpolicy_errors "ecommerce/pkg/errors/merchant_policy_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantPolicyHandleGrpc struct {
	pb.UnimplementedMerchantPoliciesServiceServer
	merchantPolicyService service.MerchantPoliciesService
}

func NewMerchantPolicyHandleGrpc(
	merchantPolicyService service.MerchantPoliciesService,
) *merchantPolicyHandleGrpc {
	return &merchantPolicyHandleGrpc{
		merchantPolicyService: merchantPolicyService,
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

	policies, totalRecords, err := s.merchantPolicyService.FindAllMerchantPolicy(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbPolicies []*pb.MerchantPoliciesResponse
	for _, p := range policies {
		pbPolicies = append(pbPolicies, &pb.MerchantPoliciesResponse{
			Id:           int32(p.MerchantPolicyID),
			MerchantId:   int32(p.MerchantID),
			PolicyType:   p.PolicyType,
			Title:        p.Title,
			Description:  p.Description,
			CreatedAt:    p.CreatedAt.Time.String(),
			UpdatedAt:    p.UpdatedAt.Time.String(),
			MerchantName: p.MerchantName,
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantPolicies{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       pbPolicies,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantPolicyHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	policy, err := s.merchantPolicyService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbPolicy := &pb.MerchantPoliciesResponse{
		Id:          int32(policy.MerchantPolicyID),
		MerchantId:  int32(policy.MerchantID),
		PolicyType:  policy.PolicyType,
		Title:       policy.Title,
		Description: policy.Description,
		CreatedAt:   policy.CreatedAt.Time.String(),
		UpdatedAt:   policy.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    pbPolicy,
	}, nil
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

	policies, totalRecords, err := s.merchantPolicyService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbPolicies []*pb.MerchantPoliciesResponseDeleteAt
	for _, p := range policies {
		pbPolicies = append(pbPolicies, &pb.MerchantPoliciesResponseDeleteAt{
			Id:           int32(p.MerchantID),
			MerchantId:   int32(p.MerchantID),
			PolicyType:   p.PolicyType,
			Title:        p.Title,
			Description:  p.Description,
			CreatedAt:    p.CreatedAt.Time.String(),
			UpdatedAt:    p.UpdatedAt.Time.String(),
			MerchantName: p.MerchantName,
			DeletedAt:    &wrapperspb.StringValue{Value: p.DeletedAt.Time.String()},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       pbPolicies,
		Pagination: paginationMeta,
	}, nil
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

	policies, totalRecords, err := s.merchantPolicyService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbPolicies []*pb.MerchantPoliciesResponseDeleteAt
	for _, p := range policies {
		pbPolicies = append(pbPolicies, &pb.MerchantPoliciesResponseDeleteAt{
			Id:           int32(p.MerchantPolicyID),
			MerchantId:   int32(p.MerchantID),
			PolicyType:   p.PolicyType,
			Title:        p.Title,
			Description:  p.Description,
			CreatedAt:    p.CreatedAt.Time.String(),
			UpdatedAt:    p.UpdatedAt.Time.String(),
			DeletedAt:    &wrapperspb.StringValue{Value: p.DeletedAt.Time.String()},
			MerchantName: p.MerchantName,
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       pbPolicies,
		Pagination: paginationMeta,
	}, nil
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

	policy, err := s.merchantPolicyService.CreateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbPolicy := &pb.MerchantPoliciesResponse{
		Id:          int32(policy.MerchantPolicyID),
		MerchantId:  int32(policy.MerchantID),
		PolicyType:  policy.PolicyType,
		Title:       policy.Title,
		Description: policy.Description,
		CreatedAt:   policy.CreatedAt.Time.String(),
		UpdatedAt:   policy.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully created merchant policy",
		Data:    pbPolicy,
	}, nil
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

	policy, err := s.merchantPolicyService.UpdateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbPolicy := &pb.MerchantPoliciesResponse{
		Id:          int32(policy.MerchantPolicyID),
		MerchantId:  int32(policy.MerchantID),
		PolicyType:  policy.PolicyType,
		Title:       policy.Title,
		Description: policy.Description,
		CreatedAt:   policy.CreatedAt.Time.String(),
		UpdatedAt:   policy.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully updated merchant policy",
		Data:    pbPolicy,
	}, nil
}

func (s *merchantPolicyHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	policy, err := s.merchantPolicyService.TrashedMerchantPolicy(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbPolicy := &pb.MerchantPoliciesResponseDeleteAt{
		Id:          int32(policy.MerchantPolicyID),
		MerchantId:  int32(policy.MerchantID),
		PolicyType:  policy.PolicyType,
		Title:       policy.Title,
		Description: policy.Description,
		CreatedAt:   policy.CreatedAt.Time.String(),
		UpdatedAt:   policy.UpdatedAt.Time.String(),
		DeletedAt:   &wrapperspb.StringValue{Value: policy.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    pbPolicy,
	}, nil
}

func (s *merchantPolicyHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	policy, err := s.merchantPolicyService.RestoreMerchantPolicy(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbPolicy := &pb.MerchantPoliciesResponseDeleteAt{
		Id:          int32(policy.MerchantPolicyID),
		MerchantId:  int32(policy.MerchantID),
		PolicyType:  policy.PolicyType,
		Title:       policy.Title,
		Description: policy.Description,
		CreatedAt:   policy.CreatedAt.Time.String(),
		UpdatedAt:   policy.UpdatedAt.Time.String(),
		DeletedAt:   &wrapperspb.StringValue{Value: policy.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    pbPolicy,
	}, nil
}

func (s *merchantPolicyHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	_, err := s.merchantPolicyService.DeleteMerchantPolicyPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantPolicyHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantPolicyService.RestoreAllMerchantPolicy(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantPolicyHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantPolicyService.DeleteAllMerchantPolicyPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant permanently",
	}, nil
}
