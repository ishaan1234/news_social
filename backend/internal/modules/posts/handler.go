package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ishaan1234/news_social/backend/internal/models"
	"github.com/ishaan1234/news_social/backend/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPosts(w, r)
	case http.MethodPost:
		h.createPost(w, r)
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getComments(w, r)
	case http.MethodPost:
		h.createComment(w, r)
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		PostID  int    `json:"post_id"`
		VoterID string `json:"voter_id"`
		Value   int    `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	summary, err := h.service.SetVote(r.Context(), req.PostID, req.VoterID, req.Value)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, summary)
}

func (h *Handler) Share(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		PostID int `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	shareCount, err := h.service.SharePost(r.Context(), req.PostID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]int{
		"post_id":     req.PostID,
		"share_count": shareCount,
	})
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.CreatePost(r.Context(), post)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetPosts(r.Context(), r.URL.Query().Get("viewer_id"))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	var comment models.PostComment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.CreateComment(r.Context(), comment)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) getComments(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || postID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "valid post_id is required")
		return
	}

	comments, err := h.service.GetComments(r.Context(), postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, comments)
}
