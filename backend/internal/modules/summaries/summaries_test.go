package summaries

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

/* Mock AI Client */

type mockAIClient struct {
	summarizeFunc func(text string) (string, error)
}

func (m *mockAIClient) Summarize(text string) (string, error) {
	return m.summarizeFunc(text)
}

/* Mock Repository */

type mockRepository struct {
	saveFunc func(summary *models.Summary) error
	findFunc func(headlineID int64) (*models.Summary, error)
}

func (m *mockRepository) Save(summary *models.Summary) error {
	return m.saveFunc(summary)
}

func (m *mockRepository) FindByHeadline(headlineID int64) (*models.Summary, error) {
	if m.findFunc != nil {
		return m.findFunc(headlineID)
	}
	return nil, nil
}

/* Service Tests */

func TestService_GenerateSummary_Success(t *testing.T) {

	expectedContent := "AI generated summary"

	mockAI := &mockAIClient{
		summarizeFunc: func(text string) (string, error) {
			return expectedContent, nil
		},
	}

	mockRepo := &mockRepository{
		saveFunc: func(summary *models.Summary) error {
			return nil
		},
	}

	service := NewService(mockRepo, mockAI)

	result, err := service.GenerateSummary("1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &models.Summary{
		Content: expectedContent,
		Model:   "gpt-4",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}

func TestService_GenerateSummary_AIError(t *testing.T) {

	mockAI := &mockAIClient{
		summarizeFunc: func(text string) (string, error) {
			return "", errors.New("AI failure")
		},
	}

	mockRepo := &mockRepository{}

	service := NewService(mockRepo, mockAI)

	_, err := service.GenerateSummary("1")

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

/* Handler Tests */

func TestHandler_Generate(t *testing.T) {

	mockAI := &mockAIClient{
		summarizeFunc: func(text string) (string, error) {
			return "AI summary", nil
		},
	}

	mockRepo := &mockRepository{
		saveFunc: func(summary *models.Summary) error {
			return nil
		},
	}

	service := NewService(mockRepo, mockAI)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/generate?headline_id=1", nil)
	rr := httptest.NewRecorder()

	handler.Generate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rr.Code)
	}
}

func TestHandler_Generate_Error(t *testing.T) {

	mockAI := &mockAIClient{
		summarizeFunc: func(text string) (string, error) {
			return "", errors.New("AI error")
		},
	}

	mockRepo := &mockRepository{}

	service := NewService(mockRepo, mockAI)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/generate?headline_id=1", nil)
	rr := httptest.NewRecorder()

	handler.Generate(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 got %d", rr.Code)
	}
}