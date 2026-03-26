## 1. Issues planned to address in Sprint 2
Sprint 2 focus was **UI Refining + Backend Orchestration** (placeholders allowed).

### Planned Sprint 2 items (UI Refining + LLM Incorporation and DB Setup)
- Backend: LLM API integration
- Backend: Setting up the db
- Backend: Integrating the Database into the workflow
- Frontend: Incorporating a chat feature
- Frontend & Backend: User Authentication

## 2. Completed successfully
- **LLM API Integration**
  - Created `frontend/` (React/CRA TS) and `backend/` (Go module) layout
  - Basic run commands established for local dev
- **Setting up the DB**
  - schemas  implemented
  - Migration feature added
- **Incorporating chat feature**
- Extended the news fetching pipeline built in Sprint 1.
- Implemented server-side web scraping to extract full article content from URLs using readability parsing.
- Added preprocessing to clean raw HTML content into structured, readable text.
- Integrated Groq LLM API to generate concise summaries from scraped content.
- Enforced a summary length constraint (40 to 60 words).
- Designed pipeline flow: Fetch → Scrape → Clean → Summarize → Respond.
- Secured all API keys (NewsAPI, Groq) using environment variables.
- Enhanced `/news` endpoint to return enriched article data with summaries.
- Improved error handling for external API failures and scraping issues
  - 

---

## 3. Not completed this sprint (and why)
- **Integrating the Database into the workflow**
  - **Why:** as the schema and migrations are implemented, the current focus is on Article aggregation using the API and LLM after which we will introduce the persistence.
  - **Next step:** Database integration will be incorporated in t to enable caching, and social features.

- **??**
  - **Why:** ??
  - **Next step:** ??

---

## 4. List of Unit Tests & Cypress Tests for frontend
- **? Modules**
  - 
- **?? Modules**
  - 

---

## 5. List of Unit Tests for backend
- **Backend Modules**
  - internal\modules\auth\auth_test
  - internal\modules\articles\article_test
  - internal\modules\headlines\headline_test
  - internal\modules\social\social_test
  - internal\modules\summaries\summaries_test
    
- **The following components were covered with unit tests:
  - `summarizeWithGroq` (LLM integration with mocked API)
  - `extractArticleText` (web scraping logic)
  - API response handling (mocked external calls)
  - Utility functions (`parseDotEnv`, `loadDotEnv`, error handling)
  

---

## 6. Documentation for backend API

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

---
