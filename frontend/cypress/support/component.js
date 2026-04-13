import './component.css';
import { mount } from 'cypress/react';

Cypress.Commands.add('mount', mount);

beforeEach(() => {
  window.localStorage.clear();
  window.location.hash = '#/';
});
