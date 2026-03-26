# `Posts.tsx` Cypress Test Summary

## Direct test coverage

| Test name | Broad feature / module | Description | Test file |
|---|---|---|---|
| `supports follow, like, share, and comments` | `Posts page social interactions` | Mounts the `Posts` page, verifies the composer is visible, toggles follow state for `post-1`, increments the like count, increments the share count, submits a new comment, and confirms the comment count updates. | [`frontend/cypress/component/posts.cy.jsx`](../cypress/component/posts.cy.jsx) |

## Indirect coverage related to `Posts.tsx`

| Test name | Broad feature / module | Description | Test file |
|---|---|---|---|
| `clicks through every navbar link one by one` | `App routing / Posts page entry` | Verifies that clicking the `Posts` navbar item routes to `#/posts` and renders the `Share an opinion` heading, confirming the page is wired into app navigation. | [`frontend/cypress/component/home.cy.jsx`](../cypress/component/home.cy.jsx) |

## Broad feature breakdown for the direct `Posts` test

| Feature / module | Description | Relevant source |
|---|---|---|
| `Composer visibility` | Checks that the post composer section renders by asserting `Share an opinion` is visible. | [`frontend/src/pages/Posts.tsx`](./pages/Posts.tsx) |
| `Follow / unfollow` | Verifies the follow button for `post-1` changes from `Following` to `Follow` when clicked. | [`frontend/src/pages/Posts.tsx`](./pages/Posts.tsx) |
| `Like reaction` | Verifies clicking like updates the count from `14` to `15`. | [`frontend/src/pages/Posts.tsx`](./pages/Posts.tsx) |
| `Share reaction` | Verifies clicking share updates the count from `3` to `4`. | [`frontend/src/pages/Posts.tsx`](./pages/Posts.tsx) |
| `Comments` | Verifies typing and submitting a comment adds it to the UI and updates the comment count from `1` to `2`. | [`frontend/src/pages/Posts.tsx`](./pages/Posts.tsx) |

## Not currently covered in Cypress for `Posts.tsx`

- Creating a brand new post from the composer
- Selecting a different attached news summary
- Empty-input validation for post submission
- Disabled-state behavior for submit buttons
