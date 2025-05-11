package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	merchant_errors "ecommerce/pkg/errors/merchant"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type merchantService struct {
	merchantRepository repository.MerchantRepository
	logger             logger.LoggerInterface
	mapping            response_service.MerchantResponseMapper
}

func NewMerchantService(
	merchantRepository repository.MerchantRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantResponseMapper,
) *merchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *merchantService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindAllMerchants(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, merchant_errors.ErrFailedFindAllMerchants
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsResponse(merchants), totalRecords, nil
}

func (s *merchantService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchants active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchant_errors.ErrFailedFindActiveMerchants
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all merchants trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchant_errors.ErrFailedFindTrashedMerchants
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantService) FindById(merchantID int) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantRepository.FindById(merchantID)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	return s.mapping.ToMerchantResponse(merchant), nil
}

func (s *merchantService) CreateMerchant(req *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantRepository.CreateMerchant(req)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchant_errors.ErrFailedCreateMerchant
	}

	return s.mapping.ToMerchantResponse(merchant), nil
}

func (s *merchantService) UpdateMerchant(req *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantID))

	merchant, err := s.merchantRepository.UpdateMerchant(req)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchant_errors.ErrFailedUpdateMerchant
	}

	return s.mapping.ToMerchantResponse(merchant), nil
}

func (s *merchantService) TrashedMerchant(merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantRepository.TrashedMerchant(merchantID)

	if err != nil {
		s.logger.Error("Failed to move merchant to trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchant_errors.ErrFailedTrashedMerchant
	}

	return s.mapping.ToMerchantResponseDeleteAt(merchant), nil
}

func (s *merchantService) RestoreMerchant(merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantRepository.RestoreMerchant(merchantID)

	if err != nil {
		s.logger.Error("Failed to restore merchant from trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchant_errors.ErrFailedRestoreMerchant
	}

	return s.mapping.ToMerchantResponseDeleteAt(merchant), nil
}

func (s *merchantService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantRepository.DeleteMerchantPermanent(merchantID)

	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return false, merchant_errors.ErrFailedDeleteMerchantPermanent
	}

	return success, nil
}

func (s *merchantService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantRepository.RestoreAllMerchant()

	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err))

		return false, merchant_errors.ErrFailedRestoreAllMerchants
	}

	return success, nil
}

func (s *merchantService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantRepository.DeleteAllMerchantPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants",
			zap.Error(err))

		return false, merchant_errors.ErrFailedDeleteAllMerchantsPermanent
	}

	return success, nil
}
