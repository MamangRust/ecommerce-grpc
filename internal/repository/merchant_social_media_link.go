package repository

import (
	"context"
	"ecommerce/internal/domain/requests"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

type merchantSocialMediaLinkRepository struct {
	db  *db.Queries
	ctx context.Context
}

func NewMerchantSocialMediaLinkRepository(db *db.Queries, ctx context.Context) *merchantSocialMediaLinkRepository {
	return &merchantSocialMediaLinkRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *merchantSocialMediaLinkRepository) CreateSocialLink(req *requests.CreateMerchantSocialRequest) (bool, error) {
	params := db.CreateMerchantSocialMediaLinkParams{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.CreateMerchantSocialMediaLink(r.ctx, params)
	if err != nil {
		return false, fmt.Errorf("failed to create social media link: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) UpdateSocialLink(req *requests.UpdateMerchantSocialRequest) (bool, error) {
	params := db.UpdateMerchantSocialMediaLinkParams{
		MerchantSocialID: int32(req.ID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.UpdateMerchantSocialMediaLink(r.ctx, params)
	if err != nil {
		return false, fmt.Errorf("failed to update social media link: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) TrashSocialLink(socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(r.ctx, int32(socialID))
	if err != nil {
		return false, fmt.Errorf("failed to trash social media link: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) RestoreSocialLink(socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(r.ctx, int32(socialID))
	if err != nil {
		return false, fmt.Errorf("failed to restore social media link: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) DeletePermanentSocialLink(socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(r.ctx, int32(socialID))
	if err != nil {
		return false, fmt.Errorf("failed to permanently delete social media link: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) RestoreAllSocialLink() (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(r.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to restore all social media links: %w", err)
	}

	return true, nil
}

func (r *merchantSocialMediaLinkRepository) DeleteAllPermanentSocialLink() (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(r.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to permanently delete all social media links: %w", err)
	}

	return true, nil
}
