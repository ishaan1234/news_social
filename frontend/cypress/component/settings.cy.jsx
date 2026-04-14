import React from 'react';
import Settings from '../../src/pages/Settings';

describe('Settings page', () => {
  beforeEach(() => {
    cy.viewport(900, 700);
    cy.mount(<Settings />);
  });

  it('clicks feed and audience controls and updates the account snapshot', () => {
    cy.get('[data-cy="settings-feed-mode-latest"]').click();
    cy.get('[data-cy="settings-snapshot-feed-mode"]').should(
      'contain.text',
      'Latest'
    );

    cy.get('[data-cy="settings-summary-style"]').select('Deep dive');
    cy.get('[data-cy="settings-snapshot-summary-style"]').should(
      'contain.text',
      'Deep dive'
    );

    cy.get('[data-cy="settings-audience-followers"]').click().should(
      'have.class',
      'bg-slate-900'
    );
  });

  it('toggles notification and privacy switches', () => {
    cy.get('[data-cy="settings-trend-alerts"]')
      .should('have.attr', 'aria-pressed', 'false')
      .click()
      .should('have.attr', 'aria-pressed', 'true');

    cy.get('[data-cy="settings-read-receipts"]')
      .should('have.attr', 'aria-pressed', 'false')
      .click()
      .should('have.attr', 'aria-pressed', 'true');

    cy.get('[data-cy="settings-dm-access-nobody"]').click();
    cy.get('[data-cy="settings-snapshot-dm-access"]').should(
      'contain.text',
      'Nobody'
    );
  });

  it('scrolls through the page and clicks a quick link', () => {
    cy.scrollTo('bottom');
    cy.get('[data-cy="settings-quick-links"]').should('be.visible');

    cy.get('[data-cy="settings-hide-new-account-requests"]')
      .scrollIntoView()
      .should('be.visible')
      .click()
      .should('have.attr', 'aria-pressed', 'true');

    cy.get('[data-cy="settings-quick-link-chat"]').click();
    cy.window().its('location.hash').should('eq', '#/chat');
  });
});
