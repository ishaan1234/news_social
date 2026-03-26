import React from 'react';
import App from '../../src/App';
import Home from '../../src/pages/Home';

const mockNewsResponse = {
  status: 'ok',
  totalResults: 2,
  articles: [
    {
      source: { name: 'Electrek' },
      title: 'Tesla expands in Canada',
      description: 'Expansion news',
      url: 'https://example.com/article-1',
      urlToImage: 'https://picsum.photos/seed/live1/800/600',
      publishedAt: '2026-03-25T10:00:00Z',
      summary: 'Tesla is expanding in Canada through new locations and market growth.',
    },
    {
      source: { name: 'Reuters' },
      title: 'EV market momentum continues',
      description: 'Momentum update',
      url: 'https://example.com/article-2',
      urlToImage: 'https://picsum.photos/seed/live2/800/600',
      publishedAt: '2026-03-25T09:00:00Z',
      summary: 'Electric vehicle demand remains strong across major markets.',
    },
  ],
};

describe('Home page', () => {
  beforeEach(() => {
    cy.intercept('GET', '**/news*', mockNewsResponse).as('getNews');
  });

  it('shows placeholder news and moves to the next story', () => {
    cy.viewport(900, 700);
    cy.mount(<Home />);
    cy.wait('@getNews');

    cy.contains('Tesla expands in Canada').should('be.visible');
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
    cy.wait('@getNews');

    cy.get('[data-cy="nav-chat"]').filter(':visible').first().click();
    cy.window().its('location.hash').should('eq', '#/chat');
    cy.get('[data-cy="chat-page"]').should('be.visible');
  });

  it('clicks through every navbar link one by one', () => {
    cy.viewport(900, 700);
    cy.mount(<App />);
    cy.wait('@getNews');

    const navChecks = [
      {
        navId: 'nav-news',
        expectedHash: '#/',
        expectedText: 'Tesla expands in Canada',
      },
      {
        navId: 'nav-posts',
        expectedHash: '#/posts',
        expectedText: 'Share an opinion',
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
