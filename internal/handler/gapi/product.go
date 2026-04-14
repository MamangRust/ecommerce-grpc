package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/utils"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type productHandleGrpc struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewProductHandleGrpc(
	productService service.ProductService,
) *productHandleGrpc {
	return &productHandleGrpc{
		productService: productService,
	}
}

func (s *productHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productService.FindAllProducts(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		protoProducts[i] = &pb.ProductResponse{
			Id:           int32(product.ProductID),
			MerchantId:   int32(product.MerchantID),
			CategoryId:   int32(product.CategoryID),
			Name:         product.Name,
			Description:  utils.StringPtrToString(product.Description),
			Price:        int32(product.Price),
			CountInStock: int32(product.CountInStock),
			Brand:        utils.StringPtrToString(product.Brand),
			Weight:       utils.Int32PtrToInt32(product.Weight),
			Rating:       utils.Float64PtrToFloat32(product.Rating),
			SlugProduct:  utils.StringPtrToString(product.SlugProduct),
			ImageProduct: utils.StringPtrToString(product.ImageProduct),
			CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched products",
		Data:       protoProducts,
		Pagination: paginationMeta,
	}, nil
}

func (s *productHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllProductMerchantRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())
	min_price := int(request.GetMinPrice())
	max_price := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if min_price <= 0 {
		min_price = 0
	}

	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.FindAllProductByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MinPrice:   &min_price,
		MaxPrice:   &max_price,
	}

	products, totalRecords, err := s.productService.FindByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		protoProducts[i] = &pb.ProductResponse{
			Id:           int32(product.ProductID),
			MerchantId:   int32(product.MerchantID),
			CategoryId:   int32(product.CategoryID),
			Name:         product.Name,
			Description:  utils.StringPtrToString(product.Description),
			Price:        int32(product.Price),
			CountInStock: int32(product.CountInStock),
			Brand:        utils.StringPtrToString(product.Brand),
			Weight:       utils.Int32PtrToInt32(product.Weight),
			Rating:       utils.Float64PtrToFloat32(product.Rating),
			SlugProduct:  utils.StringPtrToString(product.SlugProduct),
			ImageProduct: utils.StringPtrToString(product.ImageProduct),
			CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched merchant products",
		Data:       protoProducts,
		Pagination: paginationMeta,
	}, nil
}

func (s *productHandleGrpc) FindByCategory(ctx context.Context, request *pb.FindAllProductCategoryRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	category_name := request.GetCategoryName()
	min_price := int(request.GetMinPrice())
	max_price := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if min_price <= 0 {
		min_price = 0
	}

	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.FindAllProductByCategory{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		CategoryName: category_name,
		MinPrice:     &min_price,
		MaxPrice:     &max_price,
	}

	products, totalRecords, err := s.productService.FindByCategory(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		protoProducts[i] = &pb.ProductResponse{
			Id:           int32(product.ProductID),
			MerchantId:   int32(product.MerchantID),
			CategoryId:   int32(product.CategoryID),
			Name:         product.Name,
			Description:  utils.StringPtrToString(product.Description),
			Price:        int32(product.Price),
			CountInStock: int32(product.CountInStock),
			Brand:        utils.StringPtrToString(product.Brand),
			Weight:       utils.Int32PtrToInt32(product.Weight),
			Rating:       utils.Float64PtrToFloat32(product.Rating),
			SlugProduct:  utils.StringPtrToString(product.SlugProduct),
			ImageProduct: utils.StringPtrToString(product.ImageProduct),
			CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched category products",
		Data:       protoProducts,
		Pagination: paginationMeta,
	}, nil
}

func (s *productHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProduct := &pb.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  utils.StringPtrToString(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        utils.StringPtrToString(product.Brand),
		Weight:       utils.Int32PtrToInt32(product.Weight),
		Rating:       utils.Float64PtrToFloat32(product.Rating),
		SlugProduct:  utils.StringPtrToString(product.SlugProduct),
		ImageProduct: utils.StringPtrToString(product.ImageProduct),
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully fetched product",
		Data:    protoProduct,
	}, nil
}

func (s *productHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProducts := make([]*pb.ProductResponseDeleteAt, len(products))
	for i, product := range products {
		var deletedAt *wrapperspb.StringValue
		if product.DeletedAt.Valid {
			deletedAt = &wrapperspb.StringValue{Value: product.DeletedAt.Time.Format("2006-01-02")}
		}

		protoProducts[i] = &pb.ProductResponseDeleteAt{
			Id:           int32(product.ProductID),
			MerchantId:   int32(product.MerchantID),
			CategoryId:   int32(product.CategoryID),
			Name:         product.Name,
			Description:  utils.StringPtrToString(product.Description),
			Price:        int32(product.Price),
			CountInStock: int32(product.CountInStock),
			Brand:        utils.StringPtrToString(product.Brand),
			Weight:       utils.Int32PtrToInt32(product.Weight),
			Rating:       utils.Float64PtrToFloat32(product.Rating),
			SlugProduct:  utils.StringPtrToString(product.SlugProduct),
			ImageProduct: utils.StringPtrToString(product.ImageProduct),
			CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:    deletedAt,
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active products",
		Data:       protoProducts,
		Pagination: paginationMeta,
	}, nil
}

func (s *productHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProducts := make([]*pb.ProductResponseDeleteAt, len(products))
	for i, product := range products {
		var deletedAt *wrapperspb.StringValue
		if product.DeletedAt.Valid {
			deletedAt = &wrapperspb.StringValue{Value: product.DeletedAt.Time.Format("2006-01-02")}
		}

		protoProducts[i] = &pb.ProductResponseDeleteAt{
			Id:           int32(product.ProductID),
			MerchantId:   int32(product.MerchantID),
			CategoryId:   int32(product.CategoryID),
			Name:         product.Name,
			Description:  utils.StringPtrToString(product.Description),
			Price:        int32(product.Price),
			CountInStock: int32(product.CountInStock),
			Brand:        utils.StringPtrToString(product.Brand),
			Weight:       utils.Int32PtrToInt32(product.Weight),
			Rating:       utils.Float64PtrToFloat32(product.Rating),
			SlugProduct:  utils.StringPtrToString(product.SlugProduct),
			ImageProduct: utils.StringPtrToString(product.ImageProduct),
			CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:    deletedAt,
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed products",
		Data:       protoProducts,
		Pagination: paginationMeta,
	}, nil
}

func (s *productHandleGrpc) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.ApiResponseProduct, error) {
	rating := int(request.GetRating())
	slug := request.GetSlugProduct()

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := s.productService.CreateProduct(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProduct := &pb.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  utils.StringPtrToString(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        utils.StringPtrToString(product.Brand),
		Weight:       utils.Int32PtrToInt32(product.Weight),
		Rating:       utils.Float64PtrToFloat32(product.Rating),
		SlugProduct:  utils.StringPtrToString(product.SlugProduct),
		ImageProduct: utils.StringPtrToString(product.ImageProduct),
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully created product",
		Data:    protoProduct,
	}, nil
}

func (s *productHandleGrpc) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetProductId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	rating := int(request.GetRating())
	slug := request.GetSlugProduct()

	req := &requests.UpdateProductRequest{
		ProductID:    &id,
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productService.UpdateProduct(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoProduct := &pb.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  utils.StringPtrToString(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        utils.StringPtrToString(product.Brand),
		Weight:       utils.Int32PtrToInt32(product.Weight),
		Rating:       utils.Float64PtrToFloat32(product.Rating),
		SlugProduct:  utils.StringPtrToString(product.SlugProduct),
		ImageProduct: utils.StringPtrToString(product.ImageProduct),
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product",
		Data:    protoProduct,
	}, nil
}

func (s *productHandleGrpc) TrashedProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productService.TrashedProduct(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt *wrapperspb.StringValue
	if product.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: product.DeletedAt.Time.Format("2006-01-02")}
	}

	protoProduct := &pb.ProductResponseDeleteAt{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  utils.StringPtrToString(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        utils.StringPtrToString(product.Brand),
		Weight:       utils.Int32PtrToInt32(product.Weight),
		Rating:       utils.Float64PtrToFloat32(product.Rating),
		SlugProduct:  utils.StringPtrToString(product.SlugProduct),
		ImageProduct: utils.StringPtrToString(product.ImageProduct),
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully trashed product",
		Data:    protoProduct,
	}, nil
}

func (s *productHandleGrpc) RestoreProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productService.RestoreProduct(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt *wrapperspb.StringValue
	if product.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: product.DeletedAt.Time.Format("2006-01-02")}
	}

	protoProduct := &pb.ProductResponseDeleteAt{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  utils.StringPtrToString(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        utils.StringPtrToString(product.Brand),
		Weight:       utils.Int32PtrToInt32(product.Weight),
		Rating:       utils.Float64PtrToFloat32(product.Rating),
		SlugProduct:  utils.StringPtrToString(product.SlugProduct),
		ImageProduct: utils.StringPtrToString(product.ImageProduct),
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully restored product",
		Data:    protoProduct,
	}, nil
}

func (s *productHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	_, err := s.productService.DeleteProductPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductDelete{
		Status:  "success",
		Message: "Successfully deleted product permanently",
	}, nil
}

func (s *productHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	_, err := s.productService.RestoreAllProducts(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully restored all products",
	}, nil
}

func (s *productHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	_, err := s.productService.DeleteAllProductPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully deleted all products permanently",
	}, nil
}
