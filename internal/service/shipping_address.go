package service

import (
	"context"
	shippingaddress_cache "ecommerce/internal/cache/shipping_address"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type shippingAddressService struct {
	shippingRepository repository.ShippingAddressRepository
	logger             logger.LoggerInterface
	observability      observability.TraceLoggerObservability
	cache              shippingaddress_cache.ShippingAddressMencache
}

type ShippingAddressServiceDeps struct {
	ShippingRepository repository.ShippingAddressRepository
	Logger             logger.LoggerInterface
	Observability      observability.TraceLoggerObservability
	Cache              shippingaddress_cache.ShippingAddressMencache
}

func NewShippingAddressService(deps ShippingAddressServiceDeps) ShippingAddressService {
	return &shippingAddressService{
		shippingRepository: deps.ShippingRepository,
		logger:             deps.Logger,
		observability:      deps.Observability,
		cache:              deps.Cache,
	}
}

func (s *shippingAddressService) FindAllShippingAddress(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, *int, error) {
	const method = "FindAllShippingAddresses"

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

	if data, total, found := s.cache.GetShippingAddressAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingRepository.FindAllShippingAddress(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindAllShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressAllCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched all shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}

func (s *shippingAddressService) FindById(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, error) {
	const method = "FindShippingAddressById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedShippingAddressCache(ctx, shipping_id); found {
		logSuccess("Successfully retrieved shipping address by ID from cache",
			zap.Int("shipping_id", shipping_id))
		return data, nil
	}

	shippingAddress, err := s.shippingRepository.FindById(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetShippingByIDRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindShippingAddressByID,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.SetCachedShippingAddressCache(ctx, shippingAddress)

	logSuccess("Successfully fetched shipping address by ID",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressService) FindByOrder(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	const method = "FindShippingAddressByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedShippingAddressByOrderCache(ctx, order_id); found {
		logSuccess("Successfully retrieved shipping address by order ID from cache",
			zap.Int("order_id", order_id))
		return data, nil
	}

	shippingAddress, err := s.shippingRepository.FindByOrder(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetShippingAddressByOrderIDRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindShippingAddressByOrder,
			method,
			span,

			zap.Int("order_id", order_id),
		)
	}

	s.cache.SetCachedShippingAddressByOrderCache(ctx, shippingAddress)

	logSuccess("Successfully fetched shipping address by order ID",
		zap.Int("order_id", order_id))

	return shippingAddress, nil
}

func (s *shippingAddressService) FindByActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, *int, error) {
	const method = "FindActiveShippingAddresses"

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

	if data, total, found := s.cache.GetShippingAddressActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressActiveRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindActiveShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressActiveCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched active shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}

func (s *shippingAddressService) FindByTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, *int, error) {
	const method = "FindTrashedShippingAddresses"

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

	if data, total, found := s.cache.GetShippingAddressTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressTrashedRow](
			s.logger,
			shippingaddress_errors.ErrFailedFindTrashedShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressTrashedCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched trashed shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}

func (s *shippingAddressService) TrashShippingAddress(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "TrashShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingRepository.TrashShippingAddress(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ShippingAddress](
			s.logger,
			shippingaddress_errors.ErrFailedTrashShippingAddress,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.DeleteShippingAddressCache(ctx, shipping_id)

	logSuccess("Successfully trashed shipping address",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressService) RestoreShippingAddress(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "RestoreShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingRepository.RestoreShippingAddress(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ShippingAddress](
			s.logger,
			shippingaddress_errors.ErrFailedRestoreShippingAddress,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.DeleteShippingAddressCache(ctx, shipping_id)

	logSuccess("Successfully restored shipping address",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressService) DeleteShippingAddressPermanently(ctx context.Context, shipping_id int) (bool, error) {
	const method = "DeleteShippingAddressPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingRepository.DeleteShippingAddressPermanently(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.DeleteShippingAddressCache(ctx, shipping_id)

	logSuccess("Successfully permanently deleted shipping address",
		zap.Int("shipping_id", shipping_id))

	return success, nil
}

func (s *shippingAddressService) RestoreAllShippingAddress(ctx context.Context) (bool, error) {
	const method = "RestoreAllShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingRepository.RestoreAllShippingAddress(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shippingaddress_errors.ErrFailedRestoreAllShippingAddresses,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed shipping addresses")

	return success, nil
}

func (s *shippingAddressService) DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error) {
	const method = "DeleteAllPermanentShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingRepository.DeleteAllPermanentShippingAddress(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shippingaddress_errors.ErrFailedDeleteAllShippingAddressesPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed shipping addresses")

	return success, nil
}
