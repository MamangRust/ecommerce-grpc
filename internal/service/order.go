package service

import (
	"database/sql"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type orderService struct {
	orderRepository     repository.OrderRepository
	orderItemRepository repository.OrderItemRepository
	productRepository   repository.ProductRepository
	userRepository      repository.UserRepository
	merchantRepository  repository.MerchantRepository
	shippingRepository  repository.ShippingAddressRepository
	logger              logger.LoggerInterface
	mapping             response_service.OrderResponseMapper
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	orderItemRepository repository.OrderItemRepository,
	userRepository repository.UserRepository,
	merchantRepository repository.MerchantRepository,
	productRepository repository.ProductRepository,
	shippingRepository repository.ShippingAddressRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderService {
	return &orderService{
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
		productRepository:   productRepository,
		userRepository:      userRepository,
		merchantRepository:  merchantRepository,
		shippingRepository:  shippingRepository,
		logger:              logger,
		mapping:             mapping,
	}
}

func (s *orderService) FindAll(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all orders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindAllOrders(req)

	if err != nil {
		s.logger.Error("Failed to retrieve order list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order list",
			Code:    http.StatusInternalServerError,
		}
	}

	orderResponse := s.mapping.ToOrdersResponse(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderService) FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order by ID", zap.Int("order_id", order_id))

	order, err := s.orderRepository.FindById(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve order details",
			zap.Error(err),
			zap.Int("order_id", order_id))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Order with ID %d not found", order_id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order details",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderService) FindByActive(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active orders",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active orders",
			Code:    http.StatusInternalServerError,
		}
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderService) FindByTrashed(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed orders from database",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed orders",
			Code:    http.StatusInternalServerError,
		}
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderService) FindMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.orderRepository.GetMonthlyTotalRevenue(req)

	if err != nil {
		s.logger.Error("failed to get monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderMonthlyTotalRevenues(res), nil
}

func (s *orderService) FindYearlyTotalRevenue(year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	res, err := s.orderRepository.GetYearlyTotalRevenue(year)

	if err != nil {
		s.logger.Error("failed to get yearly total revenue",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderYearlyTotalRevenues(res), nil
}

func (s *orderService) FindMonthlyTotalRevenueById(req *requests.MonthTotalRevenueOrder) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.orderRepository.GetMonthlyTotalRevenueById(req)

	if err != nil {
		s.logger.Error("failed to get monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderMonthlyTotalRevenues(res), nil
}

func (s *orderService) FindYearlyTotalRevenueById(req *requests.YearTotalRevenueOrder) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	year := req.Year
	orderId := req.OrderID

	res, err := s.orderRepository.GetYearlyTotalRevenueById(req)

	if err != nil {
		s.logger.Error("failed to get yearly total revenue",
			zap.Int("year", year),
			zap.Int("order_id", orderId),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderYearlyTotalRevenues(res), nil
}

func (s *orderService) FindMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	res, err := s.orderRepository.GetMonthlyTotalRevenueByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderMonthlyTotalRevenues(res), nil
}

func (s *orderService) FindYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	year := req.Year
	merchantId := req.MerchantID

	res, err := s.orderRepository.GetYearlyTotalRevenueByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get yearly total revenue",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly total revenue data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderYearlyTotalRevenues(res), nil
}

func (s *orderService) FindMonthlyOrder(year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	res, err := s.orderRepository.GetMonthlyOrder(year)

	if err != nil {
		s.logger.Error("failed to get monthly orders",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve monthly orders data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderMonthlyPrices(res), nil
}

func (s *orderService) FindYearlyOrder(year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	res, err := s.orderRepository.GetYearlyOrder(year)

	if err != nil {
		s.logger.Error("failed to get yearly orders",
			zap.Int("year", year),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to retrieve yearly orders data",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderYearlyPrices(res), nil
}

func (s *orderService) FindMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	year := req.Year
	merchant_id := req.MerchantID

	res, err := s.orderRepository.GetMonthlyOrderByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get monthly orders by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: fmt.Sprintf("Failed to retrieve monthly orders for merchant %d", merchant_id),
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderMonthlyPrices(res), nil
}

func (s *orderService) FindYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	year := req.Year
	merchant_id := req.MerchantID

	res, err := s.orderRepository.GetYearlyOrderByMerchant(req)

	if err != nil {
		s.logger.Error("failed to get yearly orders by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "server_error",
			Message: fmt.Sprintf("Failed to retrieve yearly orders for merchant %d", merchant_id),
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderYearlyPrices(res), nil
}

func (s *orderService) CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new order with items", zap.Int("merchantID", req.MerchantID), zap.Int("userID", req.UserID))

	_, err := s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		s.logger.Error("Merchant not found for order creation",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "not_found",
			Message: fmt.Sprintf("Merchant with ID %d not found", req.MerchantID),
			Code:    http.StatusNotFound,
		}
	}

	_, err = s.userRepository.FindById(req.UserID)

	if err != nil {
		s.logger.Error("User not found for order creation",
			zap.Int("user_id", req.UserID),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "not_found",
			Message: fmt.Sprintf("User with ID %d not found", req.UserID),
			Code:    http.StatusNotFound,
		}
	}

	order, err := s.orderRepository.CreateOrder(&requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})

	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create order record",
			Code:    http.StatusInternalServerError,
		}
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(item.ProductID)

		if err != nil {
			s.logger.Error("Product not found for order item",
				zap.Int("productID", item.ProductID),
				zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Product with ID %d not found", item.ProductID),
				Code:    http.StatusNotFound,
			}
		}

		if product.CountInStock < item.Quantity {
			s.logger.Error("Insufficient product stock",
				zap.Int("productID", item.ProductID),
				zap.Int("requested", item.Quantity),
				zap.Int("available", product.CountInStock))

			return nil, &response.ErrorResponse{
				Status: "invalid_request",
				Message: fmt.Sprintf("Insufficient stock for product %d (requested %d, available %d)",
					item.ProductID, item.Quantity, product.CountInStock),
				Code: http.StatusBadRequest,
			}
		}

		_, err = s.orderItemRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		if err != nil {
			s.logger.Error("Failed to create order item",
				zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to create order item",
				Code:    http.StatusInternalServerError,
			}
		}

		product.CountInStock -= item.Quantity
		_, err = s.productRepository.UpdateProductCountStock(product.ID, product.CountInStock)

		if err != nil {
			s.logger.Error("Failed to update product stock",
				zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to update product stock",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	_, err = s.shippingRepository.CreateShippingAddress(&requests.CreateShippingAddressRequest{
		OrderID:        order.ID,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		s.logger.Error("Failed to create shipping address", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create shipping address",
			Code:    http.StatusInternalServerError,
		}
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(order.ID)

	if err != nil {
		s.logger.Error("Failed to calculate total price", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to calculate order total",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.orderRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		s.logger.Error("Failed to update order total price",
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update order total price",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderService) UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating order with items", zap.Int("orderID", *req.OrderID))

	existingOrder, err := s.orderRepository.FindById(*req.OrderID)

	if err != nil {
		s.logger.Error("Order not found for update",
			zap.Int("orderID", *req.OrderID),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "not_found",
			Message: fmt.Sprintf("Order with ID %d not found", req.OrderID),
			Code:    http.StatusNotFound,
		}
	}

	_, err = s.userRepository.FindById(req.UserID)
	if err != nil {
		s.logger.Error("Order not found for order creation",
			zap.Int("orderID", *req.OrderID),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "not_found",
			Message: fmt.Sprintf("Order with ID %d not found", *req.OrderID),
			Code:    http.StatusNotFound,
		}
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(item.ProductID)

		if err != nil {
			s.logger.Error("Product not found for order item",
				zap.Int("productID", item.ProductID),
				zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Product with ID %d not found", item.ProductID),
				Code:    http.StatusNotFound,
			}
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemRepository.UpdateOrderItem(&requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})

			if err != nil {
				s.logger.Error("Failed to update order item",
					zap.Error(err))
				return nil, &response.ErrorResponse{
					Status:  "error",
					Message: "Failed to update order item",
					Code:    http.StatusInternalServerError,
				}
			}
		} else {
			if product.CountInStock < item.Quantity {
				s.logger.Error("Insufficient product stock for new order item",
					zap.Int("productID", item.ProductID),
					zap.Int("requested", item.Quantity),
					zap.Int("available", product.CountInStock))
				return nil, &response.ErrorResponse{
					Status: "invalid_request",
					Message: fmt.Sprintf("Insufficient stock for product %d (requested %d, available %d)",
						item.ProductID, item.Quantity, product.CountInStock),
					Code: http.StatusBadRequest,
				}
			}

			_, err := s.orderItemRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})

			if err != nil {
				s.logger.Error("Failed to add new order item",
					zap.Error(err))
				return nil, &response.ErrorResponse{
					Status:  "error",
					Message: "Failed to add new order item",
					Code:    http.StatusInternalServerError,
				}
			}

			product.CountInStock -= item.Quantity
			_, err = s.productRepository.UpdateProductCountStock(product.ID, product.CountInStock)

			if err != nil {
				s.logger.Error("Failed to update product stock",
					zap.Error(err))
				return nil, &response.ErrorResponse{
					Status:  "error",
					Message: "Failed to update product stock",
					Code:    http.StatusInternalServerError,
				}
			}
		}
	}

	_, err = s.shippingRepository.UpdateShippingAddress(&requests.UpdateShippingAddressRequest{
		ShippingID:     req.ShippingAddress.ShippingID,
		OrderID:        *req.OrderID,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		s.logger.Error("Failed to update shipping address", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update shipping address",
			Code:    http.StatusInternalServerError,
		}
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(*req.OrderID)

	if err != nil {
		s.logger.Error("Failed to calculate updated order total",
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to calculate order total",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.orderRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		s.logger.Error("Failed to update order total price",
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update order total price",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderResponse(existingOrder), nil
}

func (s *orderService) TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Moving order to trash",
		zap.Int("order_id", order_id))

	order, err := s.orderRepository.FindById(order_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Order with ID %d not found", order_id),
				Code:    http.StatusNotFound,
			}
		}
		s.logger.Error("Failed to fetch order",
			zap.Int("order_id", order_id),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to verify order existence",
			Code:    http.StatusInternalServerError,
		}
	}

	if order.DeletedAt != nil {
		return nil, &response.ErrorResponse{
			Status:  "already_trashed",
			Message: fmt.Sprintf("Order with ID %d is already trashed", order_id),
			Code:    http.StatusBadRequest,
		}
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve order items for trashing",
			zap.Int("order_id", order_id),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order items",
			Code:    http.StatusInternalServerError,
		}
	}

	for _, item := range orderItems {
		if item.DeletedAt != nil {
			s.logger.Debug("Order item already trashed, skipping",
				zap.Int("order_item_id", item.ID))
			continue
		}

		trashedItem, err := s.orderItemRepository.TrashedOrderItem(item.ID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.logger.Debug("Order item not found - may have been deleted",
					zap.Int("order_item_id", item.ID))
				continue
			}

			s.logger.Error("Failed to move order item to trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: fmt.Sprintf("Failed to move order item %d to trash", item.ID),
				Code:    http.StatusInternalServerError,
			}
		}

		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", trashedItem.ID),
			zap.String("deleted_at", *trashedItem.DeletedAt))
	}

	trashedOrder, err := s.orderRepository.TrashedOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to move order to trash",
			zap.Int("order_id", order_id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Order with ID %d not found", order_id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move order to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Order moved to trash successfully",
		zap.Int("order_id", order_id),
		zap.String("deleted_at", *trashedOrder.DeletedAt))

	return s.mapping.ToOrderResponseDeleteAt(trashedOrder), nil
}

func (s *orderService) RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve order items for restoration",
			zap.Int("order_id", order_id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order items",
			Code:    http.StatusInternalServerError,
		}
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.RestoreOrderItem(item.ID)

		if err != nil {
			s.logger.Error("Failed to restore order item from trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err))

			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: fmt.Sprintf("Failed to restore order item %d", item.ID),
				Code:    http.StatusInternalServerError,
			}
		}
	}

	order, err := s.orderRepository.RestoreOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to restore order from trash",
			zap.Int("order_id", order_id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Order with ID %d not found in trash", order_id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore order from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToOrderResponseDeleteAt(order), nil
}

func (s *orderService) DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve order items for permanent deletion",
			zap.Int("order_id", order_id),
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve order items",
			Code:    http.StatusInternalServerError,
		}
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.
			DeleteOrderItemPermanent(item.ID)

		if err != nil {
			s.logger.Error("Failed to permanently delete order item",
				zap.Int("order_item_id", item.ID),
				zap.Error(err))
			return false, &response.ErrorResponse{
				Status:  "error",
				Message: fmt.Sprintf("Failed to permanently delete order item %d", item.ID),
				Code:    http.StatusInternalServerError,
			}
		}
	}

	success, err := s.orderRepository.DeleteOrderPermanent(order_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete order",
			zap.Int("order_id", order_id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return false, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Order with ID %d not found", order_id),
				Code:    http.StatusNotFound,
			}
		}

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete order",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *orderService) RestoreAllOrder() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed orders and related order items")

	successItems, err := s.orderItemRepository.RestoreAllOrderItem()

	if err != nil || !successItems {
		s.logger.Error("Failed to restore all order items",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all order items",
			Code:    http.StatusInternalServerError,
		}
	}

	success, err := s.orderRepository.RestoreAllOrder()
	if err != nil || !success {
		s.logger.Error("Failed to restore all orders",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all orders",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *orderService) DeleteAllOrderPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all orders and related order items")

	successItems, err := s.orderItemRepository.DeleteAllOrderPermanent()

	if err != nil || !successItems {
		s.logger.Error("Failed to permanently delete all order items",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all order items",
			Code:    http.StatusInternalServerError,
		}
	}

	success, err := s.orderRepository.DeleteAllOrderPermanent()

	if err != nil || !success {
		s.logger.Error("Failed to permanently delete all orders",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all orders",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}
