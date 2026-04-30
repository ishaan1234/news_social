## 1. Issues planned to address in Sprint 4
Sprint 4 focus is **Feature Completion + System Refinement**.

### Planned Sprint 4 items (Feature Completion + System Refinement)

- Frontend: Replacing placeholders in Profile and Settings pages with real data
- Frontend: Improving UI/UX consistency across Home, Posts, Profile, and Settings
- Frontend + Backend Integration: Connecting posts, feed, social interactions, and news flows end-to-end
- Backend: Completing implementation of all core internal modules, including articles, headlines, posts, social, summaries, and authentication
- Backend: Implementing social API endpoints for posts, feed, following/unfollowing, likes/unlikes, and comments
- Backend + Database: Integrating the backend with Supabase PostgreSQL for persistent cloud storage
- Backend + Database: Connecting posts, articles, follows, likes, and comments to database-backed repositories
- Backend: Optimizing repository layer and database interactions
- Backend: Strengthening error handling and request validation
- Backend: Finalizing API endpoints and ensuring consistency with frontend expectations
- Backend: Enhancing middleware, including authentication, logging, and rate limiting
- Testing: Expanding unit tests for backend modules, including `social_handlers_test.go`, and Cypress tests for frontend workflows
- Documentation: Finalizing API documentation and overall project README

---

## 2. Completed successfully
- **Frontend Enhancements**
  - Improved UI consistency across all major pages (Home, Posts, Chat, Profile, Settings)
  - Integrated frontend pages with backend APIs (Posts, Profile, Home)
  - Reduced reliance on placeholders by incorporating real backend data

- **Backend (Internal Modules Completion)**
  - Completed implementation of:
    - Articles module
    - Headlines module (including aggregation logic)
    - Posts module (with both in-memory and persistent repositories)
    - Social module (follow/interactions)
    - Summaries module (AI-based summarization pipeline)
  - Maintained clean architecture using handler → service → repository pattern

- **Backend API Stabilization**
  - Finalized endpoints for:
    - Authentication
    - News retrieval and summarization
    - Post creation
    - Feed retrieval
    - Following and unfollowing users
    - Post likes and unlikes
    - Post comments and comment retrieval
  - Standardized API responses using utility functions

- **Database & Migrations**
  - Completed schema design and verified migrations for:
    - Users, Articles, Summaries, Posts, Comments, Votes
  - Ensured database consistency and integration with repository layer

- **Middleware Improvements**
  - Enhanced:
    - Authentication middleware
    - Logging middleware
    - Rate limiting
  - Improved request tracing and debugging

- **End-to-End Integration**
  - Fully integrated frontend ↔ API layer ↔ backend/internal modules
  - Verified workflows:
    - User authentication
    - News retrieval with summaries
    - Creating posts from saved articles
    - Feed generation for a user and followed users
    - Following/unfollowing users
    - Liking/unliking posts
    - Creating and viewing comments

---

## 3. List of Unit Tests & Cypress Tests for frontend
- Cypress component tests:
  - home.cy.jsx
  - posts.cy.jsx
  - profile.cy.jsx
  - settings.cy.jsx

- Functional tests:
  - Navigation across pages
  - Button interactions in Profile and Settings
  - Form validation and user input handling

---

## 4. List of Unit Tests for backend
- **API Layer Tests**
  - cmd/api/backend_test
  - cmd/api/firebase_email_password_auth_test
    - TestNormalizeEmail
    - TestNewFirebaseIdentityClient
    - TestSignupHandlerValidation
    - TestLoginHandlerValidation
    - TestResendVerificationEmailHandlerValidation
    - TestRegisterFirebaseEmailPasswordRoutes
    - TestValidationResponsesAreJSON
  - cmd/api/social_handlers_test.go
    - Covers create post handler validation and response handling
    - Covers feed handler query validation and feed response behavior
    - Covers follow and unfollow handler validation
    - Covers post like and unlike handler validation
    - Covers create comment and get comments handler behavior
- **Backend Modules**
  - internal/server/server_http_test
  - internal/utils/errors_test
  - internal/utils/http_response_test

- **Module Tests**
  - internal/modules/articles/article_test
  - internal/modules/headlines/headline_test
  - internal/modules/posts/posts_test
  - internal/modules/social/social_test
  - internal/modules/summaries/summaries_test
  - internal/modules/auth/auth_test

---

## 5. Documentation for backend API

### Endpoint: `GET /news`

This endpoint is used to fetch news articles for a given search query and return them along with AI-generated summaries.

#### Purpose
The goal of this endpoint is to go beyond simple news retrieval. After fetching article metadata, the backend scrapes the main readable content from each article, processes that text, and sends it to the LLM to generate a short summary. The final response returns both the original article information and the generated summary.

#### Request Method
`GET`

#### Route
`/news`

#### Query Parameters

| Parameter | Type   | Required | Description |
|----------|--------|----------|-------------|
| `q`      | string | Yes      | Search keyword used to fetch related news articles |

#### Example Request
```http
GET /news?q=tesla
GET /news?q=apple
```

## Base URL

The base URL for this API is configurable through environment variables but is assumed to be the default server address in this documentation.
```plaintext
http://localhost:8080
```

## Authentication

### Firebase Authentication

The API relies on Firebase Authentication for user registration and login. All authenticated routes require a valid Firebase ID token, which should be sent in the `Authorization` header as a Bearer token.

---

## Endpoints

### 1. POST `/auth/signup`

**Description**:
Registers a new user with an email/password combination and sends a verification email.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123",
      "username": "ritik",
      "display_name": "Ritik Raj",
      "avatar_url": "https://example.com/avatar.png"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful; `400 Bad Request` or `409 Conflict` if an error occurs.

---

### 2. POST `/auth/login`

**Description**:
Logs in a user with email and password. The response includes Firebase authentication tokens.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful.
- Returns the ID token, refresh token, expiration time, and user details.

**Errors**:
- `400 Bad Request`: Missing or invalid request data.
- `401 Unauthorized`: Invalid credentials.
- `404 Not Found`: User not found.

---

### 3. POST `/auth/verify-email/resend`

**Description**:
Resends the Firebase verification email for a user who has not verified their email yet.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful.
- Returns a success message and user verification status.

**Errors**:
- `400 Bad Request`: Missing email or password.
- `401 Unauthorized`: Invalid credentials.
- `404 Not Found`: User not found.

---

### 4. GET `/news`

**Description**:
Fetches news articles for a search query, generates AI summaries, and saves the article records to the `articles` table.

**Request**:
- **Query Parameters**:

| Parameter | Type   | Required | Description |
|----------|--------|----------|-------------|
| `q`      | string | No       | Search keyword used to fetch related news articles. Defaults to `tesla` if omitted. |

**Example Request**:
```http
GET /news?q=tesla
GET /news?q=apple
```

**Response**:
- **Status**: `200 OK` if successful.
- Returns article metadata, source information, article content, and generated summaries.

**Errors**:
- `500 Internal Server Error`: Failed to fetch data from NewsAPI or summarize article content.

---

### 5. POST `/posts`

**Description**:
Creates a post for a user using an existing article from the `articles` table.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "user_email": "user@example.com",
      "article_id": "article-uuid",
      "caption": "My thoughts on this article"
    }
    ```

**Notes**:
- `article_id` must reference an existing article in the `articles` table.

**Response**:
- **Status**: `200 OK` or `201 Created` if successful, depending on handler implementation.
- Returns the created post record.

**Errors**:
- `400 Bad Request`: Missing user email, article ID, or invalid request body.
- `500 Internal Server Error`: Failed to create the post.

---

### 6. GET `/feed`

**Description**:
Returns the personalized feed for a user. The feed includes the user's own posts plus posts from users they follow.

**Request**:
- **Query Parameters**:

| Parameter    | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `user_email` | string | Yes      | Email of the user requesting the feed. |

**Example Request**:
```http
GET /feed?user_email=user@example.com
```

**Response**:
- **Status**: `200 OK` if successful.
- Returns feed items with article details, `like_count`, `comment_count`, and `liked_by_me`.

**Errors**:
- `400 Bad Request`: Missing `user_email`.
- `500 Internal Server Error`: Failed to load feed data.

---

### 7. POST `/following`

**Description**:
Allows one user to follow another user.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "follower_email": "user@example.com",
      "following_email": "other@example.com"
    }
    ```

**Response**:
- **Status**: `200 OK` or `201 Created` if successful.
- Returns a success message or follow relationship data.

**Errors**:
- `400 Bad Request`: Missing follower or following email, or user attempted to follow themselves.
- `500 Internal Server Error`: Failed to create follow relationship.

---

### 8. DELETE `/following`

**Description**:
Allows one user to unfollow another user.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "follower_email": "user@example.com",
      "following_email": "other@example.com"
    }
    ```

**Response**:
- **Status**: `200 OK` or `204 No Content` if successful.

**Errors**:
- `400 Bad Request`: Missing follower or following email.
- `500 Internal Server Error`: Failed to remove follow relationship.

---

### 9. POST `/post-likes`

**Description**:
Adds a like to a post for a user.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "user_email": "user@example.com",
      "post_id": "post-uuid"
    }
    ```

**Notes**:
- `post_id` must reference an existing post in the `posts` table.

**Response**:
- **Status**: `200 OK` or `201 Created` if successful.
- Returns a success message or updated like state.

**Errors**:
- `400 Bad Request`: Missing user email or post ID.
- `500 Internal Server Error`: Failed to like the post.

---

### 10. DELETE `/post-likes`

**Description**:
Removes a user's like from a post.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "user_email": "user@example.com",
      "post_id": "post-uuid"
    }
    ```

**Response**:
- **Status**: `200 OK` or `204 No Content` if successful.

**Errors**:
- `400 Bad Request`: Missing user email or post ID.
- `500 Internal Server Error`: Failed to unlike the post.

---

### 11. POST `/post-comments`

**Description**:
Creates a comment on a post.

**Request**:
- **Content-Type**: `application/json`
- **Body**:
    ```json
    {
      "post_id": "post-uuid",
      "user_email": "user@example.com",
      "content": "This is my comment"
    }
    ```

**Response**:
- **Status**: `200 OK` or `201 Created` if successful.
- Returns the created comment.

**Errors**:
- `400 Bad Request`: Missing post ID, user email, or comment content.
- `500 Internal Server Error`: Failed to create the comment.

---

### 12. GET `/post-comments`

**Description**:
Returns all comments for a given post, including user details.

**Request**:
- **Query Parameters**:

| Parameter | Type   | Required | Description |
|----------|--------|----------|-------------|
| `post_id` | string | Yes      | ID of the post whose comments should be returned. |

**Example Request**:
```http
GET /post-comments?post_id=post-uuid
```

**Response**:
- **Status**: `200 OK` if successful.
- Returns the list of comments for the post with user details.

**Errors**:
- `400 Bad Request`: Missing `post_id`.
- `500 Internal Server Error`: Failed to load comments.

---

## Utility Functions

### `extractArticleText`

**Description**:
Extracts the main readable content of a news article from a provided URL using the `go-readability` library, stripping out ads, sidebars, and other non-essential content.

**Parameters**:
- `articleURL` (string): The URL of the article.

**Returns**:
- A string containing the plain text content of the article or an error if extraction fails.

---

### `summarizeWithGroq`

**Description**:
Sends a news article's content to Groq's API for summarization.

**Parameters**:
- `content` (string): The content of the news article.

**Returns**:
- A string containing the summarized version of the article or an error if summarization fails.

---

## Environment Variables

- `NEWSAPI_KEY`: API key for NewsAPI.
- `FIREBASE_WEB_API_KEY`: Firebase Web API key.
- `GROQ_API_KEY`: API key for Groq API.
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to the Firebase service account credentials JSON file.

---
