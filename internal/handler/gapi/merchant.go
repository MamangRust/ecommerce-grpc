package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantService service.MerchantService
}

func NewMerchantHandleGrpc(
	merchantService service.MerchantService,
) *merchantHandleGrpc {
	return &merchantHandleGrpc{
		merchantService: merchantService,
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

	merchants, totalRecords, err := s.merchantService.FindAllMerchants(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantResponse
	for _, m := range merchants {
		pbMerchants = append(pbMerchants, &pb.MerchantResponse{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  *m.Description,
			Address:      *m.Address,
			ContactEmail: *m.ContactEmail,
			ContactPhone: *m.ContactPhone,
			Status:       m.Status,
			CreatedAt:    m.CreatedAt.Time.String(),
			UpdatedAt:    m.UpdatedAt.Time.String(),
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantResponse{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  *merchant.Description,
		Address:      *merchant.Address,
		ContactEmail: *merchant.ContactEmail,
		ContactPhone: *merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt.Time.String(),
		UpdatedAt:    merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    pbMerchant,
	}, nil
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

	// Tambahkan ctx ke pemanggilan service
	merchants, totalRecords, err := s.merchantService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantResponseDeleteAt
	for _, m := range merchants {
		var deletedAt string

		if m.DeletedAt.Valid {
			deletedAt = m.DeletedAt.Time.Format("2006-01-02")
		}

		pbMerchants = append(pbMerchants, &pb.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  *m.Description,
			Address:      *m.Address,
			ContactEmail: *m.ContactEmail,
			ContactPhone: *m.ContactPhone,
			Status:       m.Status,
			CreatedAt:    m.CreatedAt.Time.String(),
			UpdatedAt:    m.UpdatedAt.Time.String(),
			DeletedAt:    &wrapperspb.StringValue{Value: deletedAt},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
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

	merchants, totalRecords, err := s.merchantService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var pbMerchants []*pb.MerchantResponseDeleteAt
	for _, m := range merchants {
		var deletedAt string

		if m.DeletedAt.Valid {
			deletedAt = m.DeletedAt.Time.Format("2006-01-02")
		}

		pbMerchants = append(pbMerchants, &pb.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  *m.Description,
			Address:      *m.Address,
			ContactEmail: *m.ContactEmail,
			ContactPhone: *m.ContactPhone,
			Status:       m.Status,
			CreatedAt:    m.CreatedAt.Time.String(),
			UpdatedAt:    m.UpdatedAt.Time.String(),
			DeletedAt:    &wrapperspb.StringValue{Value: deletedAt},
		})
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
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

	merchant, err := s.merchantService.CreateMerchant(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantResponse{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  *merchant.Description,
		Address:      *merchant.Address,
		ContactEmail: *merchant.ContactEmail,
		ContactPhone: *merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt.Time.String(),
		UpdatedAt:    merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    pbMerchant,
	}, nil
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

	merchant, err := s.merchantService.UpdateMerchant(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantResponse{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  *merchant.Description,
		Address:      *merchant.Address,
		ContactEmail: *merchant.ContactEmail,
		ContactPhone: *merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt.Time.String(),
		UpdatedAt:    merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.TrashedMerchant(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantResponseDeleteAt{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  *merchant.Description,
		Address:      *merchant.Address,
		ContactEmail: *merchant.ContactEmail,
		ContactPhone: *merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt.Time.String(),
		UpdatedAt:    merchant.UpdatedAt.Time.String(),
		DeletedAt:    &wrapperspb.StringValue{Value: merchant.DeletedAt.Time.Format("2006-01-02")},
	}

	return &pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantService.RestoreMerchant(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchant := &pb.MerchantResponseDeleteAt{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  *merchant.Description,
		Address:      *merchant.Address,
		ContactEmail: *merchant.ContactEmail,
		ContactPhone: *merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt.Time.String(),
		UpdatedAt:    merchant.UpdatedAt.Time.String(),
	}

	return &pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    pbMerchant,
	}, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	_, err := s.merchantService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.RestoreAllMerchant(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantService.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete all merchant permanently",
	}, nil
}
