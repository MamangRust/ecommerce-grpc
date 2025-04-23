package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"net/http"

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

func (s *orderItemService) FindAllOrderItems(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order items",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.orderItemRepository.FindAllOrderItems(req)

	if err != nil {
		s.logger.Error("Failed to fetch all order items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve all order items list",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched order-item",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponse(orderItems), totalRecords, nil
}

func (s *orderItemService) FindByActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order items active",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.orderItemRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active order-items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active order-items",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched order-items",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), totalRecords, nil
}

func (s *orderItemService) FindByTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order items trashed",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.orderItemRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed order-items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed order-items",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched order-items",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), totalRecords, nil
}

func (s *orderItemService) FindOrderItemByOrder(orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order items by order", zap.Int("order_id", orderID))

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(orderID)

	if err != nil {
		s.logger.Error("Failed to retrieve order items",
			zap.Error(err),
			zap.Int("order_id", orderID))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Unable to retrieve order items",
			Code:    http.StatusInternalServerError,
		}
	}
	return s.mapping.ToOrderItemsResponse(orderItems), nil
}
