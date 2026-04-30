# News Social

Our app is a centralized platform designed to bridge the gap between information consumption and public opinion. It eliminates the friction of context-switching between third-party news sources and social commentary platforms by bringing the article and the conversation into a single, cohesive interface.

In the current digital landscape, users often have to link external articles or blog posts to provide context for their opinions. Our app solves this by hosting the news content natively, allowing users to read, analyze, and share their perspectives without leaving the ecosystem. It is designed for users who want to stay informed and engage in meaningful, context-rich discussions in real-time.

**News Social** is a full-stack platform that combines **news consumption and social interaction** into a single, seamless experience. Users can read news articles, engage in discussions, and share opinions—all within one unified ecosystem.

---

## Overview

In today’s digital landscape, users often switch between multiple platforms to read news and participate in discussions. News Social eliminates this friction by:

- Hosting news content natively
- Enabling real-time discussions tied to articles
- Providing a centralized platform for informed conversations

The goal is to create a **context-rich, interactive news experience**.

---

## Team Members

| Name | Role | Email |
|------|------|-------|
| Ishaan Gupta | Frontend | guptaishaan@ufl.edu |
| Ritik Raj | Backend | ritikraj.lnu@ufl.edu |
| Vittal Chintamaneni | Backend | chintamaneni.v@ufl.edu |

---

## Project Architecture
```
news_social/
│
├── .env
├── .gitignore
├── README.md
├── Sprint1.md
├── Sprint2.md
├── Sprint3.md
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── backend/
│   ├── .gitignore
│   ├── API_ENDPOINTS.md
│   ├── README_API.md
│   ├── go.mod
│   ├── go.sum
│   │
│   ├── cmd/
│   │   └── api/
│   │       ├── main.go
│   │       ├── ai.go
│   │       ├── feed.go
│   │       ├── news.go
│   │       ├── posts.go
│   │       ├── post_comments.go
│   │       ├── post_likes.go
│   │       ├── following.go
│   │       ├── scrape.go
│   │       ├── firebase_auth.go
│   │       ├── firebase_email_password_auth.go
│   │       ├── utils.go
│   │       ├── backend_test.go
│   │       ├── firebase_email_password_auth_test.go
│   │       └── social_handlers_test.go
│   │
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   │
│   │   ├── db/
│   │   │   ├── postgres.go
│   │   │   ├── migrate.go
│   │   │   └── migrations/
│   │   │       ├── 001_users.sql
│   │   │       ├── 002_headlines.sql
│   │   │       ├── 003_articles.sql
│   │   │       ├── 004_summaries.sql
│   │   │       ├── 005_comments.sql
│   │   │       ├── 006_votes.sql
│   │   │       └── 007_posts.sql
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── rate_limit.go
│   │   │
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── article.go
│   │   │   ├── headline.go
│   │   │   ├── summary.go
│   │   │   ├── comment.go
│   │   │   ├── post.go
│   │   │   └── vote.go
│   │   │
│   │   ├── modules/
│   │   │   ├── articles/
│   │   │   │   ├── handler.go
│   │   │   │   ├── service.go
│   │   │   │   ├── repository.go
│   │   │   │   └── article_test.go
│   │   │   │
│   │   │   ├── auth/
│   │   │   │   ├── service.go
│   │   │   │   └── auth_test.go
│   │   │   │
│   │   │   ├── headlines/
│   │   │   │   ├── handler.go
│   │   │   │   ├── service.go
│   │   │   │   ├── repository.go
│   │   │   │   ├── aggregate.go
│   │   │   │   └── headline_test.go
│   │   │   │
│   │   │   ├── posts/
│   │   │   │   ├── handler.go
│   │   │   │   ├── service.go
│   │   │   │   ├── repository.go
│   │   │   │   ├── memory_repository.go
│   │   │   │   └── posts_test.go
│   │   │   │
│   │   │   ├── social/
│   │   │   │   ├── handler.go
│   │   │   │   ├── service.go
│   │   │   │   ├── repository.go
│   │   │   │   └── social_test.go
│   │   │   │
│   │   │   └── summaries/
│   │   │       ├── handler.go
│   │   │       ├── service.go
│   │   │       ├── repository.go
│   │   │       └── summaries_test.go
│   │   │
│   │   ├── server/
│   │   │   ├── server_http.go
│   │   │   └── server_http_test.go
│   │   │
│   │   └── utils/
│   │       ├── errors.go
│   │       ├── errors_test.go
│   │       ├── http_response.go
│   │       └── http_response_test.go
│   │
│   └── pkg/
│       └── clients/
│           ├── ai/
│           │   └── openai.go
│           │
│           └── newsapi/
│               └── client.go
│
└── frontend/
    ├── package.json
    ├── package-lock.json
    ├── tsconfig.json
    ├── tailwind.config.js
    ├── postcss.config.js
    ├── cypress.config.js
    │
    ├── cypress/
    │   ├── component/
    │   │   ├── chat.cy.jsx
    │   │   ├── home.cy.jsx
    │   │   ├── posts.cy.jsx
    │   │   ├── profile.cy.jsx
    │   │   └── settings.cy.jsx
    │   │
    │   └── support/
    │       ├── component.js
    │       └── component-index.html
    │
    ├── public/
    │   ├── index.html
    │   ├── favicon.ico
    │   ├── logo192.png
    │   ├── logo512.png
    │   ├── manifest.json
    │   └── robots.txt
    │
    └── src/
        ├── App.tsx
        ├── App.css
        ├── App.test.tsx
        ├── index.tsx
        ├── index.css
        ├── setupTests.ts
        ├── setupProxy.js
        ├── reportWebVitals.ts
        ├── react-app-env.d.ts
        ├── auth.ts
        ├── postArticleDraft.ts
        ├── Posts-Cypress-Test-Summary.md
        ├── logo.svg
        │
        ├── components/
        │   ├── Navbar.tsx
        │   └── NewsCard.tsx
        │
        └── pages/
            ├── Home.tsx
            ├── Auth.tsx
            ├── Chat.tsx
            ├── Posts.tsx
            ├── Profile.tsx
            ├── Settings.tsx
            └── PlaceholderPage.tsx
```

---

## Tech Stack

### Frontend
- React.js
- JavaScript / TypeScript
- CSS / Tailwind (if used)

### Backend
- Go (Golang)
- RESTful API architecture

### Authentication
- Firebase Authentication (JWT-based)

### Database
- Supabase PostgreSQL cloud database

---

## Authentication

The system uses **Firebase Authentication**:

- Users sign up/login via email & password
- Backend validates requests using **Firebase ID Tokens**
- All protected routes require:
    Authorization: Bearer <ID_TOKEN>

---

## API Overview

Base URL: [LocalHost](http://localhost:8080)


### Key Endpoints

#### Authentication
- `POST /auth/signup` → Register a new user and send verification email
- `POST /auth/login` → Log in an existing user
- `POST /auth/verify-email/resend` → Resend email verification link

#### News / Articles
- `GET /news?q=tesla` → Fetch news articles, generate summaries, and save articles to the database

#### Posts
- `POST /posts` → Create a post from an existing saved article

#### Feed
- `GET /feed?user_email=user@example.com` → Retrieve the user’s own posts plus posts from followed users

#### Following
- `POST /following` → Follow another user
- `DELETE /following` → Unfollow another user

#### Likes
- `POST /post-likes` → Like a post
- `DELETE /post-likes` → Unlike a post

#### Comments
- `POST /post-comments` → Add a comment to a post
- `GET /post-comments?post_id=post-uuid` → Retrieve comments for a post

Full API details are available in: `backend/README_API.md`


---

## Setup Instructions

#### 1. Clone the Repository

```bash
git clone <repo-url>
cd news_social
```
#### 2. Backend Setup
```bash
cd backend
go mod tidy
go run main.go
```

#### 3. Setting up environment variables
```bash
PORT=8080
NEWSAPI_KEY=your_newsapi_key
GROQ_API_KEY=your_groq_api_key
FIREBASE_WEB_API_KEY=your_firebase_web_api_key
GOOGLE_APPLICATION_CREDENTIALS=path/to/firebase-service-account.json
DATABASE_URL=your_supabase_postgresql_connection_string
```

#### 4. Frontend Setup
```bash
cd frontend
npm install
npm start
```

### Running the tests

#### Backend:
```bash
cd backend
go test ./...
```
#### Frontend:
```bash
npm test
```
