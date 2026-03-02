package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/category_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type categoryHandleGrpc struct {
	pb.UnimplementedCategoryServiceServer
	categoryService service.CategoryService
}

func NewCategoryHandleGrpc(
	categoryService service.CategoryService,
) *categoryHandleGrpc {
	return &categoryHandleGrpc{
		categoryService: categoryService,
	}
}

func (s *categoryHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategories := make([]*pb.CategoryResponse, len(categories))
	for i, category := range categories {
		protoCategories[i] = &pb.CategoryResponse{
			Id:            int32(category.CategoryID),
			Name:          category.Name,
			Description:   *category.Description,
			SlugCategory:  *category.SlugCategory,
			ImageCategory: *category.ImageCategory,
			CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       protoCategories,
		Pagination: paginationMeta,
	}, nil
}

func (s *categoryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategory := &pb.CategoryResponse{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   *category.Description,
		SlugCategory:  *category.SlugCategory,
		ImageCategory: *category.ImageCategory,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    protoCategory,
	}, nil
}

func (s *categoryHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategories := make([]*pb.CategoryResponseDeleteAt, len(categories))
	for i, category := range categories {
		var deletedAt string
		if category.DeletedAt.Valid {
			deletedAt = category.DeletedAt.Time.Format("2006-01-02")
		}

		protoCategories[i] = &pb.CategoryResponseDeleteAt{
			Id:            int32(category.CategoryID),
			Name:          category.Name,
			Description:   *category.Description,
			SlugCategory:  *category.SlugCategory,
			ImageCategory: *category.ImageCategory,
			CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     category.UpdatedAt.Time.Format("2006-0-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       protoCategories,
		Pagination: paginationMeta,
	}, nil
}

func (s *categoryHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategories := make([]*pb.CategoryResponseDeleteAt, len(categories))
	for i, category := range categories {
		var deletedAt string
		if category.DeletedAt.Valid {
			deletedAt = category.DeletedAt.Time.Format("2006-01-02")
		}

		protoCategories[i] = &pb.CategoryResponseDeleteAt{
			Id:            int32(category.CategoryID),
			Name:          category.Name,
			Description:   *category.Description,
			SlugCategory:  *category.SlugCategory,
			ImageCategory: *category.ImageCategory,
			CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       protoCategories,
		Pagination: paginationMeta,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	serviceResults, err := s.categoryService.FindMonthlyTotalPrice(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesMonthlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesMonthlyTotalPriceResponse{
			Year:         result.Year,
			Month:        result.Month,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := s.categoryService.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesYearlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesYearlyTotalPriceResponse{
			Year:         result.Year,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	serviceResults, err := s.categoryService.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesMonthlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesMonthlyTotalPriceResponse{
			Year:         result.Year,
			Month:        result.Month,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := s.categoryService.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesYearlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesYearlyTotalPriceResponse{
			Year:         result.Year,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalPriceByMerchant) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	serviceResults, err := s.categoryService.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesMonthlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesMonthlyTotalPriceResponse{
			Year:         result.Year,
			Month:        result.Month,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearTotalPriceByMerchant) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := s.categoryService.FindYearlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoriesYearlyTotalPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoriesYearlyTotalPriceResponse{
			Year:         result.Year,
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := s.categoryService.FindMonthPrice(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryMonthPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryMonthPriceResponse{
			Month:        result.Month,
			CategoryId:   int32(result.CategoryID),
			CategoryName: result.CategoryName,
			OrderCount:   int32(result.OrderCount),
			ItemsSold:    int32(result.ItemsSold),
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := s.categoryService.FindYearPrice(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryYearPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryYearPriceResponse{
			Year:               result.Year,
			CategoryId:         int32(result.CategoryID),
			CategoryName:       result.CategoryName,
			OrderCount:         int32(result.OrderCount),
			ItemsSold:          int32(result.ItemsSold),
			TotalRevenue:       int32(result.TotalRevenue),
			UniqueProductsSold: int32(result.UniqueProductsSold),
		})
	}

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := s.categoryService.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryMonthPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryMonthPriceResponse{
			Month:        result.Month,
			CategoryId:   int32(result.CategoryID),
			CategoryName: result.CategoryName,
			OrderCount:   int32(result.OrderCount),
			ItemsSold:    int32(result.ItemsSold),
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := s.categoryService.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryYearPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryYearPriceResponse{
			Year:               result.Year,
			CategoryId:         int32(result.CategoryID),
			CategoryName:       result.CategoryName,
			OrderCount:         int32(result.OrderCount),
			ItemsSold:          int32(result.ItemsSold),
			TotalRevenue:       int32(result.TotalRevenue),
			UniqueProductsSold: int32(result.UniqueProductsSold),
		})
	}

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := s.categoryService.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryMonthPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryMonthPriceResponse{
			Month:        result.Month,
			CategoryId:   int32(result.CategoryID),
			CategoryName: result.CategoryName,
			OrderCount:   int32(result.OrderCount),
			ItemsSold:    int32(result.ItemsSold),
			TotalRevenue: int32(result.TotalRevenue),
		})
	}

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Category monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := s.categoryService.FindYearPriceById(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pb.CategoryYearPriceResponse
	for _, result := range serviceResults {
		data = append(data, &pb.CategoryYearPriceResponse{
			Year:               result.Year,
			CategoryId:         int32(result.CategoryID),
			CategoryName:       result.CategoryName,
			OrderCount:         int32(result.OrderCount),
			ItemsSold:          int32(result.ItemsSold),
			TotalRevenue:       int32(result.TotalRevenue),
			UniqueProductsSold: int32(result.UniqueProductsSold),
		})
	}

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Category yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *categoryHandleGrpc) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
	req := &requests.CreateCategoryRequest{
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryService.CreateCategory(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategory := &pb.CategoryResponse{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   *category.Description,
		SlugCategory:  *category.SlugCategory,
		ImageCategory: *category.ImageCategory,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    protoCategory,
	}, nil
}

func (s *categoryHandleGrpc) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetCategoryId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:    &id,
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryService.UpdateCategory(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCategory := &pb.CategoryResponse{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   *category.Description,
		SlugCategory:  *category.SlugCategory,
		ImageCategory: *category.ImageCategory,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    protoCategory,
	}, nil
}

func (s *categoryHandleGrpc) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryService.TrashedCategory(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if category.DeletedAt.Valid {
		deletedAt = category.DeletedAt.Time.Format("2006-01-02")
	}

	protoCategory := &pb.CategoryResponseDeleteAt{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   *category.Description,
		SlugCategory:  *category.SlugCategory,
		ImageCategory: *category.ImageCategory,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-0-01-02"),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    protoCategory,
	}, nil
}

func (s *categoryHandleGrpc) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryService.RestoreCategory(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if category.DeletedAt.Valid {
		deletedAt = category.DeletedAt.Time.Format("2006-01-02")
	}

	protoCategory := &pb.CategoryResponseDeleteAt{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   *category.Description,
		SlugCategory:  *category.SlugCategory,
		ImageCategory: *category.ImageCategory,
		CreatedAt:     category.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     category.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    protoCategory,
	}, nil
}

func (s *categoryHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	_, err := s.categoryService.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (s *categoryHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := s.categoryService.RestoreAllCategories(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restored all categories",
	}, nil
}

func (s *categoryHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := s.categoryService.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully deleted all categories permanently",
	}, nil
}
