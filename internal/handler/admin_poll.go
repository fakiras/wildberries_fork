package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GeneratePoll — POST /admin/promotions/{id}/poll/generate
func (h *AdminHandler) GeneratePoll(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Type string `json:"type"` // "questions" | "answer_tree"
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	// TODO: call AI, then save via PollRepository
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"questions": nil, "answerTree": nil})
}

// SetPollQuestions — POST /admin/promotions/{id}/poll/questions
func (h *AdminHandler) SetPollQuestions(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Questions []struct {
			Text    string `json:"text"`
			Options []struct {
				Text  string `json:"text"`
				Value string `json:"value"`
			} `json:"options"`
		} `json:"questions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	// TODO: call PollRepository.SaveQuestions
	_ = body
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// SetAnswerTree — POST /admin/promotions/{id}/poll/answer-tree
func (h *AdminHandler) SetAnswerTree(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	var body struct {
		Nodes []struct {
			NodeID      string `json:"node_id"`
			ParentNodeID string `json:"parent_node_id"`
			Label       string `json:"label"`
			Value       string `json:"value"`
		} `json:"nodes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	// TODO: call PollRepository.SaveAnswerTree
	_ = body
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
