import React from 'react';
import Chat from '../../src/pages/Chat';

describe('Chat page', () => {
  it('filters the conversation list', () => {
    cy.mount(<Chat />);

    cy.get('[data-cy="chat-search"]').type('Nina');
    cy.get('[data-cy="conversation-nina-patel"]').should('be.visible');
    cy.get('[data-cy="conversation-maya-chen"]').should('not.exist');
  });

  it('sends a message from the composer', () => {
    const message = 'Simple Cypress test message';

    cy.mount(<Chat />);
    cy.get('[data-cy="chat-send"]').should('be.disabled');
    cy.get('[data-cy="chat-draft"]').type(message);
    cy.get('[data-cy="chat-send"]').click();
    cy.get('[data-cy="chat-messages"]').contains(message).should('be.visible');
    cy.get('[data-cy="chat-draft"]').should('have.value', '');
  });
});
