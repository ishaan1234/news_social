package summaries

import (
	"encoding/json"
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

type GenerateSummaryRequest struct {
	HeadlineID int    `json:"headline_id"`
	Content    string `json:"content"`
}

func (h *Handler) GenerateSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var req GenerateSummaryRequest
	if r.Method == http.MethodGet {
		id, err := strconv.Atoi(r.URL.Query().Get("headline_id"))
		if err != nil || id <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "valid headline_id is required")
			return
		}
		req.HeadlineID = id
		req.Content = r.URL.Query().Get("content")
	} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	summary, err := h.service.GenerateAndSaveSummary(r.Context(), req.HeadlineID, req.Content)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, summary)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	headlineID, err := strconv.Atoi(r.URL.Query().Get("headline_id"))
	if err != nil || headlineID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "valid headline_id is required")
		return
	}

	summary, err := h.service.GetSummary(r.Context(), headlineID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, summary)
}
