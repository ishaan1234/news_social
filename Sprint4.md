## 1. Issues planned to address in Sprint 4
Sprint 4 focus is **Feature Completion + System Refinement**.

### Planned Sprint 4 items (Feature Completion + System Refinement)
- Frontend: Replacing placeholders in Profile and Settings pages with real data
- Frontend: Improving UI/UX consistency across Home, Posts, Profile, and Settings
- Frontend + Backend Integration: Connecting posts, and news flows end-to-end
- Backend: Completing implementation of all core internal modules (articles, headlines, posts, social, summaries)
- Backend: Optimizing repository layer and database interactions
- Backend: Strengthening error handling and request validation
- Backend: Finalizing API endpoints and ensuring consistency with frontend expectations
- Backend: Enhancing middleware (authentication, logging, rate limiting)
- Testing: Expanding unit tests for backend modules and Cypress tests for frontend workflows
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
    - Posts and interactions
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
    - Posts and social interactions

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
This endpoint registers a new user with an email/password combination. It sends a verification email to the user.

**Request**:
- **Content-Type**: `application/json`
- **Body** (JSON):
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123",
      "display_name": "John Doe"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful, `400` or `409` if error occurs.
- **Body** (JSON):
    ```json
    {
      "success": true,
      "message": "signup successful; verification email sent",
      "id_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "1//0gG7_2oaQeP2S...",
      "expires_in": "3600",
      "user": {
        "uid": "12345",
        "email": "user@example.com",
        "display_name": "John Doe",
        "email_verified": false
      }
    }
    ```

**Errors**:
- `400 Bad Request`: Invalid request body.
- `409 Conflict`: User already exists.

---

### 2. POST `/auth/login`

**Description**:
This endpoint allows users to log in with their email and password. The response will include a Firebase ID token and a refresh token.

**Request**:
- **Content-Type**: `application/json`
- **Body** (JSON):
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful, `400` if error occurs.
- **Body** (JSON):
    ```json
    {
      "success": true,
      "message": "login successful",
      "id_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "1//0gG7_2oaQeP2S...",
      "expires_in": "3600",
      "user": {
        "uid": "12345",
        "email": "user@example.com",
        "display_name": "John Doe",
        "email_verified": true
      }
    }
    ```

**Errors**:
- `400 Bad Request`: Missing email or password.
- `404 Not Found`: User not found.
- `401 Unauthorized`: Invalid credentials.

---

### 3. POST `/auth/verify-email/resend`

**Description**:
This endpoint allows users to request a resend of the verification email if the email was not verified during the signup process.

**Request**:
- **Content-Type**: `application/json`
- **Body** (JSON):
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```

**Response**:
- **Status**: `200 OK` if successful.
- **Body** (JSON):
    ```json
    {
      "success": true,
      "message": "verification email sent",
      "user": {
        "uid": "12345",
        "email": "user@example.com",
        "display_name": "John Doe",
        "email_verified": false
      }
    }
    ```

**Errors**:
- `400 Bad Request`: Missing email or password.
- `404 Not Found`: User not found.
- `401 Unauthorized`: Invalid credentials.

---

### 4. GET `/news`

**Description**:
This endpoint returns the latest news articles from the NewsAPI based on a search query.

**Request**:
- **Query Parameters**:
    - `q` (string, optional): Search query to filter news articles. Default is "tesla".
    ```http
    GET /news?q=technology
    ```

**Response**:
- **Status**: `200 OK` if successful.
- **Body** (JSON):
    ```json
    {
      "status": "ok",
      "totalResults": 100,
      "articles": [
        {
          "source": {
            "id": null,
            "name": "TechCrunch"
          },
          "author": "John Doe",
          "title": "Breaking Tech News",
          "description": "This is a description of the article.",
          "url": "https://techcrunch.com/...",
          "urlToImage": "https://image.url",
          "publishedAt": "2023-04-13T08:00:00Z",
          "content": "Article content...",
          "summary": "Summary of the article..."
        }
      ]
    }
    ```

**Errors**:
- `500 Internal Server Error`: Failed to fetch data from NewsAPI.

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
- 'GOOGLE_APPLICATION_CREDENTIALS': Path to the Firebase service account credentials JSON file.

---
