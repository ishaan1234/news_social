import React from 'react';
import App from '../../src/App';
import Home from '../../src/pages/Home';

describe('Home page', () => {
  it('shows placeholder news and moves to the next story', () => {
    cy.viewport(900, 700);
    cy.mount(<Home />);

    cy.contains('Headline placeholder').should('be.visible');
    cy.get('[data-cy="news-prev"]').should('be.disabled');
    cy.get('[data-cy="news-next"]').click();
    cy.get('[data-cy="news-prev"]').should('not.be.disabled');
    cy.get('[data-cy="news-track"]')
      .invoke('attr', 'style')
      .should('include', 'translateY(-100%)');
  });

  it('navigates to chat from the navbar', () => {
    cy.viewport(900, 700);
    cy.mount(<App />);

    cy.get('[data-cy="nav-chat"]').filter(':visible').first().click();
    cy.window().its('location.hash').should('eq', '#/chat');
    cy.get('[data-cy="chat-page"]').should('be.visible');
  });

  it('clicks through every navbar link one by one', () => {
    cy.viewport(900, 700);
    cy.mount(<App />);

    const navChecks = [
      {
        navId: 'nav-news',
        expectedHash: '#/',
        expectedText: 'Headline placeholder',
      },
      {
        navId: 'nav-posts',
        expectedHash: '#/posts',
        expectedText: 'Posts are not wired up yet.',
      },
      {
        navId: 'nav-chat',
        expectedHash: '#/chat',
        expectedText: 'Chat with your people',
      },
      {
        navId: 'nav-profile',
        expectedHash: '#/profile',
        expectedText: 'Profile is still a placeholder.',
      },
      {
        navId: 'nav-settings',
        expectedHash: '#/settings',
        expectedText: 'Settings are still pending.',
      },
    ];

    navChecks.forEach(({ navId, expectedHash, expectedText }) => {
      cy.get(`[data-cy="${navId}"]`).filter(':visible').first().click();
      cy.window().its('location.hash').should('eq', expectedHash);
      cy.contains(expectedText).should('be.visible');
    });
  });
});
