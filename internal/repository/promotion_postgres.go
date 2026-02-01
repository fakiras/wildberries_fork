package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PromotionPostgres struct {
	pool *pgxpool.Pool
}

func NewPromotionPostgres(pool *pgxpool.Pool) *PromotionPostgres {
	return &PromotionPostgres{pool: pool}
}

func (r *PromotionPostgres) GetByID(ctx context.Context, id int64) (*PromotionRow, error) {
	// TODO: SELECT ... FROM promotion WHERE id = $1 AND deleted_at IS NULL
	_ = id
	return nil, nil
}

func (r *PromotionPostgres) GetActive(ctx context.Context) (*PromotionRow, error) {
	// TODO: SELECT ... WHERE status = 'RUNNING' AND date_from <= now() AND date_to >= now() AND deleted_at IS NULL
	return nil, nil
}

func (r *PromotionPostgres) Create(ctx context.Context, row *PromotionRow) (int64, error) {
	// TODO: INSERT INTO promotion (...) RETURNING id
	_ = row
	return 0, nil
}

func (r *PromotionPostgres) Update(ctx context.Context, row *PromotionRow) error {
	// TODO: UPDATE promotion SET ... WHERE id = $1
	_ = row
	return nil
}

func (r *PromotionPostgres) SoftDelete(ctx context.Context, id int64) error {
	// TODO: UPDATE promotion SET deleted_at = now() WHERE id = $1
	_ = id
	return nil
}

func (r *PromotionPostgres) SetFixedPrices(ctx context.Context, id int64, prices []byte) error {
	// TODO: UPDATE promotion SET fixed_prices = $1 WHERE id = $2
	_, _ = id, prices
	return nil
}

func (r *PromotionPostgres) SetStatus(ctx context.Context, id int64, status string) error {
	// TODO: UPDATE promotion SET status = $1 WHERE id = $2
	_, _ = id, status
	return nil
}

// Ensure interface
var _ PromotionRepository = (*PromotionPostgres)(nil)

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
