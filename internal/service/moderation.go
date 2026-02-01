package service

import (
	"context"

	"wildberries/internal/repository"
)

type ModerationService struct {
	repo    repository.ModerationRepository
	slotRepo repository.SlotRepository
}

func NewModerationService(repo repository.ModerationRepository, slotRepo repository.SlotRepository) *ModerationService {
	return &ModerationService{repo: repo, slotRepo: slotRepo}
}

func (s *ModerationService) GetApplications(ctx context.Context, promotionID int64, status string) ([]*repository.ModerationRow, error) {
	return s.repo.ListByPromotion(ctx, promotionID, status)
}

func (s *ModerationService) Approve(ctx context.Context, applicationID int64, moderatorID *int64) error {
	// TODO: UPDATE moderation SET status='approved'; UPDATE slot SET status='occupied', product_id=...
	return s.repo.SetStatus(ctx, applicationID, "approved", moderatorID)
}

func (s *ModerationService) Reject(ctx context.Context, applicationID int64, reason string, moderatorID *int64) error {
	// TODO: UPDATE moderation SET status='rejected'; UPDATE slot SET status='available', seller_id=NULL, product_id=NULL
	_ = reason
	return s.repo.SetStatus(ctx, applicationID, "rejected", moderatorID)
}
