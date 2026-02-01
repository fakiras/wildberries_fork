package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BetPostgres struct {
	pool *pgxpool.Pool
}

func NewBetPostgres(pool *pgxpool.Pool) *BetPostgres {
	return &BetPostgres{pool: pool}
}

func (r *BetPostgres) Create(ctx context.Context, auctionID, productID int64, bet int64) (int64, error) {
	_, _, _ = auctionID, productID, bet
	return 0, nil
}

func (r *BetPostgres) TopByAuction(ctx context.Context, auctionID int64) (productID int64, bet int64, err error) {
	_ = auctionID
	return 0, 0, nil
}
