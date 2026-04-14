package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	merchantsociallink_errors "ecommerce/pkg/errors/merchant_social_link_errors"
)

type merchantSocialMediaLinkRepository struct {
	db *db.Queries
}

func NewMerchantSocialLinkRepository(db *db.Queries) *merchantSocialMediaLinkRepository {
	return &merchantSocialMediaLinkRepository{
		db: db,
	}
}

func (r *merchantSocialMediaLinkRepository) CreateSocialLink(ctx context.Context, req *requests.CreateMerchantSocialRequest) (bool, error) {
	params := db.CreateMerchantSocialMediaLinkParams{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.CreateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrCreateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) UpdateSocialLink(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (bool, error) {
	params := db.UpdateMerchantSocialMediaLinkParams{
		MerchantSocialID: int32(req.ID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.UpdateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrUpdateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) TrashSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrTrashMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) RestoreSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrDeletePermanentMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) RestoreAllSocialLink(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreAllMerchantSocialLinks
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) DeleteAllPermanentSocialLink(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrDeleteAllPermanentMerchantSocialLinks
	}

	return true, nil
}
