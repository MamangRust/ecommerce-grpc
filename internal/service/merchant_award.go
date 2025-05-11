package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	merchantaward_errors "ecommerce/pkg/errors/merchant_award"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type merchantAwardService struct {
	merchantAwardRepository repository.MerchantAwardRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.MerchantAwardResponseMapper
}

func NewMerchantAwardService(
	merchantAwardRepository repository.MerchantAwardRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantAwardResponseMapper,
) *merchantAwardService {
	return &merchantAwardService{
		merchantAwardRepository: merchantAwardRepository,
		logger:                  logger,
		mapping:                 mapping,
	}
}

func (s *merchantAwardService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantAwardRepository.FindAllMerchants(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, merchantaward_errors.ErrFailedFindAllMerchantAwards
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsAwardResponse(merchants), totalRecords, nil
}

func (s *merchantAwardService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantAwardRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchantaward_errors.ErrFailedFindActiveMerchantAwards
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsAwardResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantAwardService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantAwardRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchantaward_errors.ErrFailedFindTrashedMerchantAwards
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsAwardResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantAwardService) FindById(merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantAwardRepository.FindById(merchantID)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))
		return nil, merchantaward_errors.ErrFailedFindMerchantAwardById
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardService) CreateMerchant(req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantAwardRepository.CreateMerchantAward(req)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantaward_errors.ErrFailedCreateMerchantAward
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardService) UpdateMerchant(req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantCertificationID))

	merchant, err := s.merchantAwardRepository.UpdateMerchantAward(req)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantaward_errors.ErrFailedUpdateMerchantAward
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardService) TrashedMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantAwardRepository.TrashedMerchantAward(merchantID)

	if err != nil {
		s.logger.Error("Failed to move merchant to trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantaward_errors.ErrFailedTrashedMerchantAward
	}

	return s.mapping.ToMerchantAwardResponseDeleteAt(merchant), nil
}

func (s *merchantAwardService) RestoreMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantAwardRepository.RestoreMerchantAward(merchantID)

	if err != nil {
		s.logger.Error("Failed to restore merchant from trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantaward_errors.ErrFailedRestoreMerchantAward
	}

	return s.mapping.ToMerchantAwardResponseDeleteAt(merchant), nil
}

func (s *merchantAwardService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantAwardRepository.DeleteMerchantPermanent(merchantID)

	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))
		return false, merchantaward_errors.ErrFailedDeleteMerchantAwardPermanent
	}

	return success, nil
}

func (s *merchantAwardService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantAwardRepository.RestoreAllMerchantAward()

	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err))

		return false, merchantaward_errors.ErrFailedRestoreAllMerchantAwards
	}

	return success, nil
}

func (s *merchantAwardService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantAwardRepository.DeleteAllMerchantAwardPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants",
			zap.Error(err))

		return false, merchantaward_errors.ErrFailedDeleteAllMerchantAwardsPermanent
	}

	return success, nil
}
