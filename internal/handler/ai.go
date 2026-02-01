package handler

import (
	"encoding/json"
	"net/http"

	"wildberries/internal/service"
)

type AIHandler struct {
	ai      *service.AIService
	segment *service.SegmentService
}

func NewAIHandler(ai *service.AIService, segment *service.SegmentService) *AIHandler {
	return &AIHandler{ai: ai, segment: segment}
}

// GenerateThemes — POST /ai/themes
func (h *AIHandler) GenerateThemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	themes, err := h.ai.GenerateThemes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]map[string]interface{}, 0, len(themes))
	for _, t := range themes {
		out = append(out, map[string]interface{}{"value": t.Value, "label": t.Label})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"themes": out})
}

// GenerateSegments — POST /ai/segments
func (h *AIHandler) GenerateSegments(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Theme string `json:"theme"`
		Limit int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	segments, err := h.ai.GenerateSegments(r.Context(), body.Theme, body.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]map[string]interface{}, 0, len(segments))
	for _, s := range segments {
		out = append(out, map[string]interface{}{"name": s.Name, "category_name": s.CategoryName})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"segments": out})
}

// GenerateQuestions — POST /ai/questions
func (h *AIHandler) GenerateQuestions(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Theme string `json:"theme"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	questions, err := h.ai.GenerateQuestions(r.Context(), body.Theme)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]map[string]interface{}, 0, len(questions))
	for _, q := range questions {
		opts := make([]map[string]interface{}, 0, len(q.Options))
		for _, o := range q.Options {
			opts = append(opts, map[string]interface{}{"text": o.Text, "value": o.Value})
		}
		out = append(out, map[string]interface{}{"text": q.Text, "options": opts})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"questions": out})
}

// GenerateAnswerTree — POST /ai/answer-tree
func (h *AIHandler) GenerateAnswerTree(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Theme string `json:"theme"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	nodes, err := h.ai.GenerateAnswerTree(r.Context(), body.Theme)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]map[string]interface{}, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, map[string]interface{}{
			"node_id": n.NodeID, "parent_node_id": n.ParentNodeID, "label": n.Label, "value": n.Value,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"nodes": out})
}

// GetText — POST /ai/get-text
func (h *AIHandler) GetText(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Params    map[string]string `json:"params"`
		SegmentID int64             `json:"segment_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	var segID *int64
	if body.SegmentID > 0 {
		segID = &body.SegmentID
	}
	text, err := h.ai.GenerateText(r.Context(), body.Params, segID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if segID != nil && text != "" {
		_ = h.segment.UpdateSegmentText(r.Context(), *segID, text)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"text": text})
}
