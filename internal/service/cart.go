package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"fmt"

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

func (s *cartService) FindAll(user_id int, page int, pageSize int, search string) ([]*response.CartResponse, int, *response.ErrorResponse) {
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

	cart, totalRecords, err := s.cartRepository.FindCarts(user_id, search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch cart",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch cart",
		}
	}

	cartRes := s.mapping.ToCartsResponse(cart)

	s.logger.Debug("Successfully fetched cart",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return cartRes, int(totalRecords), nil
}

func (s *cartService) CreateCart(req *requests.CreateCartRequest) (*response.CartResponse, error) {
	product, err := s.productRepository.FindById(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	_, err = s.userRepository.FindById(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
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
		return nil, fmt.Errorf("failed to create cart: %v", err)
	}

	so := s.mapping.ToCartResponse(res)

	return so, nil
}

func (s *cartService) DeletePermanent(cart_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting cart", zap.Int("cart_id", cart_id))

	success, err := s.cartRepository.DeletePermanent(cart_id)
	if err != nil {
		s.logger.Error("Failed to permanently delete cart", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete cart"}
	}

	return success, nil
}

func (s *cartService) DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting cart all", zap.Any("cart_id", req.CartIds))

	success, err := s.cartRepository.DeleteAllPermanently(req)
	if err != nil {
		s.logger.Error("Failed to permanently delete cart all", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete cart all"}
	}

	return success, nil
}
