package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"wildberries/internal/repository"
)

// GenerateSegments — POST /admin/promotions/{id}/segments/generate
func (h *AdminHandler) GenerateSegments(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		UseTheme bool `json:"use_theme"`
		Limit    int  `json:"limit"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	// TODO: call AI, then segment.Create for each
	segments, _ := h.segment.ByPromotionID(r.Context(), id)
	out := make([]map[string]interface{}, 0, len(segments))
	for _, s := range segments {
		out = append(out, map[string]interface{}{"id": s.ID, "name": s.Name, "category_name": s.CategoryName})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"segments": out})
}

// CreateSegment — POST /admin/promotions/{id}/segments
func (h *AdminHandler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Name         string `json:"name"`
		CategoryName string `json:"category_name"`
		OrderIndex   int    `json:"order_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	row := &repository.SegmentRow{
		PromotionID: id,
		Name:        body.Name,
		CategoryName: strPtr(body.CategoryName),
		OrderIndex:  body.OrderIndex,
	}
	segID, err := h.segment.Create(r.Context(), row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": segID, "name": body.Name, "category_name": body.CategoryName})
}

// UpdateSegment — PATCH /admin/promotions/{id}/segments/{segmentId}
func (h *AdminHandler) UpdateSegment(w http.ResponseWriter, r *http.Request) {
	promoID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	segID, _ := strconv.ParseInt(chi.URLParam(r, "segmentId"), 10, 64)
	if promoID == 0 || segID == 0 {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	var body struct {
		Name         *string `json:"name"`
		CategoryName *string `json:"category_name"`
		OrderIndex   *int    `json:"order_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	seg, _ := h.segment.GetByPromoAndSegment(r.Context(), promoID, segID)
	if seg == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if body.Name != nil {
		seg.Name = *body.Name
	}
	if body.CategoryName != nil {
		seg.CategoryName = body.CategoryName
	}
	if body.OrderIndex != nil {
		seg.OrderIndex = *body.OrderIndex
	}
	if err := h.segment.Update(r.Context(), seg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DeleteSegment — DELETE /admin/promotions/{id}/segments/{segmentId}
func (h *AdminHandler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	segID, _ := strconv.ParseInt(chi.URLParam(r, "segmentId"), 10, 64)
	if segID == 0 {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	if err := h.segment.Delete(r.Context(), segID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ShuffleSegmentCategories — POST /admin/promotions/{id}/segments/shuffle-categories
func (h *AdminHandler) ShuffleSegmentCategories(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := h.segment.ShuffleCategories(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
