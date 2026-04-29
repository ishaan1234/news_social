package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestCreatePostHandlerValidation(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantText   string
	}{
		{
			name:       "reject non post",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantText:   "method not allowed",
		},
		{
			name:       "reject invalid json",
			method:     http.MethodPost,
			body:       "{bad json",
			wantStatus: http.StatusBadRequest,
			wantText:   "invalid request body",
		},
		{
			name:       "reject invalid email",
			method:     http.MethodPost,
			body:       `{"user_email":"bad","article_id":"10558852-aea8-459e-8bc1-08ea0eab7714","caption":"hello"}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "invalid email address",
		},
		{
			name:       "reject invalid article id",
			method:     http.MethodPost,
			body:       `{"user_email":"user@example.com","article_id":"bad","caption":"hello"}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "valid article_id is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, _, cleanup := newSQLMockDB(t)
			defer cleanup()

			req := httptest.NewRequest(tc.method, "/posts", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			createPostHandler(db)(rr, req)

			assertStatusAndBody(t, rr, tc.wantStatus, tc.wantText)
		})
	}
}

func TestCreatePostHandlerSuccess(t *testing.T) {
	db, mock, cleanup := newSQLMockDB(t)
	defer cleanup()

	createdAt := time.Date(2026, 4, 29, 16, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs("user@example.com", "10558852-aea8-459e-8bc1-08ea0eab7714", "hello").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_email", "article_id", "caption", "created_at"}).
			AddRow("post-id", "user@example.com", "10558852-aea8-459e-8bc1-08ea0eab7714", "hello", createdAt))

	req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(`{
		"user_email":"user@example.com",
		"article_id":"10558852-aea8-459e-8bc1-08ea0eab7714",
		"caption":"hello"
	}`))
	rr := httptest.NewRecorder()

	createPostHandler(db)(rr, req)

	assertStatusAndBody(t, rr, http.StatusCreated, `"success":true`)
	assertJSONPath(t, rr.Body.Bytes(), "post", "id", "post-id")
}

func TestPostInsertErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"duplicate", &pq.Error{Code: "23505"}, http.StatusConflict, "user has already posted this article"},
		{"foreign key", &pq.Error{Code: "23503"}, http.StatusBadRequest, "user_email or article_id does not exist"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStatus, gotMsg := postInsertError(tc.err)
			if gotStatus != tc.wantStatus || gotMsg != tc.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", gotStatus, gotMsg, tc.wantStatus, tc.wantMsg)
			}
		})
	}
}

func TestFollowingHandlerValidation(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantText   string
	}{
		{
			name:       "reject method",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantText:   "method not allowed",
		},
		{
			name:       "reject invalid follower",
			method:     http.MethodPost,
			body:       `{"follower_email":"bad","following_email":"other@example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "valid follower_email is required",
		},
		{
			name:       "reject invalid following",
			method:     http.MethodPost,
			body:       `{"follower_email":"user@example.com","following_email":"bad"}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "valid following_email is required",
		},
		{
			name:       "reject self follow",
			method:     http.MethodPost,
			body:       `{"follower_email":"user@example.com","following_email":"user@example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "users cannot follow themselves",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, _, cleanup := newSQLMockDB(t)
			defer cleanup()

			req := httptest.NewRequest(tc.method, "/following", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			followingHandler(db)(rr, req)

			assertStatusAndBody(t, rr, tc.wantStatus, tc.wantText)
		})
	}
}

func TestFollowingHandlerFollowAndUnfollow(t *testing.T) {
	db, mock, cleanup := newSQLMockDB(t)
	defer cleanup()

	createdAt := time.Date(2026, 4, 29, 16, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`INSERT INTO following`).
		WithArgs("user@example.com", "other@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"follower_email", "following_email", "created_at"}).
			AddRow("user@example.com", "other@example.com", createdAt))
	mock.ExpectExec(`DELETE FROM following`).
		WithArgs("user@example.com", "other@example.com").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"follower_email":"user@example.com","following_email":"other@example.com"}`

	followReq := httptest.NewRequest(http.MethodPost, "/following", strings.NewReader(body))
	followRR := httptest.NewRecorder()
	followingHandler(db)(followRR, followReq)
	assertStatusAndBody(t, followRR, http.StatusCreated, `"success":true`)

	unfollowReq := httptest.NewRequest(http.MethodDelete, "/following", strings.NewReader(body))
	unfollowRR := httptest.NewRecorder()
	followingHandler(db)(unfollowRR, unfollowReq)
	assertStatusAndBody(t, unfollowRR, http.StatusOK, "unfollowed user")
}

func TestFollowInsertErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"duplicate", &pq.Error{Code: "23505"}, http.StatusConflict, "user is already following this account"},
		{"foreign key", &pq.Error{Code: "23503"}, http.StatusBadRequest, "follower_email or following_email does not exist"},
		{"check", &pq.Error{Code: "23514"}, http.StatusBadRequest, "users cannot follow themselves"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStatus, gotMsg := followInsertError(tc.err)
			if gotStatus != tc.wantStatus || gotMsg != tc.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", gotStatus, gotMsg, tc.wantStatus, tc.wantMsg)
			}
		})
	}
}

func TestPostLikesHandlerLikeAndUnlike(t *testing.T) {
	db, mock, cleanup := newSQLMockDB(t)
	defer cleanup()

	createdAt := time.Date(2026, 4, 29, 16, 0, 0, 0, time.UTC)
	postID := "10558852-aea8-459e-8bc1-08ea0eab7714"
	mock.ExpectQuery(`INSERT INTO post_likes`).
		WithArgs("user@example.com", postID).
		WillReturnRows(sqlmock.NewRows([]string{"user_email", "post_id", "created_at"}).
			AddRow("user@example.com", postID, createdAt))
	mock.ExpectExec(`DELETE FROM post_likes`).
		WithArgs("user@example.com", postID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"user_email":"user@example.com","post_id":"` + postID + `"}`

	likeReq := httptest.NewRequest(http.MethodPost, "/post-likes", strings.NewReader(body))
	likeRR := httptest.NewRecorder()
	postLikesHandler(db)(likeRR, likeReq)
	assertStatusAndBody(t, likeRR, http.StatusCreated, `"success":true`)

	unlikeReq := httptest.NewRequest(http.MethodDelete, "/post-likes", strings.NewReader(body))
	unlikeRR := httptest.NewRecorder()
	postLikesHandler(db)(unlikeRR, unlikeReq)
	assertStatusAndBody(t, unlikeRR, http.StatusOK, "unliked post")
}

func TestPostLikesHandlerValidation(t *testing.T) {
	db, _, cleanup := newSQLMockDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/post-likes", strings.NewReader(`{"user_email":"bad","post_id":"bad"}`))
	rr := httptest.NewRecorder()

	postLikesHandler(db)(rr, req)

	assertStatusAndBody(t, rr, http.StatusBadRequest, "valid user_email is required")
}

func TestPostLikeInsertErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"duplicate", &pq.Error{Code: "23505"}, http.StatusConflict, "user has already liked this post"},
		{"foreign key", &pq.Error{Code: "23503"}, http.StatusBadRequest, "user_email or post_id does not exist"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotStatus, gotMsg := postLikeInsertError(tc.err)
			if gotStatus != tc.wantStatus || gotMsg != tc.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", gotStatus, gotMsg, tc.wantStatus, tc.wantMsg)
			}
		})
	}
}

func TestPostCommentsHandlerCreateAndFetch(t *testing.T) {
	db, mock, cleanup := newSQLMockDB(t)
	defer cleanup()

	createdAt := time.Date(2026, 4, 29, 16, 0, 0, 0, time.UTC)
	postID := "10558852-aea8-459e-8bc1-08ea0eab7714"
	mock.ExpectQuery(`INSERT INTO post_comments`).
		WithArgs(postID, "user@example.com", "hello comment").
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_email", "content", "created_at"}).
			AddRow("comment-id", postID, "user@example.com", "hello comment", createdAt))
	mock.ExpectQuery(`SELECT`).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_email", "username", "display_name", "avatar_url", "content", "created_at"}).
			AddRow("comment-id", postID, "user@example.com", "user", "User", "", "hello comment", createdAt))

	createReq := httptest.NewRequest(http.MethodPost, "/post-comments", strings.NewReader(`{
		"post_id":"`+postID+`",
		"user_email":"user@example.com",
		"content":"hello comment"
	}`))
	createRR := httptest.NewRecorder()
	postCommentsHandler(db)(createRR, createReq)
	assertStatusAndBody(t, createRR, http.StatusCreated, `"success":true`)

	getReq := httptest.NewRequest(http.MethodGet, "/post-comments?post_id="+postID, nil)
	getRR := httptest.NewRecorder()
	postCommentsHandler(db)(getRR, getReq)
	assertStatusAndBody(t, getRR, http.StatusOK, "hello comment")
}

func TestPostCommentsHandlerValidation(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		body       string
		wantStatus int
		wantText   string
	}{
		{
			name:       "reject empty content",
			method:     http.MethodPost,
			target:     "/post-comments",
			body:       `{"post_id":"10558852-aea8-459e-8bc1-08ea0eab7714","user_email":"user@example.com","content":" "}`,
			wantStatus: http.StatusBadRequest,
			wantText:   "content is required",
		},
		{
			name:       "reject bad get post id",
			method:     http.MethodGet,
			target:     "/post-comments?post_id=bad",
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantText:   "valid post_id is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, _, cleanup := newSQLMockDB(t)
			defer cleanup()

			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			postCommentsHandler(db)(rr, req)

			assertStatusAndBody(t, rr, tc.wantStatus, tc.wantText)
		})
	}
}

func TestPostCommentInsertErrorMapping(t *testing.T) {
	gotStatus, gotMsg := postCommentInsertError(&pq.Error{Code: "23503"})
	if gotStatus != http.StatusBadRequest || gotMsg != "post_id or user_email does not exist" {
		t.Fatalf("got (%d, %q), want (%d, %q)", gotStatus, gotMsg, http.StatusBadRequest, "post_id or user_email does not exist")
	}
}

func TestFeedHandlerValidationAndSuccess(t *testing.T) {
	t.Run("reject invalid user email", func(t *testing.T) {
		db, _, cleanup := newSQLMockDB(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/feed?user_email=bad", nil)
		rr := httptest.NewRecorder()

		feedHandler(db)(rr, req)

		assertStatusAndBody(t, rr, http.StatusBadRequest, "valid user_email is required")
	})

	t.Run("success", func(t *testing.T) {
		db, mock, cleanup := newSQLMockDB(t)
		defer cleanup()

		createdAt := time.Date(2026, 4, 29, 16, 0, 0, 0, time.UTC)
		publishedAt := time.Date(2026, 4, 29, 15, 0, 0, 0, time.UTC)
		mock.ExpectQuery(`SELECT`).
			WithArgs("user@example.com").
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_email", "username", "display_name", "avatar_url", "caption", "created_at",
				"like_count", "comment_count", "liked_by_me",
				"article_id", "title", "description", "content", "summary", "author", "source_name", "source_id", "url", "image_url", "published_at", "article_created_at",
			}).AddRow(
				"post-id", "user@example.com", "user", "User", "", "caption", createdAt,
				2, 1, true,
				"article-id", "Title", "Description", "Content", "Summary", "Author", "Source", "source-id", "https://example.com", "", publishedAt, createdAt,
			))

		req := httptest.NewRequest(http.MethodGet, "/feed?user_email=user@example.com", nil)
		rr := httptest.NewRecorder()

		feedHandler(db)(rr, req)

		assertStatusAndBody(t, rr, http.StatusOK, `"success":true`)
		assertStatusAndBody(t, rr, http.StatusOK, `"like_count":2`)
		assertStatusAndBody(t, rr, http.StatusOK, `"comment_count":1`)
		assertStatusAndBody(t, rr, http.StatusOK, `"liked_by_me":true`)
	})
}

func newSQLMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sql mock: %v", err)
	}

	return db, mock, func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}
}

func assertStatusAndBody(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantText string) {
	t.Helper()

	if rr.Code != wantStatus {
		t.Fatalf("got status %d, want %d. body=%s", rr.Code, wantStatus, rr.Body.String())
	}
	if wantText != "" && !strings.Contains(rr.Body.String(), wantText) {
		t.Fatalf("expected body to contain %q, got %q", wantText, rr.Body.String())
	}
}

func assertJSONPath(t *testing.T, body []byte, topKey, nestedKey string, want any) {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	nested, ok := decoded[topKey].(map[string]any)
	if !ok {
		t.Fatalf("expected %q to be an object in %#v", topKey, decoded)
	}
	if got := nested[nestedKey]; got != want {
		t.Fatalf("got %s.%s=%v, want %v", topKey, nestedKey, got, want)
	}
}
