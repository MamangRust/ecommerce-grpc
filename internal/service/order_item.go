package service

import (
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type orderItemService struct {
	orderItemRepository repository.OrderItemRepository
	logger              logger.LoggerInterface
	mapping             response_service.OrderItemResponseMapper
}

func NewOrderItemService(
	orderItemRepository repository.OrderItemRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderItemResponseMapper,
) *orderItemService {
	return &orderItemService{
		orderItemRepository: orderItemRepository,
		logger:              logger,
		mapping:             mapping,
	}
}

func (s *orderItemService) FindAllOrderItems(search string, page, pageSize int) ([]*response.OrderItemResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching all order items", zap.String("search", search), zap.Int("page", page), zap.Int("pageSize", pageSize))

	orderItems, total, err := s.orderItemRepository.FindAllOrderItems(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch order items", zap.Error(err))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items"}
	}

	return s.mapping.ToOrderItemsResponse(orderItems), total, nil
}

func (s *orderItemService) FindByActive(search string, page, pageSize int) ([]*response.OrderItemResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active order items", zap.String("search", search), zap.Int("page", page), zap.Int("pageSize", pageSize))

	orderItems, total, err := s.orderItemRepository.FindByActive(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch active order items", zap.Error(err))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active order items"}
	}

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), total, nil
}

func (s *orderItemService) FindByTrashed(search string, page, pageSize int) ([]*response.OrderItemResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed order items", zap.String("search", search), zap.Int("page", page), zap.Int("pageSize", pageSize))

	orderItems, total, err := s.orderItemRepository.FindByTrashed(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch trashed order items", zap.Error(err))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed order items"}
	}

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), total, nil
}

func (s *orderItemService) FindOrderItemByOrder(orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order items by order", zap.Int("order_id", orderID))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(orderID)
	if err != nil {
		s.logger.Error("Failed to fetch order items by order", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to fetch order items by order"}
	}

	return s.mapping.ToOrderItemsResponse(orderItems), nil
}
