package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentPostgres struct {
	pool *pgxpool.Pool
}

func NewSegmentPostgres(pool *pgxpool.Pool) *SegmentPostgres {
	return &SegmentPostgres{pool: pool}
}

func (r *SegmentPostgres) ByPromotionID(ctx context.Context, promotionID int64) ([]*SegmentRow, error) {
	_ = promotionID
	return nil, nil
}

func (r *SegmentPostgres) GetByID(ctx context.Context, id int64) (*SegmentRow, error) {
	_ = id
	return nil, nil
}

func (r *SegmentPostgres) GetByPromoAndSegment(ctx context.Context, promotionID, segmentID int64) (*SegmentRow, error) {
	_, _ = promotionID, segmentID
	return nil, nil
}

func (r *SegmentPostgres) Create(ctx context.Context, row *SegmentRow) (int64, error) {
	_ = row
	return 0, nil
}

func (r *SegmentPostgres) Update(ctx context.Context, row *SegmentRow) error {
	_ = row
	return nil
}

func (r *SegmentPostgres) Delete(ctx context.Context, id int64) error {
	_ = id
	return nil
}

func (r *SegmentPostgres) ShuffleCategories(ctx context.Context, promotionID int64) error {
	_ = promotionID
	return nil
}

func (r *SegmentPostgres) UpdateText(ctx context.Context, segmentID int64, text string) error {
	_, _ = segmentID, text
	return nil
}

var _ SegmentRepository = (*SegmentPostgres)(nil)
