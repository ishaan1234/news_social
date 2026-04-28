package articles

import (
	"net/http"
	"strconv"

	"github.com/ishaan1234/news_social/backend/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	headlineID, err := strconv.Atoi(r.URL.Query().Get("headline_id"))
	if err != nil || headlineID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "valid headline_id is required")
		return
	}

	articles, err := h.service.GetArticles(r.Context(), headlineID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, articles)
}
