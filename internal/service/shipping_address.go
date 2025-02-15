package service

import (
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type shippingAddressService struct {
	shippingRepository repository.ShippingAddressRepository
	logger             logger.LoggerInterface
	mapping            response_service.ShippingAddressResponseMapper
}

func NewShippingAddressService(
	shippingRepository repository.ShippingAddressRepository,
	logger logger.LoggerInterface,
	mapping response_service.ShippingAddressResponseMapper,
) *shippingAddressService {
	return &shippingAddressService{
		shippingRepository: shippingRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *shippingAddressService) FindAll(page int, pageSize int, search string) ([]*response.ShippingAddressResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching category",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	shipping, totalRecords, err := s.shippingRepository.FindAllShippingAddress(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch category",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch category",
		}
	}

	shippingRes := s.mapping.ToShippingAddressesResponse(shipping)

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingRes, int(totalRecords), nil
}

func (s *shippingAddressService) FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching shipping address by ID", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingRepository.FindById(shipping_id)
	if err != nil {
		s.logger.Error("Failed to fetch shipping address", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "shipping address not found"}
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressService) FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching shipping address by order id", zap.Int("shipping_id", order_id))

	shipping, err := s.shippingRepository.FindByOrder(order_id)
	if err != nil {
		s.logger.Error("Failed to fetch shipping address", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "shipping address not found"}
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressService) FindByActive(search string, page, pageSize int) ([]*response.ShippingAddressResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	cashiers, totalRecords, err := s.shippingRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch shipping address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch active shipping address"}
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *shippingAddressService) FindByTrashed(search string, page, pageSize int) ([]*response.ShippingAddressResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching shipping address",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	shipping, totalRecords, err := s.shippingRepository.FindByTrashed(search, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch shipping address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))
		return nil, 0, &response.ErrorResponse{Status: "error", Message: "Failed to fetch trashed shipping address"}
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(shipping), totalRecords, nil
}

func (s *shippingAddressService) TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing shipping address", zap.Int("category", shipping_id))

	category, err := s.shippingRepository.TrashShippingAddress(shipping_id)
	if err != nil {
		s.logger.Error("Failed to trash shipping address", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to trash shipping address"}
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(category), nil
}

func (s *shippingAddressService) RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring shipping address", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingRepository.RestoreShippingAddress(shipping_id)
	if err != nil {
		s.logger.Error("Failed to restore category", zap.Error(err))
		return nil, &response.ErrorResponse{Status: "error", Message: "Failed to restore category"}
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(shipping), nil
}

func (s *shippingAddressService) DeleteShippingAddressPermanently(categoryID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting category", zap.Int("categoryID", categoryID))

	success, err := s.shippingRepository.DeleteShippingAddressPermanently(categoryID)
	if err != nil {
		s.logger.Error("Failed to permanently delete category", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete category"}
	}

	return success, nil
}

func (s *shippingAddressService) RestoreAllShippingAddress() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed shipping address")

	success, err := s.shippingRepository.RestoreAllShippingAddress()
	if err != nil {
		s.logger.Error("Failed to restore all shipping address", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to restore all shipping address"}
	}

	return success, nil
}

func (s *shippingAddressService) DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all shipping address")

	success, err := s.shippingRepository.DeleteAllPermanentShippingAddress()
	if err != nil {
		s.logger.Error("Failed to permanently delete all shipping address", zap.Error(err))
		return false, &response.ErrorResponse{Status: "error", Message: "Failed to permanently delete all shipping address"}
	}

	return success, nil
}
