package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type orderServiceMapper struct {
	orderRepository     repository.OrderRepository
	orderItemRepository repository.OrderItemRepository
	productRepository   repository.ProductRepository
	userRepository      repository.UserRepository
	merchantRepository  repository.MerchantRepository
	shippingRepository  repository.ShippingAddressRepository
	logger              logger.LoggerInterface
	mapping             response_service.OrderResponseMapper
}

func NewOrderServiceMapper(
	orderRepository repository.OrderRepository,
	orderItemRepository repository.OrderItemRepository,
	userRepository repository.UserRepository,
	merchantRepository repository.MerchantRepository,
	productRepository repository.ProductRepository,
	shippingRepository repository.ShippingAddressRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderServiceMapper {
	return &orderServiceMapper{
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

func (s *orderServiceMapper) FindAll(page int, pageSize int, search string) ([]*response.OrderResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching orders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindAllOrders(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch user",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch users",
		}
	}

	orderResponse := s.mapping.ToOrdersResponse(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, int(totalRecords), nil
}

func (s *orderServiceMapper) FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order by ID", zap.Int("order_id", order_id))

	order, err := s.orderRepository.FindById(order_id)
	if err != nil {
		s.logger.Error("Failed to fetch cashier", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Cashier not found"}
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderServiceMapper) FindByActive(page int, pageSize int, search string) ([]*response.OrderResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching orders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch order",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch order",
		}
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, int(totalRecords), nil
}

func (s *orderServiceMapper) FindByTrashed(page int, pageSize int, search string) ([]*response.OrderResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching orders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch order",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch order",
		}
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, int(totalRecords), nil
}

func (s *orderServiceMapper) CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new order with items", zap.Int("merchantID", req.MerchantID), zap.Int("userID", req.UserID))

	_, err := s.merchantRepository.FindById(req.MerchantID)
	if err != nil {
		s.logger.Error("Merchant not found", zap.Int("merchantID", req.MerchantID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Merchant not found"}
	}

	_, err = s.userRepository.FindById(req.UserID)
	if err != nil {
		s.logger.Error("User not found", zap.Int("userID", req.UserID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "User not found"}
	}

	order, err := s.orderRepository.CreateOrder(&requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
		TotalPrice: 0,
	})
	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create order"}
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(item.ProductID)
		if err != nil {
			s.logger.Error("Product not found", zap.Int("productID", item.ProductID), zap.Error(err))
			return nil, &response.ErrorResponse{Status: "error", Message: "Product not found"}
		}

		if product.CountInStock < 1 {
			s.logger.Error("Product out of stock", zap.Int("productID", item.ProductID))
			return nil, &response.ErrorResponse{Status: "error", Message: "Product out of stock"}
		}

		_, err = s.orderItemRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		if err != nil {
			s.logger.Error("Failed to create order item", zap.Error(err))
			return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create order item"}
		}

		product.CountInStock -= item.Quantity
		_, err = s.productRepository.UpdateProductCountStock(product.ID, product.CountInStock)
		if err != nil {
			s.logger.Error("Failed to update product stock", zap.Error(err))
			return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update product stock"}
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
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create shipping address"}
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(order.ID)
	if err != nil {
		s.logger.Error("Failed to calculate total price", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to calculate total price"}
	}

	_, err = s.orderRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		s.logger.Error("Failed to update order total price", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update order total price"}
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderServiceMapper) UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating order with items", zap.Int("orderID", req.OrderID))

	existingOrder, err := s.orderRepository.FindById(req.OrderID)
	if err != nil {
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Order not found"}
	}

	_, err = s.userRepository.FindById(req.UserID)
	if err != nil {
		s.logger.Error("User not found", zap.Int("userID", req.UserID), zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "User not found"}
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(item.ProductID)
		if err != nil {
			s.logger.Error("Product not found", zap.Int("productID", item.ProductID), zap.Error(err))
			return nil, &response.ErrorResponse{Status: "error", Message: "Product not found"}
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemRepository.UpdateOrderItem(&requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})
			if err != nil {
				s.logger.Error("Failed to update order item", zap.Error(err))
				return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update order item"}
			}
		} else {
			_, err := s.orderItemRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
				OrderID:   req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
			if err != nil {
				s.logger.Error("Failed to create order item", zap.Error(err))
				return nil, &response.ErrorResponse{Status: "error", Message: "Failed to create order item"}
			}

			product.CountInStock -= item.Quantity
			_, err = s.productRepository.UpdateProductCountStock(product.ID, product.CountInStock)
			if err != nil {
				s.logger.Error("Failed to update product stock", zap.Error(err))
				return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update product stock"}
			}
		}
	}

	_, err = s.shippingRepository.UpdateShippingAddress(&requests.UpdateShippingAddressRequest{
		ShippingID:     req.ShippingAddress.ShippingID,
		OrderID:        req.OrderID,
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
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update shipping address"}
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(req.OrderID)
	if err != nil {
		s.logger.Error("Failed to calculate total price", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to calculate total price"}
	}

	_, err = s.orderRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		s.logger.Error("Failed to update order total price", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to update order total price"}
	}

	return s.mapping.ToOrderResponse(existingOrder), nil
}

func (s *orderServiceMapper) TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to fetch order items for trashing", zap.Error(err))

		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.TrashedOrderItem(item.ID)

		if err != nil {
			s.logger.Error("Failed to trash order item", zap.Int("order_item_id", item.ID), zap.Error(err))

			return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash order item"}
		}
	}

	order, err := s.orderRepository.TrashedOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to trash order", zap.Error(err))

		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash order"}
	}

	return s.mapping.ToOrderResponseDeleteAt(order), nil
}

func (s *orderServiceMapper) RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to fetch order items for restoring", zap.Error(err))

		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.RestoreOrderItem(item.ID)

		if err != nil {
			s.logger.Error("Failed to restore order item", zap.Int("order_item_id", item.ID), zap.Error(err))

			return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore order item"}
		}
	}

	order, err := s.orderRepository.RestoreOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to restore order", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore order"}
	}

	return s.mapping.ToOrderResponseDeleteAt(order), nil
}

func (s *orderServiceMapper) DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to fetch order items for permanent deletion", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.
			DeleteOrderItemPermanent(item.ID)

		if err != nil {
			s.logger.Error("Failed to permanently delete order item", zap.Int("order_item_id", item.ID), zap.Error(err))
			return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete order item"}
		}
	}

	success, err := s.orderRepository.DeleteOrderPermanent(order_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete order", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete order"}
	}

	return success, nil
}

func (s *orderServiceMapper) RestoreAllOrder() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed orders and related order items")

	successItems, err := s.orderItemRepository.RestoreAllOrderItem()

	if err != nil || !successItems {
		s.logger.Error("Failed to restore all order items", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all order items"}
	}

	success, err := s.orderRepository.RestoreAllOrder()

	if err != nil || !success {
		s.logger.Error("Failed to restore all orders", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all orders"}
	}

	return success, nil
}

func (s *orderServiceMapper) DeleteAllOrderPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all orders and related order items")

	successItems, err := s.orderItemRepository.DeleteAllOrderPermanent()

	if err != nil || !successItems {
		s.logger.Error("Failed to permanently delete all order items", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all order items"}
	}

	success, err := s.orderRepository.DeleteAllOrderPermanent()

	if err != nil || !success {
		s.logger.Error("Failed to permanently delete all orders", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all orders"}
	}

	return success, nil
}
