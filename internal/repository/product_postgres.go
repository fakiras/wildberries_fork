package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductPostgres struct {
	pool *pgxpool.Pool
}

func NewProductPostgres(pool *pgxpool.Pool) *ProductPostgres {
	return &ProductPostgres{pool: pool}
}

func (r *ProductPostgres) GetByIDs(ctx context.Context, ids []int64, filters ProductFilters) ([]*ProductRow, error) {
	_, _ = ids, filters
	return nil, nil
}

func (r *ProductPostgres) ListBySeller(ctx context.Context, sellerID int64, categoryID string, page, perPage int) ([]*ProductRow, int, error) {
	_, _, _, _ = sellerID, categoryID, page, perPage
	return nil, 0, nil
}

var _ ProductRepository = (*ProductPostgres)(nil)
