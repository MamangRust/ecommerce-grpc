package service

import (
	"context"
	order_cache "ecommerce/internal/cache/order"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/errors/order_errors"
	orderitem_errors "ecommerce/pkg/errors/order_item_errors"
	"ecommerce/pkg/errors/product_errors"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
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
	observability       observability.TraceLoggerObservability
	cache               order_cache.OrderMencache
}

type OrderServiceDeps struct {
	OrderRepository     repository.OrderRepository
	OrderItemRepository repository.OrderItemRepository
	ProductRepository   repository.ProductRepository
	UserRepository      repository.UserRepository
	MerchantRepository  repository.MerchantRepository
	ShippingRepository  repository.ShippingAddressRepository
	Logger              logger.LoggerInterface
	Observability       observability.TraceLoggerObservability
	Cache               order_cache.OrderMencache
}

func NewOrderService(deps OrderServiceDeps) *orderService {
	return &orderService{
		orderRepository:     deps.OrderRepository,
		orderItemRepository: deps.OrderItemRepository,
		productRepository:   deps.ProductRepository,
		userRepository:      deps.UserRepository,
		merchantRepository:  deps.MerchantRepository,
		shippingRepository:  deps.ShippingRepository,
		logger:              deps.Logger,
		observability:       deps.Observability,
		cache:               deps.Cache,
	}
}

func (s *orderService) FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, *int, error) {
	const method = "FindAllOrders"

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

	if data, total, found := s.cache.GetOrderAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindAllOrders(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersRow](
			s.logger,
			order_errors.ErrFailedFindAllOrders,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderAllCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched all orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderService) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, *int, error) {
	const method = "FindActiveOrders"

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

	if data, total, found := s.cache.GetOrderActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersActiveRow](
			s.logger,
			order_errors.ErrFailedFindOrdersByActive,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderActiveCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched active orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderService) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, *int, error) {
	const method = "FindTrashedOrders"

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

	if data, total, found := s.cache.GetOrderTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersTrashedRow](
			s.logger,
			order_errors.ErrFailedFindOrdersByTrashed,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderTrashedCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched trashed orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderService) FindById(ctx context.Context, orderID int) (*db.GetOrderByIDRow, error) {
	const method = "FindByIdOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedOrderCache(ctx, orderID); found {
		logSuccess("Successfully retrieved order by ID from cache", zap.Int("orderID", orderID))
		return data, nil
	}

	order, err := s.orderRepository.FindById(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetOrderByIDRow](
			s.logger,
			order_errors.ErrFailedFindOrderById,
			method,
			span,

			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedOrderCache(ctx, order)

	logSuccess("Successfully fetched order by ID", zap.Int("orderID", orderID))
	return order, nil
}

func (s *orderService) FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error) {
	const method = "FindMonthlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTotalRevenueCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total revenue from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.orderRepository.GetMonthlyTotalRevenue(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalRevenueRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyTotalRevenue,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthlyTotalRevenueCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total revenue from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *orderService) FindYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error) {
	const method = "FindYearlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTotalRevenueCache(ctx, year); found {
		logSuccess("Successfully retrieved yearly total revenue from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderRepository.GetYearlyTotalRevenue(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalRevenueRow](
			s.logger,
			order_errors.ErrFailedFindYearlyTotalRevenue,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTotalRevenueCache(ctx, year, res)

	logSuccess("Successfully fetched yearly total revenue from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *orderService) FindMonthlyTotalRevenueById(ctx context.Context, req *requests.MonthTotalRevenueOrder) ([]*db.GetMonthlyTotalRevenueByIdRow, error) {
	const method = "FindMonthlyTotalRevenueById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("orderID", req.OrderID))

	defer func() {
		end(status)
	}()

	res, err := s.orderRepository.GetMonthlyTotalRevenueById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalRevenueByIdRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyTotalRevenueById,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("orderID", req.OrderID),
		)
	}

	logSuccess("Successfully fetched monthly total revenue by ID",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("orderID", req.OrderID))

	return res, nil
}

func (s *orderService) FindYearlyTotalRevenueById(ctx context.Context, req *requests.YearTotalRevenueOrder) ([]*db.GetYearlyTotalRevenueByIdRow, error) {
	const method = "FindYearlyTotalRevenueById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("orderID", req.OrderID))

	defer func() {
		end(status)
	}()

	res, err := s.orderRepository.GetYearlyTotalRevenueById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalRevenueByIdRow](
			s.logger,
			order_errors.ErrFailedFindYearlyTotalRevenueById,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("orderID", req.OrderID),
		)
	}

	logSuccess("Successfully fetched yearly total revenue by ID",
		zap.Int("year", req.Year),
		zap.Int("orderID", req.OrderID))

	return res, nil
}

func (s *orderService) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	const method = "FindMonthlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total revenue by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderRepository.GetMonthlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalRevenueByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyTotalRevenueByMerchant,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total revenue by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderService) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	const method = "FindYearlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly total revenue by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderRepository.GetYearlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalRevenueByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindYearlyTotalRevenueByMerchant,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetYearlyTotalRevenueByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched yearly total revenue by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderService) FindMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error) {
	const method = "FindMonthlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyOrderCache(ctx, year); found {
		logSuccess("Successfully retrieved monthly orders from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderRepository.GetMonthlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyOrderRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyOrder,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyOrderCache(ctx, year, res)

	logSuccess("Successfully fetched monthly orders from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *orderService) FindYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error) {
	const method = "FindYearlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyOrderCache(ctx, year); found {
		logSuccess("Successfully retrieved yearly orders from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderRepository.GetYearlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyOrderRow](
			s.logger,
			order_errors.ErrFailedFindYearlyOrder,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyOrderCache(ctx, year, res)

	logSuccess("Successfully fetched yearly orders from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *orderService) FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	const method = "FindMonthlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly orders by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderRepository.GetMonthlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyOrderByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyOrderByMerchant,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetMonthlyOrderByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched monthly orders by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderService) FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	const method = "FindYearlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly orders by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderRepository.GetYearlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyOrderByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindYearlyOrderByMerchant,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetYearlyOrderByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched yearly orders by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*db.UpdateOrderRow, error) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("userID", req.UserID))

	defer func() {
		end(status)
	}()

	_, err := s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Int("merchantID", req.MerchantID),
		)
	}

	_, err = s.userRepository.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Int("user_id", req.UserID),
		)
	}

	order, err := s.orderRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			order_errors.ErrFailedCreateOrder,
			method,
			span,
		)
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](
				s.logger,
				product_errors.ErrFailedFindProductById,
				method,
				span,
				zap.Int("productID", item.ProductID),
			)
		}

		if int(product.CountInStock) < item.Quantity {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](
				s.logger,
				order_errors.ErrInsufficientProductStock,
				method,
				span,
				zap.Int("productID", item.ProductID),
				zap.Int("requested", item.Quantity),
				zap.Int("available", int(product.CountInStock)),
			)
		}

		_, err = s.orderItemRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   int(order.OrderID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     int(product.Price),
		})
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](
				s.logger,
				orderitem_errors.ErrFailedCreateOrderItem,
				method,
				span,
			)
		}

		product.CountInStock -= int32(item.Quantity)
		_, err = s.productRepository.UpdateProductCountStock(ctx, int(product.ProductID), int(product.CountInStock))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](
				s.logger,
				product_errors.ErrFailedCountStock,
				method,
				span,
			)
		}
	}

	orderId := int(order.OrderID)

	_, err = s.shippingRepository.CreateShippingAddress(ctx, &requests.CreateShippingAddressRequest{
		OrderID:        &orderId,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			shippingaddress_errors.ErrFailedCreateShippingAddress,
			method,
			span,
		)
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(ctx, int(order.OrderID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal,
			method,
			span,
		)
	}

	updatedOrder, err := s.orderRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(order.OrderID),
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			order_errors.ErrFailedUpdateOrder,
			method,
			span,
		)
	}

	logSuccess("Successfully created order with transfer",
		zap.Int("orderID", int(order.OrderID)),
		zap.Int("userID", req.UserID),
		zap.Int("merchantID", req.MerchantID),
		zap.Int("totalPrice", int(*totalPrice)))

	return updatedOrder, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*db.UpdateOrderRow, error) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", *req.OrderID),
		attribute.Int("userID", req.UserID))

	defer func() {
		end(status)
	}()

	existingOrder, err := s.orderRepository.FindById(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			order_errors.ErrFailedFindOrderById,
			method,
			span,
			zap.Int("orderID", *req.OrderID),
		)
	}

	_, err = s.userRepository.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Int("userID", req.UserID),
		)
	}

	for _, item := range req.Items {
		product, err := s.productRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](
				s.logger,
				product_errors.ErrFailedFindProductById,
				method,
				span,
				zap.Int("productID", item.ProductID),
			)
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       int(product.Price),
			})
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](
					s.logger,
					orderitem_errors.ErrFailedUpdateOrderItem,
					method,
					span,
					zap.Int("orderItemID", item.OrderItemID),
				)
			}
		} else {
			if product.CountInStock < int32(item.Quantity) {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](
					s.logger,
					order_errors.ErrInsufficientProductStock,
					method,
					span,
					zap.Int("productID", item.ProductID),
					zap.Int("requested", item.Quantity),
					zap.Int("available", int(product.CountInStock)),
				)
			}

			_, err := s.orderItemRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     int(product.Price),
			})
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](
					s.logger,
					orderitem_errors.ErrFailedCreateOrderItem,
					method,
					span,
					zap.Int("orderID", *req.OrderID),
					zap.Int("productID", item.ProductID),
				)
			}

			product.CountInStock -= int32(item.Quantity)
			_, err = s.productRepository.UpdateProductCountStock(ctx, int(product.ProductID), int(product.CountInStock))
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](
					s.logger,
					product_errors.ErrFailedCountStock,
					method,
					span,
					zap.Int("productID", int(product.ProductID)),
				)
			}
		}
	}

	orderId := int(existingOrder.OrderID)

	_, err = s.shippingRepository.UpdateShippingAddress(ctx, &requests.UpdateShippingAddressRequest{
		ShippingID:     req.ShippingAddress.ShippingID,
		OrderID:        &orderId,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			shippingaddress_errors.ErrFailedUpdateShippingAddress,
			method,
			span,
			zap.Int("shippingID", *req.ShippingAddress.ShippingID),
		)
	}

	totalPrice, err := s.orderItemRepository.CalculateTotalPrice(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal,
			method,
			span,
			zap.Int("orderID", *req.OrderID),
		)
	}

	updated, err := s.orderRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](
			s.logger,
			order_errors.ErrFailedUpdateOrder,
			method,
			span,
			zap.Int("orderID", *req.OrderID),
		)
	}

	logSuccess("Successfully updated order",
		zap.Int("orderID", *req.OrderID),
		zap.Int("userID", req.UserID),
		zap.Int("totalPrice", int(*totalPrice)))

	s.cache.DeleteOrderCache(ctx, *req.OrderID)

	return updated, nil
}

func (s *orderService) TrashedOrder(ctx context.Context, order_id int) (*db.Order, error) {
	const method = "TrashedOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	_, err := s.orderRepository.FindById(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedFindOrderById,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedNotDeleteAtOrder,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.TrashedOrderItem(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedTrashedOrderItem,
				method,
				span,
				zap.Int("order_item_id", int(item.OrderItemID)),
			)
		}
	}

	shipping, err := s.shippingRepository.FindByOrder(ctx, order_id)
	if err == nil && shipping != nil {
		_, err := s.shippingRepository.TrashShippingAddress(ctx, int(shipping.ShippingAddressID))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.Order](
				s.logger,
				shippingaddress_errors.ErrTrashShippingAddress,
				method,
				span,
				zap.Int("order_id", order_id),
			)
		}
	}

	trashedOrder, err := s.orderRepository.TrashedOrder(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedCreateOrder,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	logSuccess("Order moved to trash successfully",
		zap.Int("order_id", order_id),
		zap.String("deleted_at", trashedOrder.DeletedAt.Time.String()))

	s.cache.DeleteOrderCache(ctx, order_id)

	return trashedOrder, nil
}

func (s *orderService) RestoreOrder(ctx context.Context, order_id int) (*db.Order, error) {
	const method = "RestoreOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemRepository.FindOrderItemByOrderTrashed(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.RestoreOrderItem(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedRestoreOrderItem,
				method,
				span,
				zap.Int("order_item_id", int(item.OrderItemID)),
			)
		}
	}

	order, err := s.orderRepository.RestoreOrder(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedRestoreOrder,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	logSuccess("Order restored successfully", zap.Int("order_id", order_id))

	s.cache.DeleteOrderCache(ctx, order_id)

	return order, nil
}

func (s *orderService) DeleteOrderPermanent(ctx context.Context, order_id int) (bool, error) {
	const method = "DeleteOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemRepository.FindOrderItemByOrderTrashed(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	for _, item := range orderItems {
		_, err := s.orderItemRepository.DeleteOrderItemPermanent(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[bool](
				s.logger,
				orderitem_errors.ErrFailedDeleteOrderItem,
				method,
				span,
				zap.Int("order_item_id", int(item.OrderItemID)),
			)
		}
	}
	shipping, err := s.shippingRepository.FindByOrder(ctx, order_id)
	if err != nil {
		// Try to find trashed shipping address if active one not found
		trashedShipping, errTrashed := s.shippingRepository.FindTrashedByOrder(ctx, order_id)
		if errTrashed == nil && trashedShipping != nil {
			_, err := s.shippingRepository.DeleteShippingAddressPermanently(ctx, int(trashedShipping.ShippingAddressID))
			if err != nil {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent,
					method,
					span,
					zap.Int("order_id", order_id),
				)
			}
		}
	} else if shipping != nil {
		_, err := s.shippingRepository.DeleteShippingAddressPermanently(ctx, int(shipping.ShippingAddressID))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[bool](
				s.logger,
				shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent,
				method,
				span,
				zap.Int("order_id", order_id),
			)
		}
	}

	success, err := s.orderRepository.DeleteOrderPermanent(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedDeleteOrderPermanent,
			method,
			span,
			zap.Int("order_id", order_id),
		)
	}

	logSuccess("Order permanently deleted successfully", zap.Int("order_id", order_id))

	s.cache.DeleteOrderCache(ctx, order_id)

	return success, nil
}

func (s *orderService) RestoreAllOrder(ctx context.Context) (bool, error) {
	const method = "RestoreAllOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemRepository.RestoreAllOrderItem(ctx)
	if err != nil || !successItems {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedRestoreAllOrderItem,
			method,
			span,
		)
	}

	success, err := s.orderRepository.RestoreAllOrder(ctx)
	if err != nil || !success {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedRestoreAllOrder,
			method,
			span,
		)
	}

	logSuccess("All orders restored successfully")

	return success, nil
}

func (s *orderService) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemRepository.DeleteAllOrderPermanent(ctx)
	if err != nil || !successItems {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteAllOrderItem,
			method,
			span,
		)
	}

	success, err := s.orderRepository.DeleteAllOrderPermanent(ctx)
	if err != nil || !success {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedDeleteAllOrderPermanent,
			method,
			span,
		)
	}

	logSuccess("All orders permanently deleted successfully")

	return success, nil
}
