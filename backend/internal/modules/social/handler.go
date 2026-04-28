package social

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ishaan1234/news_social/backend/internal/middleware"
	"github.com/ishaan1234/news_social/backend/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == 0 {
		if userID, ok := r.Context().Value(middleware.UserIDKey).(int); ok {
			req.UserID = userID
		}
	}

	comment, err := h.service.CreateComment(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, comment)
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	headlineID, err := strconv.Atoi(r.URL.Query().Get("headline_id"))
	if err != nil || headlineID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "valid headline_id is required")
		return
	}

	comments, err := h.service.GetComments(r.Context(), headlineID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, comments)
}
