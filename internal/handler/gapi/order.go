package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors_custom"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleGrpc struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderService
	mapping      protomapper.OrderProtoMapper
}

func NewOrderHandleGrpc(
	orderService service.OrderService,
	mapping protomapper.OrderProtoMapper,
) *orderHandleGrpc {
	return &orderHandleGrpc{
		orderService: orderService,
		mapping:      mapping,
	}
}

func (s *orderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderService.FindAll(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationOrder(paginationMeta, "success", "Successfully fetched order", merchant)
	return so, nil
}

func (s *orderHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_request",
				Message: "Valid order ID is required",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	merchant, err := s.orderService.FindById(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully fetched order", merchant)

	return so, nil

}

func (s *orderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderService.FindByActive(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched active order", merchant)

	return so, nil
}

func (s *orderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderService.FindByTrashed(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched trashed order", users)

	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 || month >= 12 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_month",
				Message: "Month must be between 1 and 12",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderService.FindMonthlyTotalRevenue(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	methods, err := s.orderService.FindYearlyTotalRevenue(year)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenueById(ctx context.Context, req *pb.FindYearMonthTotalRevenueById) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetOrderId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 || month >= 12 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_month",
				Message: "Month must be between 1 and 12",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTotalRevenueOrder{
		OrderID: id,
		Month:   month,
		Year:    year,
	}

	methods, err := s.orderService.FindMonthlyTotalRevenueById(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenueById(ctx context.Context, req *pb.FindYearTotalRevenueById) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())
	id := int(req.GetOrderId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearTotalRevenueOrder{
		OrderID: id,
		Year:    year,
	}

	methods, err := s.orderService.FindYearlyTotalRevenueById(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if month <= 0 || month >= 12 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.orderService.FindMonthlyTotalRevenueByMerchant(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.orderService.FindYearlyTotalRevenueByMerchant(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.orderService.FindMonthlyOrder(year)
	if err != nil {
		return nil, status.Errorf(codes.Code(err.Code), "%v", &pb.ErrorResponse{
			Status:  err.Status,
			Message: err.Message,
		})
	}

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	res, err := s.orderService.FindYearlyOrder(year)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderService.FindMonthlyOrderByMerchant(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid year parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	if id <= 0 {
		return nil, status.Errorf(
			codes.Code(codes.InvalidArgument),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "invalid_input",
				Message: "Invalid id parameter",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderService.FindYearlyOrderByMerchant(&reqService)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error) {
	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.UserId),
		TotalPrice: int(request.GetTotalPrice()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Price:     int(item.GetPrice()),
		})
	}

	if request.Shipping != nil {
		req.ShippingAddress = requests.CreateShippingAddressRequest{
			OrderID:        int(request.Shipping.GetOrderId()),
			Alamat:         request.Shipping.GetAlamat(),
			Provinsi:       request.Shipping.GetProvinsi(),
			Kota:           request.Shipping.GetKota(),
			Courier:        request.Shipping.GetCourier(),
			ShippingMethod: request.Shipping.GetShippingMethod(),
			ShippingCost:   int(request.Shipping.GetShippingCost()),
			Negara:         request.Shipping.GetNegara(),
		}
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to create new order. Please check your input.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	order, err := s.orderService.CreateOrder(req)
	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully created order", order)
	return so, nil
}

func (s *orderHandleGrpc) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetOrderId())
	idShipping := int(request.GetShipping().GetShippingId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Order ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	req := &requests.UpdateOrderRequest{
		OrderID:    &id,
		TotalPrice: int(request.GetTotalPrice()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
			Price:       int(item.GetPrice()),
		})
	}

	if request.Shipping != nil {
		req.ShippingAddress = requests.UpdateShippingAddressRequest{
			ShippingID:     &idShipping,
			OrderID:        int(request.Shipping.GetOrderId()),
			Alamat:         request.Shipping.GetAlamat(),
			Provinsi:       request.Shipping.GetProvinsi(),
			Kota:           request.Shipping.GetKota(),
			Courier:        request.Shipping.GetCourier(),
			ShippingMethod: request.Shipping.GetShippingMethod(),
			ShippingCost:   int(request.Shipping.GetShippingCost()),
			Negara:         request.Shipping.GetNegara(),
		}
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Unable to process order update. Please review your data.",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	order, err := s.orderService.UpdateOrder(req)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully updated order", order)
	return so, nil
}

func (s *orderHandleGrpc) TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Order ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	merchant, err := s.orderService.TrashedOrder(id)

	if err != nil {
		return nil, status.Errorf(codes.Code(err.Code), "%v", &pb.ErrorResponse{
			Status:  err.Status,
			Message: err.Message,
			Code:    int32(err.Code),
		})
	}

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully trashed order", merchant)

	return so, nil
}

func (s *orderHandleGrpc) RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Order ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	merchant, err := s.orderService.RestoreOrder(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully restored order", merchant)

	return so, nil
}

func (s *orderHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  "validation_error",
				Message: "Order ID parameter cannot be empty and must be a positive number",
				Code:    int32(codes.InvalidArgument),
			}),
		)
	}

	_, err := s.orderService.DeleteOrderPermanent(id)

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrderDelete("success", "Successfully deleted order permanently")

	return so, nil
}

func (s *orderHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	_, err := s.orderService.RestoreAllOrder()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully restore all order")

	return so, nil
}

func (s *orderHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	_, err := s.orderService.DeleteAllOrderPermanent()

	if err != nil {
		return nil, status.Errorf(
			codes.Code(err.Code),
			errors_custom.GrpcErrorToJson(&pb.ErrorResponse{
				Status:  err.Status,
				Message: err.Message,
				Code:    int32(err.Code),
			}),
		)
	}

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully delete order permanen")

	return so, nil
}
