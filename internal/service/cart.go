package service

import (
	"context"
	cart_cache "ecommerce/internal/cache/cart"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/cart_errors"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cartService struct {
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	cartRepository    repository.CartRepository
	logger            logger.LoggerInterface
	cache             cart_cache.CartMencache
	observability     observability.TraceLoggerObservability
}

type CartServiceDeps struct {
	ProductRepository repository.ProductRepository
	UserRepository    repository.UserRepository
	CartRepository    repository.CartRepository
	Logger            logger.LoggerInterface
	Cache             cart_cache.CartMencache
	Observability     observability.TraceLoggerObservability
}

func NewCartService(deps CartServiceDeps) *cartService {
	return &cartService{
		productRepository: deps.ProductRepository,
		userRepository:    deps.UserRepository,
		cartRepository:    deps.CartRepository,
		logger:            deps.Logger,
		cache:             deps.Cache,
		observability:     deps.Observability,
	}
}

func (s *cartService) FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, *int, error) {
	const method = "FindAllCarts"

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

	if data, total, found := s.cache.GetCachedCartsCache(ctx, req); found {
		logSuccess("Successfully retrieved all cart records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	carts, err := s.cartRepository.FindCarts(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCartsRow](
			s.logger,
			cart_errors.ErrFailedFindAllCarts,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(carts) > 0 {
		totalCount = int(carts[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCartsCache(ctx, req, carts, &totalCount)

	logSuccess("Successfully fetched all carts",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return carts, &totalCount, nil
}

func (s *cartService) CreateCart(ctx context.Context, req *requests.CreateCartRequest) (*db.Cart, error) {
	const method = "CreateCart"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.FindById(ctx, req.ProductID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			product_errors.ErrFailedFindProductById,
			method,
			span,

			zap.Int("product_id", req.ProductID),
		)
	}

	_, err = s.userRepository.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", req.UserID),
		)
	}

	var imageProduct string
	if product.ImageProduct != nil {
		imageProduct = *product.ImageProduct
	}

	var weight int
	if product.Weight != nil {
		weight = int(*product.Weight)
	}

	cartRecord := &requests.CartCreateRecord{
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		Name:         product.Name,
		Price:        int(product.Price),
		ImageProduct: imageProduct,
		Quantity:     req.Quantity,
		Weight:       weight,
	}

	res, err := s.cartRepository.CreateCart(ctx, cartRecord)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			cart_errors.ErrFailedCreateCart,
			method,
			span,

			zap.Any("request", req),
		)
	}

	logSuccess("Successfully created cart", zap.Int("cartID", int(res.CartID)))
	return res, nil
}

func (s *cartService) DeletePermanent(ctx context.Context, cartID int) (bool, error) {
	const method = "DeletePermanentCart"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cartID", cartID))

	defer func() {
		end(status)
	}()

	success, err := s.cartRepository.DeletePermanent(ctx, cartID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			cart_errors.ErrFailedDeleteCart,
			method,
			span,

			zap.Int("cart_id", cartID),
		)
	}

	logSuccess("Successfully deleted cart permanently", zap.Int("cartID", cartID))
	return success, nil
}

func (s *cartService) DeleteAllPermanently(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	const method = "DeleteAllPermanentlyCarts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.cartRepository.DeleteAllPermanently(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			cart_errors.ErrFailedDeleteAllCarts,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all carts permanently")
	return success, nil
}
