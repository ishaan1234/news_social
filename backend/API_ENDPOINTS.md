# API Endpoints

Base URL for local development:

```text
http://localhost:8080
```

## Auth

### Sign Up

```http
POST /auth/signup
```

```json
{
  "email": "user@example.com",
  "password": "123456",
  "username": "ritik",
  "display_name": "Ritik Raj",
  "avatar_url": "https://example.com/avatar.png"
}
```

### Login

```http
POST /auth/login
```

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

### Resend Verification Email

```http
POST /auth/verify-email/resend
```

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

## News / Articles

### Fetch News

```http
GET /news?q=tesla
```

No JSON body.

This fetches news, summarizes articles, and saves them to the `articles` table.

## Posts

### Create Post

```http
POST /posts
```

```json
{
  "user_email": "user@example.com",
  "article_id": "article-uuid",
  "caption": "My thoughts on this article"
}
```

`article_id` must come from the `articles` table.

## Feed

### Get Feed

```http
GET /feed?user_email=user@example.com
```

No JSON body.

Returns the user's own posts plus posts from people they follow. Each feed item includes article details, `like_count`, `comment_count`, and `liked_by_me`.

## Following

### Follow User

```http
POST /following
```

```json
{
  "follower_email": "user@example.com",
  "following_email": "other@example.com"
}
```

### Unfollow User

```http
DELETE /following
```

```json
{
  "follower_email": "user@example.com",
  "following_email": "other@example.com"
}
```

## Likes

### Like Post

```http
POST /post-likes
```

```json
{
  "user_email": "user@example.com",
  "post_id": "post-uuid"
}
```

### Unlike Post

```http
DELETE /post-likes
```

```json
{
  "user_email": "user@example.com",
  "post_id": "post-uuid"
}
```

`post_id` must come from the `posts` table.

## Comments

### Create Comment

```http
POST /post-comments
```

```json
{
  "post_id": "post-uuid",
  "user_email": "user@example.com",
  "content": "This is my comment"
}
```

### Get Comments For Post

```http
GET /post-comments?post_id=post-uuid
```

No JSON body.

Returns comments for a post with user details.
