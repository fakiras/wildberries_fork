package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"wildberries/internal/repository"
	"wildberries/internal/service"
)

type BuyerHandler struct {
	promotion *service.PromotionService
	segment   *service.SegmentService
	slot      *service.SlotService
	product   *service.ProductService
	ident     *service.IdentificationService
}

func NewBuyerHandler(
	promotion *service.PromotionService,
	segment *service.SegmentService,
	slot *service.SlotService,
	product *service.ProductService,
	ident *service.IdentificationService,
) *BuyerHandler {
	return &BuyerHandler{
		promotion: promotion,
		segment:   segment,
		slot:      slot,
		product:   product,
		ident:     ident,
	}
}

// GetCurrentPromotion — GET /promotions/current
func (h *BuyerHandler) GetCurrentPromotion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	promo, segments, err := h.promotion.GetActive(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if promo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"id":                  promo.ID,
		"name":                promo.Name,
		"description":         promo.Description,
		"theme":               promo.Theme,
		"status":              promo.Status,
		"date_from":           promo.DateFrom,
		"date_to":             promo.DateTo,
		"segments":            segmentsToJSON(segments),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetSegmentProducts — GET /promotions/{promotionId}/segments/{segmentId}/products
func (h *BuyerHandler) GetSegmentProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	promoID, _ := strconv.ParseInt(chi.URLParam(r, "promotionId"), 10, 64)
	segID, _ := strconv.ParseInt(chi.URLParam(r, "segmentId"), 10, 64)
	if promoID == 0 || segID == 0 {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	_, err := h.segment.GetByPromoAndSegment(r.Context(), promoID, segID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	slots, err := h.slot.GetSlotsBySegment(r.Context(), segID, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var productIDs []int64
	for _, s := range slots {
		if s.ProductID != nil {
			productIDs = append(productIDs, *s.ProductID)
		}
	}
	filters := repository.ProductFilters{
		Category:      r.URL.Query().Get("category"),
		OnlyDiscounts: r.URL.Query().Get("onlyDiscounts") == "true",
		Sort:          r.URL.Query().Get("sort"),
	}
	products, err := h.product.GetByIDs(r.Context(), productIDs, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page, perPage := getPagePerPage(r)
	resp := map[string]interface{}{
		"items":     productsToItems(products),
		"total":     len(products),
		"page":      page,
		"per_page":  perPage,
		"completed": false,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// StartIdentification — POST /identification/start
func (h *BuyerHandler) StartIdentification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		PromotionID int64 `json:"promotionId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	method, questions, optionsByQuestion, resultSegmentID, err := h.ident.StartIdentification(r.Context(), body.PromotionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{"method": method}
	if method == "user_profile" && resultSegmentID != 0 {
		resp["result_segment_id"] = resultSegmentID
	}
	if method == "questions" && len(questions) > 0 {
		resp["poll"] = pollToJSON(questions, optionsByQuestion)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Answer — POST /identification/answer
func (h *BuyerHandler) Answer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		PromotionID int64 `json:"promotionId"`
		QuestionID  int64 `json:"questionId"`
		OptionID    int64 `json:"optionId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	nextID, resultSegID, err := h.ident.ProcessAnswer(r.Context(), body.PromotionID, body.QuestionID, body.OptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"next_question_id":   nextID,
		"result_segment_id":  resultSegID,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func segmentsToJSON(segments []*repository.SegmentRow) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(segments))
	for _, s := range segments {
		out = append(out, map[string]interface{}{
			"id":            s.ID,
			"name":          s.Name,
			"category_name": s.CategoryName,
		})
	}
	return out
}

func productsToItems(products []*repository.ProductRow) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		oldPrice := int64(0)
		if p.Discount > 0 {
			oldPrice = p.Price * int64(100) / int64(100-p.Discount)
		}
		out = append(out, map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"image":      p.Image,
			"price":      p.Price,
			"old_price":  oldPrice,
			"discount":   p.Discount,
			"badge":      "",
		})
	}
	return out
}

func pollToJSON(questions []*repository.PollQuestionRow, optionsByQuestion map[int64][]*repository.PollOptionRow) map[string]interface{} {
	qs := make([]map[string]interface{}, 0, len(questions))
	for _, q := range questions {
		opts := optionsByQuestion[q.ID]
		optMaps := make([]map[string]interface{}, 0, len(opts))
		for _, o := range opts {
			optMaps = append(optMaps, map[string]interface{}{"id": o.ID, "text": o.Text, "value": o.Value})
		}
		qs = append(qs, map[string]interface{}{"id": q.ID, "text": q.Text, "options": optMaps})
	}
	return map[string]interface{}{"questions": qs}
}

func getPagePerPage(r *http.Request) (page, perPage int) {
	page, perPage = 1, 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, _ := strconv.Atoi(p); v > 0 {
			page = v
		}
	}
	if p := r.URL.Query().Get("per_page"); p != "" {
		if v, _ := strconv.Atoi(p); v > 0 {
			perPage = v
		}
	}
	return page, perPage
}
