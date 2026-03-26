## 1. Issues planned to address in Sprint 2
Sprint 2 focus was **UI Refining + Backend Orchestration** (placeholders allowed).

### Planned Sprint 2 items (UI Refining + LLM Incorporation and DB Setup)
- Backend: LLM API integration
- Backend: Setting up the db
- Backend: Generating the News Summary
- Backend: Integrating the Database into the workflow
- Integration: The main home page showing real time news and news summary fetched through the API and the backend 
- Frontend: Implementing a social media feed section where people can post their views about a news, and like, comment and share that. (With Placeholders)
- Frontend: Incorporating a chat feature (Frontend with placeholders)
- Testing of Both Frontend and Backend.

## 2. Completed successfully
- **LLM API Integration**
  - Created `frontend/` (React/CRA TS) and `backend/` (Go module) layout
  - Basic run commands established for local dev
- **Setting up the DB**
  - schemas  implemented
  - Migration feature added
- **Generating the News Summary**
- Extended the news fetching pipeline built in Sprint 1.
- Implemented server-side web scraping to extract full article content from URLs using readability parsing.
- Added preprocessing to clean raw HTML content into structured, readable text.
- Integrated Groq LLM API to generate concise summaries from scraped content.
- Enforced a summary length constraint (40 to 60 words).
- Designed pipeline flow: Fetch → Scrape → Clean → Summarize → Respond.
- Secured all API keys (NewsAPI, Groq) using environment variables.
- Enhanced `/news` endpoint to return enriched article data with summaries.
- Improved error handling for external API failures and scraping issues

  
- **Integration of Backend and Frontend for the News Page**
- The home page or the page that displays the news summary now fetches and summarises news in real time instead of placeholders.
  
- **Frontend for social media feed**
- It has placeholder values for selecting the news and posting, with interactive buttons.
- Posts are stored in local browser for now.
- Interactive like, comment and post buttons.
- **Incorporating chat feature frontend**
- It has placeholder values for chats, with interactive buttons.
- Chats are stored in local browser for now.
  
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
- **Home page and navbar**
  - Makes sure the home page components are visible and all the buttons in navbar are clickable
- **Post feed and post interactions**
  - Makes sure 'create a post' section is visible and users can like, comment and the like count gets incremented
- **Chat page**
  - Makes sure different users for chats are clickable and the buttons to send the message work as the chat is displayed
---

## 5. List of Unit Tests for backend
- **Backend Modules**
  - internal\modules\auth\auth_test
  - internal\modules\articles\article_test
  - internal\modules\headlines\headline_test
  - internal\modules\social\social_test
  - internal\modules\summaries\summaries_test
    
- The following components were covered with unit tests:
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
