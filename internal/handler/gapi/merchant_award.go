package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantAwardHandleGrpc struct {
	pb.UnimplementedMerchantAwardServiceServer
	merchantAwardService service.MerchantAwardService
}

func NewMerchantAwardHandleGrpc(
	merchantAwardService service.MerchantAwardService,
) *merchantAwardHandleGrpc {
	return &merchantAwardHandleGrpc{
		merchantAwardService: merchantAwardService,
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

	merchants, totalRecords, err := s.merchantAwardService.FindAllMerchants(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantAwardResponse
	for _, merchant := range merchants {
		pbMerchants = append(pbMerchants, &pb.MerchantAwardResponse{
			Id:             int32(merchant.MerchantCertificationID),
			MerchantId:     int32(merchant.MerchantID),
			Title:          merchant.Title,
			Description:    *merchant.Description,
			IssuedBy:       *merchant.IssuedBy,
			IssueDate:      merchant.IssueDate.Time.String(),
			ExpiryDate:     merchant.ExpiryDate.Time.String(),
			CertificateUrl: *merchant.CertificateUrl,
			CreatedAt:      merchant.CreatedAt.Time.String(),
			UpdatedAt:      merchant.UpdatedAt.Time.String(),
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantAward{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantAwardHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantAwardResponse{
		Id:             int32(merchant.MerchantCertificationID),
		MerchantId:     int32(merchant.MerchantID),
		Title:          merchant.Title,
		Description:    *merchant.Description,
		IssuedBy:       *merchant.IssuedBy,
		IssueDate:      merchant.IssueDate.Time.String(),
		ExpiryDate:     merchant.ExpiryDate.Time.String(),
		CertificateUrl: *merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt.Time.String(),
		UpdatedAt:      merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    pbMerchant,
	}, nil
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

	merchants, totalRecords, err := s.merchantAwardService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantAwardResponseDeleteAt
	for _, merchant := range merchants {
		pbMerchants = append(pbMerchants, &pb.MerchantAwardResponseDeleteAt{
			Id:             int32(merchant.MerchantCertificationID),
			MerchantId:     int32(merchant.MerchantID),
			Title:          merchant.Title,
			Description:    *merchant.Description,
			IssuedBy:       *merchant.IssuedBy,
			IssueDate:      merchant.IssueDate.Time.String(),
			ExpiryDate:     merchant.ExpiryDate.Time.String(),
			CertificateUrl: *merchant.CertificateUrl,
			CreatedAt:      merchant.CreatedAt.Time.String(),
			UpdatedAt:      merchant.UpdatedAt.Time.String(),
			DeletedAt:      &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
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

	users, totalRecords, err := s.merchantAwardService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbUsers []*pb.MerchantAwardResponseDeleteAt
	for _, merchant := range users {
		pbUsers = append(pbUsers, &pb.MerchantAwardResponseDeleteAt{
			Id:             int32(merchant.MerchantCertificationID),
			MerchantId:     int32(merchant.MerchantID),
			Title:          merchant.Title,
			Description:    *merchant.Description,
			IssuedBy:       *merchant.IssuedBy,
			IssueDate:      merchant.IssueDate.Time.String(),
			ExpiryDate:     merchant.ExpiryDate.Time.String(),
			CertificateUrl: *merchant.CertificateUrl,
			CreatedAt:      merchant.CreatedAt.Time.String(),
			UpdatedAt:      merchant.UpdatedAt.Time.String(),
			DeletedAt:      &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       pbUsers,
		Pagination: paginationMeta,
	}, nil
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

	merchant, err := s.merchantAwardService.CreateMerchantAward(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantAwardResponse{
		Id:             int32(merchant.MerchantCertificationID),
		MerchantId:     int32(merchant.MerchantID),
		Title:          merchant.Title,
		Description:    *merchant.Description,
		IssuedBy:       *merchant.IssuedBy,
		IssueDate:      merchant.IssueDate.Time.String(),
		ExpiryDate:     merchant.ExpiryDate.Time.String(),
		CertificateUrl: *merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt.Time.String(),
		UpdatedAt:      merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully created merchant award",
		Data:    pbMerchant,
	}, nil
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

	merchant, err := s.merchantAwardService.UpdateMerchantAward(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantAwardResponse{
		Id:             int32(merchant.MerchantCertificationID),
		MerchantId:     int32(merchant.MerchantID),
		Title:          merchant.Title,
		Description:    *merchant.Description,
		IssuedBy:       *merchant.IssuedBy,
		IssueDate:      merchant.IssueDate.Time.String(),
		ExpiryDate:     merchant.ExpiryDate.Time.String(),
		CertificateUrl: *merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt.Time.String(),
		UpdatedAt:      merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully updated merchant award",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantAwardHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.TrashedMerchantAward(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantAwardResponseDeleteAt{
		Id:             int32(merchant.MerchantCertificationID),
		MerchantId:     int32(merchant.MerchantID),
		Title:          merchant.Title,
		Description:    *merchant.Description,
		IssuedBy:       *merchant.IssuedBy,
		IssueDate:      merchant.IssueDate.Time.String(),
		ExpiryDate:     merchant.ExpiryDate.Time.String(),
		CertificateUrl: *merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt.Time.String(),
		UpdatedAt:      merchant.UpdatedAt.Time.String(),
		DeletedAt:      &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantAwardHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardService.RestoreMerchantAward(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantAwardResponseDeleteAt{
		Id:             int32(merchant.MerchantCertificationID),
		MerchantId:     int32(merchant.MerchantID),
		Title:          merchant.Title,
		Description:    *merchant.Description,
		IssuedBy:       *merchant.IssuedBy,
		IssueDate:      merchant.IssueDate.Time.String(),
		ExpiryDate:     merchant.ExpiryDate.Time.String(),
		CertificateUrl: *merchant.CertificateUrl,
		CreatedAt:      merchant.CreatedAt.Time.String(),
		UpdatedAt:      merchant.UpdatedAt.Time.String(),
		DeletedAt:      &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.String()},
	}

	return &pb.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantAwardHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	_, err := s.merchantAwardService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantAwardHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardService.RestoreAllMerchantAward(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantAwardHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardService.DeleteAllMerchantAwardPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant permanently",
	}, nil
}
