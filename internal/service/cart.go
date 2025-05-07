package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/cart_errors"
	"ecommerce/pkg/errors/product_errors"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type cartService struct {
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	cartRepository    repository.CartRepository
	logger            logger.LoggerInterface
	mapping           response_service.CartResponseMapper
}

func NewCartService(
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
	cartRepository repository.CartRepository,
	logger logger.LoggerInterface,
	mapping response_service.CartResponseMapper,
) *cartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
		userRepository:    userRepository,
		logger:            logger,
		mapping:           mapping,
	}
}

func (s *cartService) FindAll(req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching cart",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	cart, totalRecords, err := s.cartRepository.FindCarts(req)

	if err != nil {
		s.logger.Error("Failed to fetch cart",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, cart_errors.ErrFailedFindAllCarts
	}

	cartRes := s.mapping.ToCartsResponse(cart)

	s.logger.Debug("Successfully fetched cart",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return cartRes, totalRecords, nil
}

func (s *cartService) CreateCart(req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse) {
	product, err := s.productRepository.FindById(req.ProductID)

	if err != nil {
		s.logger.Error("Failed to retrieve product details",
			zap.Error(err),
			zap.Int("product_id", req.UserID))

		return nil, product_errors.ErrFailedFindProductById
	}

	_, err = s.userRepository.FindById(req.UserID)

	if err != nil {
		s.logger.Error("Failed to retrieve user details",
			zap.Error(err),
			zap.Int("user_id", req.UserID))

		return nil, user_errors.ErrUserNotFoundRes
	}

	cartRecord := &requests.CartCreateRecord{
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		Name:         product.Name,
		Price:        product.Price,
		ImageProduct: product.ImageProduct,
		Quantity:     req.Quantity,
		Weight:       product.Weight,
	}

	res, err := s.cartRepository.CreateCart(cartRecord)

	if err != nil {
		s.logger.Error("Failed to create new cart",
			zap.Error(err),
			zap.Any("request", req))

		return nil, cart_errors.ErrFailedCreateCart
	}

	so := s.mapping.ToCartResponse(res)

	return so, nil
}

func (s *cartService) DeletePermanent(cart_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting cart", zap.Int("cart_id", cart_id))

	success, err := s.cartRepository.DeletePermanent(cart_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete cart",
			zap.Error(err),
			zap.Int("cart_id", cart_id))

		return false, cart_errors.ErrFailedDeleteCart
	}

	return success, nil
}

func (s *cartService) DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting cart all", zap.Any("cart_id", req.CartIds))

	success, err := s.cartRepository.DeleteAllPermanently(req)

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed cart",
			zap.Error(err))

		return false, cart_errors.ErrFailedDeleteAllCarts
	}

	return success, nil
}
