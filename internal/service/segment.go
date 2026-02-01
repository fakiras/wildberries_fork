package service

import (
	"context"

	"wildberries/internal/repository"
)

type SegmentService struct {
	repo repository.SegmentRepository
}

func NewSegmentService(repo repository.SegmentRepository) *SegmentService {
	return &SegmentService{repo: repo}
}

func (s *SegmentService) GetByPromoAndSegment(ctx context.Context, promotionID, segmentID int64) (*repository.SegmentRow, error) {
	return s.repo.GetByPromoAndSegment(ctx, promotionID, segmentID)
}

func (s *SegmentService) ByPromotionID(ctx context.Context, promotionID int64) ([]*repository.SegmentRow, error) {
	return s.repo.ByPromotionID(ctx, promotionID)
}

func (s *SegmentService) Create(ctx context.Context, row *repository.SegmentRow) (int64, error) {
	return s.repo.Create(ctx, row)
}

func (s *SegmentService) Update(ctx context.Context, row *repository.SegmentRow) error {
	return s.repo.Update(ctx, row)
}

func (s *SegmentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *SegmentService) ShuffleCategories(ctx context.Context, promotionID int64) error {
	return s.repo.ShuffleCategories(ctx, promotionID)
}

func (s *SegmentService) UpdateSegmentText(ctx context.Context, segmentID int64, text string) error {
	return s.repo.UpdateText(ctx, segmentID, text)
}
