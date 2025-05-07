package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	shippingaddress_errors "ecommerce/pkg/errors/shipping_address_errors"
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

func (s *shippingAddressService) FindAll(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	shipping, totalRecords, err := s.shippingRepository.FindAllShippingAddress(req)

	if err != nil {
		s.logger.Error("Failed to retrieve shipping address list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, shippingaddress_errors.ErrFailedFindAllShippingAddresses
	}

	shippingRes := s.mapping.ToShippingAddressesResponse(shipping)

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingRes, totalRecords, nil
}

func (s *shippingAddressService) FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching Shipping Address by ID", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingRepository.FindById(shipping_id)

	if err != nil {
		s.logger.Error("Failed to retrieve Shipping Address details",
			zap.Int("Shipping Address ID", shipping_id),
			zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByID
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressService) FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching shipping address by order id", zap.Int("shipping_id", order_id))

	shipping, err := s.shippingRepository.FindByOrder(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve Shipping Address details",
			zap.Int("Order ID", order_id),
			zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressService) FindByActive(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

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

	cashiers, totalRecords, err := s.shippingRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active Shipping Address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, shippingaddress_errors.ErrFailedFindActiveShippingAddresses
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *shippingAddressService) FindByTrashed(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Shipping Address",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	shipping, totalRecords, err := s.shippingRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed shipping address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, shippingaddress_errors.ErrFailedFindTrashedShippingAddresses
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(shipping), totalRecords, nil
}

func (s *shippingAddressService) TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing shipping address", zap.Int("category", shipping_id))

	category, err := s.shippingRepository.TrashShippingAddress(shipping_id)

	if err != nil {
		s.logger.Error("Failed to move shipping address to trash",
			zap.Int("shipping_id", shipping_id),
			zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedTrashShippingAddress
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(category), nil
}

func (s *shippingAddressService) RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring Shipping Address", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingRepository.RestoreShippingAddress(shipping_id)

	if err != nil {
		s.logger.Error("Failed to restore role from trash",
			zap.Int("shipping_id", shipping_id),
			zap.Error(err))

		return nil, shippingaddress_errors.ErrFailedRestoreShippingAddress
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(shipping), nil
}

func (s *shippingAddressService) DeleteShippingAddressPermanently(shipping_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting shipping address", zap.Int("shipping_id", shipping_id))

	success, err := s.shippingRepository.DeleteShippingAddressPermanently(shipping_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete role",
			zap.Int("shipping_address", shipping_id),
			zap.Error(err))

		return false, shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent
	}

	return success, nil
}

func (s *shippingAddressService) RestoreAllShippingAddress() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed shipping address")

	success, err := s.shippingRepository.RestoreAllShippingAddress()

	if err != nil {
		s.logger.Error("Failed to restore all trashed shipping address",
			zap.Error(err))
		return false, shippingaddress_errors.ErrFailedRestoreAllShippingAddresses
	}

	return success, nil
}

func (s *shippingAddressService) DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all shipping address")

	success, err := s.shippingRepository.DeleteAllPermanentShippingAddress()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed shipping address",
			zap.Error(err))

		return false, shippingaddress_errors.ErrFailedDeleteAllShippingAddressesPermanent
	}

	return success, nil
}
