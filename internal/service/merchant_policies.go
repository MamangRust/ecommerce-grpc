package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type merchantPoliciesService struct {
	merchantPoliciesRepository repository.MerchantPoliciesRepository
	logger                     logger.LoggerInterface
	mapping                    response_service.MerchantPolicyResponseMapper
}

func NewMerchantPoliciesService(
	merchantPoliciesRepository repository.MerchantPoliciesRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantPolicyResponseMapper,
) *merchantPoliciesService {
	return &merchantPoliciesService{
		merchantPoliciesRepository: merchantPoliciesRepository,
		logger:                     logger,
		mapping:                    mapping,
	}
}

func (s *merchantPoliciesService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantPoliciesRepository.FindAllMerchantPolicy(req)

	if err != nil {
		s.logger.Error("Failed to fetch merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve merchants list",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsPolicyResponse(merchants), totalRecords, nil
}

func (s *merchantPoliciesService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantPoliciesRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantPoliciesService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantPoliciesRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantPoliciesService) FindById(merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPoliciesRepository.FindById(merchantID)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve merchant details",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPoliciesService) CreateMerchant(req *requests.CreateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantPoliciesRepository.CreateMerchantPolicy(req)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create new merchant record",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPoliciesService) UpdateMerchant(req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantPolicy", *req.MerchantPolicyID))

	merchant, err := s.merchantPoliciesRepository.UpdateMerchantPolicy(req)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update merchant record",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPoliciesService) TrashedMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPoliciesRepository.TrashedMerchantPolicy(merchantID)

	if err != nil {
		s.logger.Error("Failed to move merchant to trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move merchant to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToMerchantPolicyResponseDeleteAt(merchant), nil
}

func (s *merchantPoliciesService) RestoreMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPoliciesRepository.RestoreMerchantPolicy(merchantID)

	if err != nil {
		s.logger.Error("Failed to restore merchant from trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore merchant from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	return s.mapping.ToMerchantPolicyResponseDeleteAt(merchant), nil
}

func (s *merchantPoliciesService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantPoliciesRepository.DeleteMerchantPolicyPermanent(merchantID)

	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *merchantPoliciesService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantPoliciesRepository.RestoreAllMerchantPolicy()

	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all trashed merchants",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}

func (s *merchantPoliciesService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantPoliciesRepository.DeleteAllMerchantPolicyPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all trashed merchants",
			Code:    http.StatusInternalServerError,
		}
	}

	return success, nil
}
