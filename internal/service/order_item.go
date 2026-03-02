package service

import (
	"context"
	orderitem_cache "ecommerce/internal/cache/order_item"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemService struct {
	orderItemRepository repository.OrderItemRepository
	logger              logger.LoggerInterface
	observability       observability.TraceLoggerObservability
	cache               orderitem_cache.OrderItemMencache
}

type OrderItemServiceDeps struct {
	OrderItemRepository repository.OrderItemRepository
	Logger              logger.LoggerInterface
	Observability       observability.TraceLoggerObservability
	Cache               orderitem_cache.OrderItemMencache
}

func NewOrderItemService(deps OrderItemServiceDeps) *orderItemService {
	return &orderItemService{
		orderItemRepository: deps.OrderItemRepository,
		logger:              deps.Logger,
		observability:       deps.Observability,
		cache:               deps.Cache,
	}
}

func (s *orderItemService) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, error) {
	const method = "FindAllOrderItems"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemsAll(ctx, req); found {
		logSuccess("Successfully retrieved all order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, err := s.orderItemRepository.FindAllOrderItems(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsRow](
			s.logger,
			orderitem_errors.ErrFailedFindAllOrderItems,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemsAll(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched all order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemService) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, error) {
	const method = "FindByActiveOrderItems"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemActive(ctx, req); found {
		logSuccess("Successfully retrieved active order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	// Fixed: Added ctx to the repository call
	orderItems, err := s.orderItemRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsActiveRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByActive,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemActive(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched active order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemService) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, error) {
	const method = "FindByTrashedOrderItems"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, err := s.orderItemRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsTrashedRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByTrashed,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemTrashed(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched trashed order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemService) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.GetOrderItemsByOrderRow, error) {
	const method = "FindOrderItemByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedOrderItems(ctx, orderID); found {
		totalCount := len(data)
		logSuccess("Successfully retrieved order items by order ID from cache",
			zap.Int("order_id", orderID),
			zap.Int("totalRecords", totalCount))
		return data, nil
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetOrderItemsByOrderRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,

			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedOrderItems(ctx, orderItems)

	logSuccess("Successfully fetched order items by order ID",
		zap.Int("order_id", orderID))

	return orderItems, nil
}
