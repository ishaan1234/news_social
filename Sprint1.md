# User Stories — news_social

## US1 — View categorized news feed
**As a user**, I want to see a feed of the latest news across different categories, so I can quickly understand what is happening in areas that interest me.  
**Acceptance Criteria**
- User can view latest headlines on the homepage
- Articles are grouped or filterable by category (e.g., technology, business, sports)
- Each article displays title, source, date, and summary

---

## US2 — Auto 50-word summary
**As a user**, I want each news article to display a concise 50-word summary by default, so I can quickly scan and understand the main idea without opening the full article.  
**Acceptance Criteria**
- Each article displays a summary of ≤ 50 words
- Summary appears directly below the title
- If summarization fails, a fallback description is shown

---

## US3 — Open full article
**As a user**, when I find an article interesting, I want to open the original source, so I can read the complete story.  
**Acceptance Criteria**
- Clicking the article title or publisher name opens the original article URL
- Article opens in a new browser tab

---

## US4 — Interact with an article (Like & Save)
**As a user**, I want to like and save an article, so I can express interest and easily access it later.  
**Acceptance Criteria**
- User can click a Like button on an article
- Like count updates immediately
- User can unlike an article
- User can click a Save (Bookmark) button
- Saved articles appear in a “Saved Articles” section
- User can remove a saved article
- Likes and saves persist per logged-in user
- A user cannot like the same article multiple times

---

## US5 — Post an opinion
**As a user**, I want to post a short opinion on an article, so I can share my perspective in a concise, social-media style format.  
**Acceptance Criteria**
- User can post an opinion (≤ 280 characters)
- Opinion is linked to a specific article
- Username and timestamp are displayed
- Empty opinions are not allowed
- Opinions are displayed in reverse chronological order

---

## US6 — Filter by category + country
**As a user**, I want to filter the news feed by category and country, so I can focus on news relevant to me.  
**Acceptance Criteria**
- User can select a category (technology, business, sports, etc.)
- User can select a country/region (e.g., US, IN, GB)
- Feed refreshes based on selected filters
- Selected filters persist during the session (page refresh optional for MVP)

---

## US7 — Search headlines
**As a user**, I want to search for news by keyword, so I can quickly find articles about a topic.  
**Acceptance Criteria**
- Search box accepts keywords
- Search triggers when user presses Enter or clicks Search
- Results show title, source, date, summary
- Empty search shows a friendly prompt (no API call)

---

## US8 — Refresh feed
**As a user**, I want to refresh the feed, so I can see the latest news without reloading the page.  
**Acceptance Criteria**
- “Refresh” button refetches the feed
- UI indicates refresh in progress
- Feed updates with newest content

---

## US9 — Send friend request
**As a user**, I want to send a friend request to another user, so I can connect with people I know.  
**Acceptance Criteria**
- User can search for another user by username
- User can click “Add Friend” to send a request
- Request appears as “Pending” for the sender
- Duplicate requests are prevented

---

## US10 — Accept/decline friend requests
**As a user**, I want to accept or decline incoming friend requests, so I control who can connect with me.  
**Acceptance Criteria**
- Incoming requests appear in a “Friend Requests” section
- User can accept or decline each request
- On accept, both users appear in each other’s friends list
- On decline, request is removed and cannot be accepted later unless re-sent

---

## US11 — View friends list
**As a user**, I want to see my friends list, so I can manage and message them.  
**Acceptance Criteria**
- Friends list page shows friend usernames + avatar (optional)
- Friends list supports basic search/filter
- User can remove a friend (unfriend) and it updates for both users

---

## US12 — Start a 1:1 chat with a friend
**As a user**, I want to start a private chat with a friend, so I can message them directly.  
**Acceptance Criteria**
- User can open chat from a friend’s profile or friends list
- If a conversation exists, it is reused (no duplicates)
- Chat screen shows message history (latest first or bottom anchored)

---

## US13 — Message delivery + read indicators (optional MVP+)
**As a user**, I want to see whether my message was delivered/read, so I know it was seen.  
**Acceptance Criteria**
- Each message has a status: sent → delivered → read
- Status updates in real-time (or on refresh)

---

## US14 — Sign up (register)
**As a user**, I want to create an account, so I can like/save/comment and chat with friends.  
**Acceptance Criteria**
- User can register with username + email + password
- Username is unique (case-insensitive)
- Password rules enforced
- On success, user is logged in automatically (token issued)
- Clear validation errors

---

## US15 — Login
**As a user**, I want to log in, so I can access my personalized features (saved, likes, friends, chat).  
**Acceptance Criteria**
- User can log in using email (or username) + password
- On success, user receives auth token (JWT or session)
- Auth state persists after refresh (token stored securely)
- Wrong credentials show a friendly error message

---

## US16 — Logout
**As a user**, I want to log out, so I can securely end my session.  
**Acceptance Criteria**
- Logout clears stored auth tokens/session
- User is redirected to login/home
- Protected pages require re-login after logout

---

## US17 — Protected access
**As a user**, I should not be able to access protected features unless logged in, so my account is secure.  
**Acceptance Criteria**
- Like/Save/Comment/Post/Chat require login
- If not logged in, user is redirected to login (or shown a prompt)

---

## US18 — View profile
**As a user**, I want to view my profile and others’ profiles, so I can see identity and activity.  
**Acceptance Criteria**
- Profile shows: username, display name, bio, avatar, joined date (optional)
- Profile shows counts: posts, likes (optional), friends
- If user not found, show 404-friendly page

---

## US19 — Edit profile
**As a user**, I want to edit my profile, so I can personalize my identity.  
**Acceptance Criteria**
- User can edit: display name, bio, avatar
- Changes persist after refresh
- Invalid inputs show inline validation errors
- Only the logged-in user can edit their own profile

---

## US20 — Change password
**As a user**, I want to change my password, so I can secure my account if needed.  
**Acceptance Criteria**
- Requires current password + new password
- Shows confirmation message on success

---

## 2. Issues planned to address in Sprint 1
Sprint 1 focus was **project setup + CI + basic end-to-end demo** (placeholders allowed).

### Planned Sprint 1 items (setup + foundation)
- Initialize monorepo structure (`frontend/` CRA TS + `backend/` Go)
- Add GitHub Actions CI (frontend build, backend test/vet)
- Backend: add `/health` endpoint and skeleton server
- Identify a “real-time news API” option and plan integration (keys protected server-side)
- Frontend: basic feed page showing **title/source/date + summary placeholder**, plus error state

### Related user stories targeted (partial / foundation)
- **US1** (partial: display feed UI with placeholder content)
- **US3** (partial: click opens link/new tab)
- **US2** (not implemented yet; only placeholder summary)
- US4–US20 are not in Sprint 1 scope

---

## 3. Completed successfully
- **Project initialization + repo structure**
  - Created `frontend/` (React/CRA TS) and `backend/` (Go module) layout
  - Basic run commands established for local dev
- **Backend skeleton created**
  - `/health` endpoint implemented
  - Backend folder structure created (`cmd/api`, `internal`)
- **Frontend foundation**
  - Basic page that can render a list of articles with placeholder summaries
  - Basic “open full article” behavior available (via link)
- **News provider finalized**
  - Selected **NewsAPI** as the real-time news provider for the project
  - Agreed that API keys will be stored server-side and never exposed in the frontend

---

## 4. Not completed this sprint (and why)
- **Real-time news API integration (provider wired end-to-end)**
  - **Why:** provider was selected, but backend proxy endpoint (`/api/news`) and normalized response shape are not fully implemented/connected to the frontend yet.
  - **Next step:** implement `/api/news` proxy in Go using NewsAPI; return normalized fields (title, source, date, url, summary); connect frontend feed to this endpoint.

- **US2 Auto 50-word summary**
  - **Why:** summarization logic not implemented yet; UI currently uses placeholders.
  - **Next step:** implement server-side summarization (or a temporary truncation fallback), then replace placeholders.

---
