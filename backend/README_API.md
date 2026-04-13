# API Documentation

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