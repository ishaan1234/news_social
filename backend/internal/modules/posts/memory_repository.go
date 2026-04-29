package posts

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type MemoryRepository struct {
	mu            sync.Mutex
	nextPostID    int
	nextCommentID int
	posts         []models.Post
	comments      map[int][]models.PostComment
	votes         map[int]map[string]int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextPostID:    1,
		nextCommentID: 1,
		comments:      make(map[int][]models.PostComment),
		votes:         make(map[int]map[string]int),
	}
}

func (r *MemoryRepository) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post.ID = r.nextPostID
	r.nextPostID++
	post.CreatedAt = time.Now().UTC()
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *MemoryRepository) GetPosts(ctx context.Context, viewerID string) ([]models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]models.Post, 0, len(r.posts))
	for _, post := range r.posts {
		result = append(result, r.withCounts(post, viewerID))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (r *MemoryRepository) GetPostByID(ctx context.Context, postID int, viewerID string) (models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, post := range r.posts {
		if post.ID == postID {
			return r.withCounts(post, viewerID), nil
		}
	}

	return models.Post{}, sql.ErrNoRows
}

func (r *MemoryRepository) CreateComment(ctx context.Context, comment models.PostComment) (models.PostComment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.hasPost(comment.PostID) {
		return models.PostComment{}, fmt.Errorf("post not found")
	}

	comment.ID = r.nextCommentID
	r.nextCommentID++
	comment.CreatedAt = time.Now().UTC()
	r.comments[comment.PostID] = append(r.comments[comment.PostID], comment)
	return comment, nil
}

func (r *MemoryRepository) GetComments(ctx context.Context, postID int) ([]models.PostComment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	comments := append([]models.PostComment(nil), r.comments[postID]...)
	return comments, nil
}

func (r *MemoryRepository) SetVote(ctx context.Context, postID int, voterID string, value int) (models.PostVoteSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.hasPost(postID) {
		return models.PostVoteSummary{}, fmt.Errorf("post not found")
	}

	if _, exists := r.votes[postID]; !exists {
		r.votes[postID] = make(map[string]int)
	}

	if value == 0 {
		delete(r.votes[postID], voterID)
	} else {
		r.votes[postID][voterID] = value
	}

	return r.voteSummary(postID, voterID), nil
}

func (r *MemoryRepository) IncrementShare(ctx context.Context, postID int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.posts {
		if r.posts[i].ID == postID {
			r.posts[i].ShareCount++
			return r.posts[i].ShareCount, nil
		}
	}

	return 0, fmt.Errorf("post not found")
}

func (r *MemoryRepository) hasPost(postID int) bool {
	for _, post := range r.posts {
		if post.ID == postID {
			return true
		}
	}
	return false
}

func (r *MemoryRepository) withCounts(post models.Post, viewerID string) models.Post {
	summary := r.voteSummary(post.ID, viewerID)
	post.VoteScore = summary.VoteScore
	post.ViewerVote = summary.ViewerVote
	post.CommentCount = len(r.comments[post.ID])
	return post
}

func (r *MemoryRepository) voteSummary(postID int, viewerID string) models.PostVoteSummary {
	summary := models.PostVoteSummary{PostID: postID}
	for voterID, value := range r.votes[postID] {
		summary.VoteScore += value
		if voterID == viewerID {
			summary.ViewerVote = value
		}
	}
	return summary
}
