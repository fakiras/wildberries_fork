package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"wildberries/internal/repository"
	"wildberries/internal/service"
)

type AdminHandler struct {
	promotion *service.PromotionService
	segment   *service.SegmentService
	slot      *service.SlotService
	moderation *service.ModerationService
	ident     *service.IdentificationService
	ai        *service.AIService
}

func NewAdminHandler(
	promotion *service.PromotionService,
	segment *service.SegmentService,
	slot *service.SlotService,
	moderation *service.ModerationService,
	ident *service.IdentificationService,
	ai *service.AIService,
) *AdminHandler {
	return &AdminHandler{
		promotion:  promotion,
		segment:    segment,
		slot:       slot,
		moderation: moderation,
		ident:      ident,
		ai:         ai,
	}
}

// CreatePromotion — POST /admin/promotions
func (h *AdminHandler) CreatePromotion(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name                string   `json:"name"`
		Description         string   `json:"description"`
		Theme               string   `json:"theme"`
		DateFrom            string   `json:"date_from"`
		DateTo              string   `json:"date_to"`
		IdentificationMode  string   `json:"identification_mode"`
		PricingModel        string   `json:"pricing_model"`
		SlotCount           int      `json:"slot_count"`
		MinDiscount         *int     `json:"min_discount"`
		MaxDiscount         *int     `json:"max_discount"`
		StopFactors         []string `json:"stop_factors"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	row := &repository.PromotionRow{
		Name:               body.Name,
		Description:        body.Description,
		Theme:              body.Theme,
		DateFrom:           body.DateFrom,
		DateTo:             body.DateTo,
		Status:             "NOT_READY",
		IdentificationMode: body.IdentificationMode,
		PricingModel:       body.PricingModel,
		SlotCount:          body.SlotCount,
		MinDiscount:        body.MinDiscount,
		MaxDiscount:        body.MaxDiscount,
		StopFactors:        mustMarshal(body.StopFactors),
	}
	id, err := h.promotion.Create(r.Context(), row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "NOT_READY"})
}

// GetPromotion — GET /admin/promotions/{id}
func (h *AdminHandler) GetPromotion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	promo, segments, err := h.promotion.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if promo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := promotionToAdminJSON(promo, segments)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// UpdatePromotion — PATCH /admin/promotions/{id}
func (h *AdminHandler) UpdatePromotion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	row := &repository.PromotionRow{ID: id}
	applyPromotionUpdate(row, body)
	if err := h.promotion.Update(r.Context(), row); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DeletePromotion — DELETE /admin/promotions/{id}
func (h *AdminHandler) DeletePromotion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := h.promotion.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// SetFixedPrices — PUT /admin/promotions/{id}/fixed-prices
func (h *AdminHandler) SetFixedPrices(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Prices []struct {
			Position int   `json:"position"`
			Price    int64 `json:"price"`
		} `json:"prices"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	prices := make(map[int32]int64)
	for _, p := range body.Prices {
		prices[int32(p.Position)] = p.Price
	}
	if err := h.promotion.SetFixedPrices(r.Context(), id, prices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ChangeStatus — PUT /admin/promotions/{id}/status
func (h *AdminHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	if err := h.promotion.ChangeStatus(r.Context(), id, body.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// SetSlotProduct — POST /horoscope/products (segmentId, slotId?, productId)
func (h *AdminHandler) SetSlotProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SegmentID int64 `json:"segment_id"`
		SlotID    int64 `json:"slot_id"`
		ProductID int64 `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	if err := h.slot.SetProduct(r.Context(), body.SegmentID, body.SlotID, body.ProductID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetModerationApplications — GET /admin/promotions/{id}/moderation/applications
func (h *AdminHandler) GetModerationApplications(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	status := r.URL.Query().Get("status")
	apps, err := h.moderation.GetApplications(r.Context(), id, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]map[string]interface{}, 0, len(apps))
	for _, a := range apps {
		resp = append(resp, moderationToJSON(a))
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// ApproveModeration — POST /admin/moderation/{applicationId}/approve
func (h *AdminHandler) ApproveModeration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "applicationId"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := h.moderation.Approve(r.Context(), id, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RejectModeration — POST /admin/moderation/{applicationId}/reject
func (h *AdminHandler) RejectModeration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "applicationId"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := h.moderation.Reject(r.Context(), id, body.Reason, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func promotionToAdminJSON(promo *repository.PromotionRow, segments []*repository.SegmentRow) map[string]interface{} {
	segMaps := make([]map[string]interface{}, 0, len(segments))
	for _, s := range segments {
		segMaps = append(segMaps, map[string]interface{}{
			"id": s.ID, "name": s.Name, "category_name": s.CategoryName, "order_index": s.OrderIndex,
		})
	}
	return map[string]interface{}{
		"id":                   promo.ID,
		"name":                 promo.Name,
		"description":          promo.Description,
		"theme":                promo.Theme,
		"status":               promo.Status,
		"date_from":            promo.DateFrom,
		"date_to":              promo.DateTo,
		"identification_mode": promo.IdentificationMode,
		"pricing_model":        promo.PricingModel,
		"slot_count":           promo.SlotCount,
		"min_discount":         promo.MinDiscount,
		"max_discount":         promo.MaxDiscount,
		"stop_factors":         promo.StopFactors,
		"segments":             segMaps,
	}
}

func applyPromotionUpdate(row *repository.PromotionRow, body map[string]interface{}) {}

func moderationToJSON(a *repository.ModerationRow) map[string]interface{} {
	return map[string]interface{}{
		"id": a.ID, "seller_id": a.SellerID, "segment_id": a.SegmentID, "slot_id": a.SlotID,
		"product_name": "", "price": int64(0), "status": a.Status, "discount": a.Discount, "stop_factors": a.StopFactors,
	}
}
