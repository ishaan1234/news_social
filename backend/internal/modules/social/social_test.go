package social

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

/* Mock Repository */

type mockRepository struct {
	createFunc func(comment *models.Comment) error
}

func (m *mockRepository) Create(comment *models.Comment) error {
	return m.createFunc(comment)
}

/* Service Tests */

func TestService_CreateComment_Success(t *testing.T) {

	mockRepo := &mockRepository{
		createFunc: func(comment *models.Comment) error {
			return nil
		},
	}

	service := NewService(mockRepo)

	req := CreateCommentRequest{
		UserID:     "1",
		HeadlineID: "10",
		Content:    "Great article!",
	}

	result, err := service.CreateComment(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &models.Comment{
		Content: "Great article!",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}

func TestService_CreateComment_Error(t *testing.T) {

	mockRepo := &mockRepository{
		createFunc: func(comment *models.Comment) error {
			return errors.New("db error")
		},
	}

	service := NewService(mockRepo)

	req := CreateCommentRequest{
		Content: "test comment",
	}

	_, err := service.CreateComment(req)

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

/* Handler Tests */

func TestHandler_CreateComment_Success(t *testing.T) {

	mockRepo := &mockRepository{
		createFunc: func(comment *models.Comment) error {
			return nil
		},
	}

	service := NewService(mockRepo)
	handler := NewHandler(service)

	body := CreateCommentRequest{
		UserID:     "1",
		HeadlineID: "10",
		Content:    "Nice article",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateComment(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rr.Code)
	}
}

func TestHandler_CreateComment_BadRequest(t *testing.T) {

	mockRepo := &mockRepository{}

	service := NewService(mockRepo)
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/comments", bytes.NewBuffer([]byte("invalid-json")))
	rr := httptest.NewRecorder()

	handler.CreateComment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 got %d", rr.Code)
	}
}