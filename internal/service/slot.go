package service

import (
	"context"

	"wildberries/internal/repository"
)

type SlotService struct {
	slotRepo   repository.SlotRepository
	modRepo    repository.ModerationRepository
	auctionRepo repository.AuctionRepository
	betRepo    repository.BetRepository
}

func NewSlotService(
	slotRepo repository.SlotRepository,
	modRepo repository.ModerationRepository,
	auctionRepo repository.AuctionRepository,
	betRepo repository.BetRepository,
) *SlotService {
	return &SlotService{
		slotRepo:   slotRepo,
		modRepo:    modRepo,
		auctionRepo: auctionRepo,
		betRepo:    betRepo,
	}
}

// GetSlotsBySegment возвращает слоты сегмента (опционально только occupied).
func (s *SlotService) GetSlotsBySegment(ctx context.Context, segmentID int64, onlyOccupied bool) ([]*repository.SlotRow, error) {
	return s.slotRepo.BySegmentID(ctx, segmentID, onlyOccupied)
}

// GetSellerSlots возвращает слоты/ставки/заявки селлера.
func (s *SlotService) GetSellerSlots(ctx context.Context, sellerID int64, promotionID *int64) ([]*repository.SlotRow, error) {
	return s.slotRepo.BySellerID(ctx, sellerID, promotionID)
}

// Participate — ставка (аукцион) или покупка слота (fixed) с созданием заявки на модерацию.
func (s *SlotService) Participate(ctx context.Context, sellerID int64, slotID int64, amount *int64, product *ProductInput) error {
	// TODO: если слот auction — INSERT bet; если fixed — UPDATE slot status=pending, INSERT moderation
	_, _, _, _ = sellerID, slotID, amount, product
	return nil
}

type ProductInput struct {
	Name    string
	Price   int64
	Discount int
	Image   string
}

// Cancel — отмена ставки/заявки.
func (s *SlotService) Cancel(ctx context.Context, sellerID int64, slotID int64) error {
	_, _ = sellerID, slotID
	return nil
}

// SetProduct — ручная установка товара в слот (админ).
func (s *SlotService) SetProduct(ctx context.Context, segmentID, slotID int64, productID int64) error {
	_, _, _ = segmentID, slotID, productID
	return nil
}
