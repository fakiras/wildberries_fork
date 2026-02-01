package service

import (
	"context"

	"wildberries/internal/repository"
)

type PromotionService struct {
	repo   repository.PromotionRepository
	segRepo repository.SegmentRepository
}

func NewPromotionService(repo repository.PromotionRepository, segRepo repository.SegmentRepository) *PromotionService {
	return &PromotionService{repo: repo, segRepo: segRepo}
}

// GetActive возвращает текущую активную акцию (status=RUNNING, даты пересекаются с now).
func (s *PromotionService) GetActive(ctx context.Context) (*repository.PromotionRow, []*repository.SegmentRow, error) {
	promo, err := s.repo.GetActive(ctx)
	if err != nil || promo == nil {
		return nil, nil, err
	}
	segments, err := s.segRepo.ByPromotionID(ctx, promo.ID)
	if err != nil {
		return promo, nil, err
	}
	return promo, segments, nil
}

// GetByID — для админки.
func (s *PromotionService) GetByID(ctx context.Context, id int64) (*repository.PromotionRow, []*repository.SegmentRow, error) {
	promo, err := s.repo.GetByID(ctx, id)
	if err != nil || promo == nil {
		return nil, nil, err
	}
	segments, err := s.segRepo.ByPromotionID(ctx, id)
	if err != nil {
		return promo, nil, err
	}
	return promo, segments, nil
}

// Create создаёт акцию в статусе NOT_READY.
func (s *PromotionService) Create(ctx context.Context, row *repository.PromotionRow) (int64, error) {
	return s.repo.Create(ctx, row)
}

// Update — частичное обновление.
func (s *PromotionService) Update(ctx context.Context, row *repository.PromotionRow) error {
	return s.repo.Update(ctx, row)
}

// Delete — мягкое удаление.
func (s *PromotionService) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}

// SetFixedPrices сохраняет цены по позициям в promotion.fixed_prices (jsonb).
func (s *PromotionService) SetFixedPrices(ctx context.Context, id int64, prices map[int32]int64) error {
	// Сериализовать prices в jsonb и вызвать repo.SetFixedPrices
	return s.repo.SetFixedPrices(ctx, id, nil) // TODO: marshal prices
}

// ChangeStatus меняет статус; при переходе в READY_TO_START создаются аукционы для слотов.
func (s *PromotionService) ChangeStatus(ctx context.Context, id int64, status string) error {
	// TODO: при status == READY_TO_START — создать записи auction для каждого слота с pricing_type=auction
	return s.repo.SetStatus(ctx, id, status)
}
