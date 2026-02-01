package service

import (
	"context"

	"wildberries/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetByIDs(ctx context.Context, ids []int64, filters repository.ProductFilters) ([]*repository.ProductRow, error) {
	return s.repo.GetByIDs(ctx, ids, filters)
}

func (s *ProductService) ListBySeller(ctx context.Context, sellerID int64, categoryID string, page, perPage int) ([]*repository.ProductRow, int, error) {
	return s.repo.ListBySeller(ctx, sellerID, categoryID, page, perPage)
}
