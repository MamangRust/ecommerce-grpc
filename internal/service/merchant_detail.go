package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type merchantDetailService struct {
	merchantDetailRepository     repository.MerchantDetailRepository
	merchantSocialLinkRepository repository.MerchantSocialLinkRepository
	logger                       logger.LoggerInterface
	mapping                      response_service.MerchantDetailResponseMapper
}

func NewMerchantDetailService(
	merchantDetailRepository repository.MerchantDetailRepository,
	merchantSocialLinkRepository repository.MerchantSocialLinkRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantDetailResponseMapper,
) *merchantDetailService {
	return &merchantDetailService{
		merchantDetailRepository:     merchantDetailRepository,
		merchantSocialLinkRepository: merchantSocialLinkRepository,
		logger:                       logger,
		mapping:                      mapping,
	}
}

func (s *merchantDetailService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantDetailRepository.FindAllMerchants(req)

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

	return s.mapping.ToMerchantsDetailResponse(merchants), totalRecords, nil
}

func (s *merchantDetailService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantDetailRepository.FindByActive(req)

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

	return s.mapping.ToMerchantsDetailResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantDetailService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse) {
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

	merchants, totalRecords, err := s.merchantDetailRepository.FindByTrashed(req)

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

	return s.mapping.ToMerchantsDetailResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantDetailService) FindById(merchantID int) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantDetailRepository.FindById(merchantID)

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

	return s.mapping.ToMerchantDetailRelationResponse(merchant), nil
}

func (s *merchantDetailService) CreateMerchant(req *requests.CreateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantDetailRepository.CreateMerchantDetail(req)
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

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.CreateSocialLink(social)
		if err != nil {
			s.logger.Error("Failed to create social media link",
				zap.Error(err),
				zap.Any("social_link", social))

			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to create merchant social media link",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	return s.mapping.ToMerchantDetailResponse(merchant), nil
}

func (s *merchantDetailService) UpdateMerchant(req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantDetailID))

	merchant, err := s.merchantDetailRepository.UpdateMerchantDetail(req)
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

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.UpdateSocialLink(social)
		if err != nil {
			s.logger.Error("Failed to update social media link",
				zap.Error(err),
				zap.Any("social_link", social))

			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to update merchant social media link",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	return s.mapping.ToMerchantDetailResponse(merchant), nil
}

func (s *merchantDetailService) TrashedMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantDetailRepository.TrashedMerchantDetail(merchantID)
	if err != nil {
		s.logger.Error("Failed to move merchant to trash", zap.Error(err), zap.Int("merchant_id", merchantID))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move merchant to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.merchantSocialLinkRepository.TrashSocialLink(merchant.ID)
	if err != nil {
		s.logger.Debug("Failed to trash merchant social link", zap.Error(err), zap.Int("merchant_id", merchantID))
	}

	return s.mapping.ToMerchantDetailResponseDeleteAt(merchant), nil
}

func (s *merchantDetailService) RestoreMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantDetailRepository.RestoreMerchantDetail(merchantID)
	if err != nil {
		s.logger.Error("Failed to restore merchant from trash", zap.Error(err), zap.Int("merchant_id", merchantID))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore merchant from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.merchantSocialLinkRepository.RestoreSocialLink(merchant.ID)
	if err != nil {
		s.logger.Debug("Failed to restore merchant social link", zap.Error(err), zap.Int("merchant_id", merchantID))
	}

	return s.mapping.ToMerchantDetailResponseDeleteAt(merchant), nil
}

func (s *merchantDetailService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	res, err := s.merchantDetailRepository.FindByIdTrashed(merchantID)
	if err != nil {
		s.logger.Error("Failed to find merchant detail",
			zap.Int("merchant_detail_id", merchantID),
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to verify merchant detail existence",
			Code:    http.StatusInternalServerError,
		}
	}

	if res.CoverImageUrl != "" {
		err := os.Remove(res.CoverImageUrl)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("Cover image file not found, skipping delete",
					zap.String("cover_image_path", res.CoverImageUrl))
			} else {
				s.logger.Error("Failed to delete cover image",
					zap.String("cover_image_path", res.CoverImageUrl),
					zap.Error(err))
				return false, &response.ErrorResponse{
					Status:  "error",
					Message: "Failed to delete merchant cover image",
					Code:    http.StatusInternalServerError,
				}
			}
		} else {
			s.logger.Debug("Successfully deleted cover image",
				zap.String("cover_image_path", res.CoverImageUrl))
		}
	}

	if res.LogoUrl != "" {
		err := os.Remove(res.LogoUrl)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("Logo file not found, skipping delete",
					zap.String("logo_path", res.LogoUrl))
			} else {
				s.logger.Error("Failed to delete logo image",
					zap.String("logo_path", res.LogoUrl),
					zap.Error(err))
				return false, &response.ErrorResponse{
					Status:  "error",
					Message: "Failed to delete merchant logo",
					Code:    http.StatusInternalServerError,
				}
			}
		} else {
			s.logger.Debug("Successfully deleted logo image",
				zap.String("logo_path", res.LogoUrl))
		}
	}

	success, err := s.merchantDetailRepository.DeleteMerchantDetailPermanent(merchantID)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant", zap.Error(err), zap.Int("merchant_id", merchantID))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete merchant",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.merchantSocialLinkRepository.DeletePermanentSocialLink(merchantID)
	if err != nil {
		s.logger.Debug("Failed to permanently delete merchant social link", zap.Error(err), zap.Int("merchant_id", merchantID))
	}

	return success, nil
}

func (s *merchantDetailService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantDetailRepository.RestoreAllMerchantDetail()
	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all trashed merchants",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.merchantSocialLinkRepository.RestoreAllSocialLink()
	if err != nil {
		s.logger.Debug("Failed to restore all social links", zap.Error(err))
	}

	return success, nil
}

func (s *merchantDetailService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantDetailRepository.DeleteAllMerchantDetailPermanent()
	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all trashed merchants",
			Code:    http.StatusInternalServerError,
		}
	}

	_, err = s.merchantSocialLinkRepository.DeleteAllPermanentSocialLink()
	if err != nil {
		s.logger.Debug("Failed to delete all social links permanently", zap.Error(err))
	}

	return success, nil
}
