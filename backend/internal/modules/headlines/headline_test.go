package headlines

import (
	"errors"
    "testing"
	/*"net/http"
	"net/http/httptest"
    "reflect"
    "context"*/
	

	"github.com/go-chi/chi/v5"
	"github.com/ishaan1234/news_social/backend/internal/models"
    /*"github.com/google/uuid"
    headlineService "github.com/ishaan1234/news_social/backend/internal/service"*/
)

/* Mock for ArticleAggregator */
type mockAggregator struct {
	listHeadlinesFunc       func() ([]models.Headline, error)
	fetchHeadlineBundleFunc func(id string) (any, error)
}

func (m *mockAggregator) ListHeadlines() ([]models.Headline, error) {
	return m.listHeadlinesFunc()
}

func (m *mockAggregator) FetchHeadlineBundle(id string) (any, error) {
	return m.fetchHeadlineBundleFunc(id)
}

/* Tests for original Service (ArticleAggregator) and Handler */
/*
func TestService_ListHeadlines(t *testing.T) {

	expected := []models.Headline{
		{ID: "1", Title: "Breaking News"},
		{ID: "2", Title: "Tech Headlines"},
	}

	mockRepo := &mockAggregator{
		listHeadlinesFunc: func() ([]models.Headline, error) {
			return expected, nil
		},
	}

	service := NewService(mockRepo)

	result, err := service.ListHeadlines()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v got %v", expected, result)
	}
}
*/

func TestService_GetHeadlineDetails_Error(t *testing.T) {

	mockRepo := &mockAggregator{
		fetchHeadlineBundleFunc: func(id string) (any, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewService(mockRepo)

	_, err := service.GetHeadlineDetails("123")

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

/* Handler Tests */
/*
func TestHandler_List(t *testing.T) {

	expected := []models.Headline{
		{ID: "1", Title: "World News"},
	}

	mockRepo := &mockAggregator{
		listHeadlinesFunc: func() ([]models.Headline, error) {
			return expected, nil
		},
	}

	service := NewService(mockRepo)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200 got %d", status)
	}
}

func TestHandler_Get(t *testing.T) {

	mockRepo := &mockAggregator{
		fetchHeadlineBundleFunc: func(id string) (any, error) {
			return map[string]string{"headline": "Breaking"}, nil
		},
	}

	service := NewService(mockRepo)
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(contextWithChi(rctx))

	handler.Get(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200 got %d", status)
	}
}
*/
/* headline_service Tests (new code) */
/*
func TestHeadlineService_CreateFromArticles(t *testing.T) {
	hService := &headlineService.HeadlineService{}

	articles := []models.Article{
		{Title: "Test Article"},
	}

	result, err := hService.CreateFromArticles(context.Background(), articles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Title != articles[0].Title {
		t.Errorf("expected title %s got %s", articles[0].Title, result.Title)
	}
}

func TestHeadlineService_GetByID(t *testing.T) {
	hService := &headlineService.HeadlineService{}

	result, err := hService.GetByID(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result got %v", result)
	}
}

func TestHeadlineService_GetTrending(t *testing.T) {
	hService := &headlineService.HeadlineService{}

	result, err := hService.GetTrending(context.Background(), 5)
	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty slice got %v", result)
	}
}
*/

/* Helper for chi context */
func contextWithChi(rctx *chi.Context) interface{} {
	return rctx
}