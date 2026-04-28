package headlines

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ishaan1234/news_social/backend/internal/utils"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.service.CreateHeadline(r.Context(), req.Title)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) GetHeadlines(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetHeadlines(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) GetHeadline(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "valid id is required")
		return
	}

	data, err := h.service.GetHeadline(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}
