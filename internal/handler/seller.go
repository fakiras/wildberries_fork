package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wildberries/internal/service"
)

type SellerHandler struct {
	product   *service.ProductService
	promotion *service.PromotionService
	slot      *service.SlotService
}

func NewSellerHandler(
	product *service.ProductService,
	promotion *service.PromotionService,
	slot *service.SlotService,
) *SellerHandler {
	return &SellerHandler{
		product:   product,
		promotion: promotion,
		slot:      slot,
	}
}

// sellerIDFromContext — заглушка: в реальности брать из JWT/сессии (внешний ID из другой системы)
func sellerIDFromContext(r *http.Request) int64 {
	// TODO: X-Seller-Id header или JWT claim
	if id := r.Header.Get("X-Seller-Id"); id != "" {
		if v, _ := strconv.ParseInt(id, 10, 64); v > 0 {
			return v
		}
	}
	return 0
}

// ListProductsBy — GET /products/list-by
func (h *SellerHandler) ListProductsBy(w http.ResponseWriter, r *http.Request) {
	sellerID := sellerIDFromContext(r)
	if sellerID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	categoryID := r.URL.Query().Get("category_id")
	page, perPage := getPagePerPage(r)
	products, total, err := h.product.ListBySeller(r.Context(), sellerID, categoryID, page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	items := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		items = append(items, map[string]interface{}{
			"id": p.ID, "nm_id": p.NmID, "category_name": p.CategoryName,
			"name": p.Name, "image": p.Image, "price": p.Price, "discount": p.Discount,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"items": items, "total": total, "page": page, "per_page": perPage,
	})
}

// GetSellerActions — GET /seller/actions
func (h *SellerHandler) GetSellerActions(w http.ResponseWriter, r *http.Request) {
	// TODO: список акций в статусе RUNNING (или READY_TO_START), доступных селлеру
	_ = r.Context()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"actions": []interface{}{}})
}

// GetSellerBetsList — GET /seller/bets/list
func (h *SellerHandler) GetSellerBetsList(w http.ResponseWriter, r *http.Request) {
	sellerID := sellerIDFromContext(r)
	if sellerID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var promotionID *int64
	if p := r.URL.Query().Get("promotion_id"); p != "" {
		if v, _ := strconv.ParseInt(p, 10, 64); v > 0 {
			promotionID = &v
		}
	}
	slots, err := h.slot.GetSellerSlots(r.Context(), sellerID, promotionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	items := make([]map[string]interface{}, 0, len(slots))
	for _, s := range slots {
		items = append(items, map[string]interface{}{
			"id": s.ID, "promotion_id": s.PromotionID, "segment_id": s.SegmentID, "slot_id": s.ID,
			"status": s.Status,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
}

// MakeBet — POST /seller/bets/make
func (h *SellerHandler) MakeBet(w http.ResponseWriter, r *http.Request) {
	sellerID := sellerIDFromContext(r)
	if sellerID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		SlotID  int64 `json:"slotId"`
		Amount  int64 `json:"amount"`
		Product *struct {
			Name     string `json:"name"`
			Price    int64  `json:"price"`
			Discount int    `json:"discount"`
			Image    string `json:"image"`
		} `json:"product"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	var product *service.ProductInput
	if body.Product != nil {
		product = &service.ProductInput{
			Name: body.Product.Name, Price: body.Product.Price,
			Discount: body.Product.Discount, Image: body.Product.Image,
		}
	}
	var amount *int64
	if body.Amount > 0 {
		amount = &body.Amount
	}
	if err := h.slot.Participate(r.Context(), sellerID, body.SlotID, amount, product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// RemoveBet — POST /seller/bets/remove
func (h *SellerHandler) RemoveBet(w http.ResponseWriter, r *http.Request) {
	sellerID := sellerIDFromContext(r)
	if sellerID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		SlotID int64 `json:"slotId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	if err := h.slot.Cancel(r.Context(), sellerID, body.SlotID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
