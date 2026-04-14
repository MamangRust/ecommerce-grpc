package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	merchantbusiness_errors "ecommerce/pkg/errors/merchant_business"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantBusinessHandleGrpc struct {
	pb.UnimplementedMerchantBusinessServiceServer
	merchantBusinessService service.MerchantBusinessService
}

func NewMerchantBusinessHandleGrpc(
	merchantBusinessService service.MerchantBusinessService,
) *merchantBusinessHandleGrpc {
	return &merchantBusinessHandleGrpc{
		merchantBusinessService: merchantBusinessService,
	}
}

func (s *merchantBusinessHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusiness, error) {
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

	merchants, totalRecords, err := s.merchantBusinessService.FindAllMerchants(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantBusinessResponse
	for _, merchant := range merchants {
		pbMerchants = append(pbMerchants, &pb.MerchantBusinessResponse{
			Id:                int32(merchant.MerchantBusinessInfoID),
			MerchantId:        int32(merchant.MerchantID),
			BusinessType:      *merchant.BusinessType,
			TaxId:             *merchant.TaxID,
			EstablishedYear:   int32(*merchant.EstablishedYear),
			NumberOfEmployees: int32(*merchant.NumberOfEmployees),
			WebsiteUrl:        *merchant.WebsiteUrl,
			CreatedAt:         merchant.CreatedAt.Time.String(),
			UpdatedAt:         merchant.UpdatedAt.Time.String(),
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantBusiness{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantBusinessHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantBusinessResponse{
		Id:                int32(merchant.MerchantBusinessInfoID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      *merchant.BusinessType,
		TaxId:             *merchant.TaxID,
		EstablishedYear:   int32(*merchant.EstablishedYear),
		NumberOfEmployees: int32(*merchant.NumberOfEmployees),
		WebsiteUrl:        *merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt.Time.String(),
		UpdatedAt:         merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantBusinessHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
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

	merchants, totalRecords, err := s.merchantBusinessService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantBusinessResponseDeleteAt
	for _, merchant := range merchants {
		pbMerchants = append(pbMerchants, &pb.MerchantBusinessResponseDeleteAt{
			Id:                int32(merchant.MerchantBusinessInfoID),
			MerchantId:        int32(merchant.MerchantID),
			BusinessType:      *merchant.BusinessType,
			TaxId:             *merchant.TaxID,
			EstablishedYear:   int32(*merchant.EstablishedYear),
			NumberOfEmployees: int32(*merchant.NumberOfEmployees),
			WebsiteUrl:        *merchant.WebsiteUrl,
			CreatedAt:         merchant.CreatedAt.Time.String(),
			UpdatedAt:         merchant.UpdatedAt.Time.String(),
			DeletedAt:         &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantBusinessHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
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

	merchants, totalRecords, err := s.merchantBusinessService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbUsers []*pb.MerchantBusinessResponseDeleteAt
	for _, merchant := range merchants {
		pbUsers = append(pbUsers, &pb.MerchantBusinessResponseDeleteAt{
			Id:                int32(merchant.MerchantBusinessInfoID),
			MerchantId:        int32(merchant.MerchantID),
			BusinessType:      *merchant.BusinessType,
			TaxId:             *merchant.TaxID,
			EstablishedYear:   int32(*merchant.EstablishedYear),
			NumberOfEmployees: int32(*merchant.NumberOfEmployees),
			WebsiteUrl:        *merchant.WebsiteUrl,
			CreatedAt:         merchant.CreatedAt.Time.String(),
			UpdatedAt:         merchant.UpdatedAt.Time.String(),
			DeletedAt:         &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       pbUsers,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantBusinessHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        int(request.GetMerchantId()),
		BusinessType:      request.GetBusinessType(),
		TaxID:             request.GetTaxId(),
		EstablishedYear:   int(request.GetEstablishedYear()),
		NumberOfEmployees: int(request.GetNumberOfEmployees()),
		WebsiteUrl:        request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateCreateMerchantBusiness
	}

	merchant, err := s.merchantBusinessService.CreateMerchantBusiness(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantBusinessResponse{
		Id:                int32(merchant.MerchantBusinessInfoID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      *merchant.BusinessType,
		TaxId:             *merchant.TaxID,
		EstablishedYear:   int32(*merchant.EstablishedYear),
		NumberOfEmployees: int32(*merchant.NumberOfEmployees),
		WebsiteUrl:        *merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt.Time.String(),
		UpdatedAt:         merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully created merchant business information",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantBusinessHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetMerchantBusinessInfoId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	req := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &id,
		BusinessType:           request.GetBusinessType(),
		TaxID:                  request.GetTaxId(),
		EstablishedYear:        int(request.GetEstablishedYear()),
		NumberOfEmployees:      int(request.GetNumberOfEmployees()),
		WebsiteUrl:             request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateUpdateMerchantBusiness
	}

	merchant, err := s.merchantBusinessService.UpdateMerchantBusiness(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantBusinessResponse{
		Id:                int32(merchant.MerchantBusinessInfoID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      *merchant.BusinessType,
		TaxId:             *merchant.TaxID,
		EstablishedYear:   int32(*merchant.EstablishedYear),
		NumberOfEmployees: int32(*merchant.NumberOfEmployees),
		WebsiteUrl:        *merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt.Time.String(),
		UpdatedAt:         merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully updated merchant business information",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantBusinessHandleGrpc) TrashedMerchantBusiness(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessService.TrashedMerchantBusiness(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantBusinessResponseDeleteAt{
		Id:                int32(merchant.MerchantBusinessInfoID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      *merchant.BusinessType,
		TaxId:             *merchant.TaxID,
		EstablishedYear:   int32(*merchant.EstablishedYear),
		NumberOfEmployees: int32(*merchant.NumberOfEmployees),
		WebsiteUrl:        *merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt.Time.String(),
		UpdatedAt:         merchant.UpdatedAt.Time.String(),
		DeletedAt:         &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantBusinessDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantBusinessHandleGrpc) RestoreMerchantBusiness(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessService.RestoreMerchantBusiness(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantBusinessResponseDeleteAt{
		Id:                int32(merchant.MerchantBusinessInfoID),
		MerchantId:        int32(merchant.MerchantID),
		BusinessType:      *merchant.BusinessType,
		TaxId:             *merchant.TaxID,
		EstablishedYear:   int32(*merchant.EstablishedYear),
		NumberOfEmployees: int32(*merchant.NumberOfEmployees),
		WebsiteUrl:        *merchant.WebsiteUrl,
		CreatedAt:         merchant.CreatedAt.Time.String(),
		UpdatedAt:         merchant.UpdatedAt.Time.String(),
		DeletedAt:         &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantBusinessDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantBusinessHandleGrpc) DeleteMerchantBusinessPermanent(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	_, err := s.merchantBusinessService.DeleteMerchantBusinessPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantBusinessHandleGrpc) RestoreAllMerchantBusiness(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessService.RestoreAllMerchantBusiness(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantBusinessHandleGrpc) DeleteAllMerchantBusinessPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessService.DeleteAllMerchantBusinessPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant permanently",
	}, nil
}
