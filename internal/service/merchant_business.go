package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	merchantbusiness_errors "ecommerce/pkg/errors/merchant_business"
	"ecommerce/pkg/logger"

	"go.uber.org/zap"
)

type merchantBusinessService struct {
	merchantBusinessRepository repository.MerchantBusinessRepository
	logger                     logger.LoggerInterface
	mapping                    response_service.MerchantBusinessResponseMapper
}

func NewMerchantBusinessService(
	merchantBusinessRepository repository.MerchantBusinessRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantBusinessResponseMapper,
) *merchantBusinessService {
	return &merchantBusinessService{
		merchantBusinessRepository: merchantBusinessRepository,
		logger:                     logger,
		mapping:                    mapping,
	}
}

func (s *merchantBusinessService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantBusinessRepository.FindAllMerchants(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, merchantbusiness_errors.ErrFailedFindAllMerchantBusiness
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsBusinessResponse(merchants), totalRecords, nil
}

func (s *merchantBusinessService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantBusinessRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchantbusiness_errors.ErrFailedFindActiveMerchantBusiness
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsBusinessResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantBusinessService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantBusinessRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, merchantbusiness_errors.ErrFailedFindTrashedMerchantBusiness
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsBusinessResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantBusinessService) FindById(merchantID int) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantBusinessRepository.FindById(merchantID)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantbusiness_errors.ErrFailedFindMerchantBusinessById
	}

	return s.mapping.ToMerchantBusinessResponseRelation(merchant), nil
}

func (s *merchantBusinessService) CreateMerchant(req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantBusinessRepository.CreateMerchantBusiness(req)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantbusiness_errors.ErrFailedCreateMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponse(merchant), nil
}

func (s *merchantBusinessService) UpdateMerchant(req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantBusinessInfoID))

	merchant, err := s.merchantBusinessRepository.UpdateMerchantBusiness(req)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantbusiness_errors.ErrFailedUpdateMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponse(merchant), nil
}

func (s *merchantBusinessService) TrashedMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantBusinessRepository.TrashedMerchantBusiness(merchantID)

	if err != nil {
		s.logger.Error("Failed to move merchant to trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantbusiness_errors.ErrFailedTrashedMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponseDeleteAt(merchant), nil
}

func (s *merchantBusinessService) RestoreMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantBusinessRepository.RestoreMerchantBusiness(merchantID)

	if err != nil {
		s.logger.Error("Failed to restore merchant from trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantbusiness_errors.ErrFailedRestoreMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponseDeleteAt(merchant), nil
}

func (s *merchantBusinessService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantBusinessRepository.DeleteMerchantBusinessPermanent(merchantID)

	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return false, merchantbusiness_errors.ErrFailedDeleteMerchantBusinessPermanent
	}

	return success, nil
}

func (s *merchantBusinessService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantBusinessRepository.RestoreAllMerchantBusiness()

	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err))

		return false, merchantbusiness_errors.ErrFailedRestoreAllMerchantBusiness
	}

	return success, nil
}

func (s *merchantBusinessService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantBusinessRepository.DeleteAllMerchantBusinessPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants",
			zap.Error(err))

		return false, merchantbusiness_errors.ErrFailedDeleteAllMerchantBusinessPermanent
	}

	return success, nil
}
