import React from 'react';
import Profile from '../../src/pages/Profile';

describe('Profile page', () => {
  it('switches between timeline tabs', () => {
    cy.viewport(1100, 800);
    cy.mount(<Profile />);

    cy.contains('Pinned').should('be.visible');

    cy.get('[data-cy="profile-tab-replies"]').click();
    cy.contains('Replying to @maya').should('be.visible');

    cy.get('[data-cy="profile-tab-media"]').click();
    cy.contains('Policy briefing card').should('be.visible');

    cy.get('[data-cy="profile-tab-likes"]').click();
    cy.contains('Liked a post from @nina').should('be.visible');
  });

  it('opens the editor, saves a profile change, and closes the modal', () => {
    cy.viewport(1100, 800);
    cy.mount(<Profile />);

    cy.get('[data-cy="profile-edit"]').click();
    cy.get('[data-cy="profile-edit-modal"]').should('be.visible');

    cy.get('[data-cy="profile-name-input"]').clear().type('Avery Stone Live');
    cy.get('[data-cy="profile-save"]').click();

    cy.get('[data-cy="profile-edit-modal"]').should('not.exist');
    cy.contains('h1', 'Avery Stone Live').should('be.visible');
  });

  it('opens the editor and closes it with cancel', () => {
    cy.viewport(1100, 800);
    cy.mount(<Profile />);

    cy.get('[data-cy="profile-edit"]').click();
    cy.get('[data-cy="profile-edit-modal"]').should('be.visible');
    cy.get('[data-cy="profile-cancel"]').click();
    cy.get('[data-cy="profile-edit-modal"]').should('not.exist');
  });

  it('toggles follow buttons in the sidebar', () => {
    cy.viewport(1100, 800);
    cy.mount(<Profile />);

    cy.get('[data-cy="profile-follow-follow-1"]')
      .should('contain.text', 'Follow')
      .click()
      .should('contain.text', 'Following');

    cy.get('[data-cy="profile-follow-follow-1"]')
      .click()
      .should('contain.text', 'Follow');
  });
});
