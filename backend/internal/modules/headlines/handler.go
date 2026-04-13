package headlines

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateHeadline(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.CreateHeadline(r.Context(), req.Title)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetHeadlines(w http.ResponseWriter, r *http.Request) {
	data, _ := h.service.GetHeadlines(r.Context())
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetHeadline(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	data, _ := h.service.GetHeadline(r.Context(), id)
	json.NewEncoder(w).Encode(data)
}