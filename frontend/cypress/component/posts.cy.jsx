import React from 'react';
import Posts from '../../src/pages/Posts';

describe('Posts page', () => {
  it('supports follow, like, share, and comments', () => {
    cy.viewport(1100, 800);
    cy.mount(<Posts />);

    cy.contains('Share an opinion').should('be.visible');

    cy.get('[data-cy="follow-post-1"]').should('contain.text', 'Following');
    cy.get('[data-cy="follow-post-1"]').click().should('contain.text', 'Follow');

    cy.get('[data-cy="like-post-1"]').should('contain.text', 'Like 14');
    cy.get('[data-cy="like-post-1"]').click().should('contain.text', 'Like 15');

    cy.get('[data-cy="share-post-1"]').click().should('contain.text', 'Share 4');

    cy.get('[data-cy="comment-draft-post-1"]').type('Frontend-only comment.');
    cy.get('[data-cy="comment-submit-post-1"]').click();
    cy.contains('Frontend-only comment.').should('be.visible');
    cy.get('[data-cy="comment-trigger-post-1"]').should('contain.text', 'Comment 2');
  });
});
