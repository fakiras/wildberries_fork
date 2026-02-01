package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ModerationPostgres struct {
	pool *pgxpool.Pool
}

func NewModerationPostgres(pool *pgxpool.Pool) *ModerationPostgres {
	return &ModerationPostgres{pool: pool}
}

func (r *ModerationPostgres) ListByPromotion(ctx context.Context, promotionID int64, status string) ([]*ModerationRow, error) {
	_, _ = promotionID, status
	return nil, nil
}

func (r *ModerationPostgres) GetByID(ctx context.Context, id int64) (*ModerationRow, error) {
	_ = id
	return nil, nil
}

func (r *ModerationPostgres) Create(ctx context.Context, row *ModerationRow) (int64, error) {
	_ = row
	return 0, nil
}

func (r *ModerationPostgres) SetStatus(ctx context.Context, id int64, status string, moderatorID *int64) error {
	_, _, _ = id, status, moderatorID
	return nil
}

var _ ModerationRepository = (*ModerationPostgres)(nil)
