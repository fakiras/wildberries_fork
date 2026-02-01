package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotPostgres struct {
	pool *pgxpool.Pool
}

func NewSlotPostgres(pool *pgxpool.Pool) *SlotPostgres {
	return &SlotPostgres{pool: pool}
}

func (r *SlotPostgres) BySegmentID(ctx context.Context, segmentID int64, onlyOccupied bool) ([]*SlotRow, error) {
	_, _ = segmentID, onlyOccupied
	return nil, nil
}

func (r *SlotPostgres) BySellerID(ctx context.Context, sellerID int64, promotionID *int64) ([]*SlotRow, error) {
	_, _ = sellerID, promotionID
	return nil, nil
}

func (r *SlotPostgres) GetByID(ctx context.Context, id int64) (*SlotRow, error) {
	_ = id
	return nil, nil
}

func (r *SlotPostgres) Create(ctx context.Context, row *SlotRow) (int64, error) {
	_ = row
	return 0, nil
}

func (r *SlotPostgres) Update(ctx context.Context, row *SlotRow) error {
	_ = row
	return nil
}

func (r *SlotPostgres) SetProduct(ctx context.Context, slotID int64, sellerID, productID int64, status string) error {
	_, _, _, _ = slotID, sellerID, productID, status
	return nil
}

var _ SlotRepository = (*SlotPostgres)(nil)
