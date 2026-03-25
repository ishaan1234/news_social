import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

afterEach(() => {
  window.location.hash = '#/';
});

test('renders the news feed by default', () => {
  window.location.hash = '#/';
  render(<App />);

  expect(screen.getByText(/newshub/i)).toBeInTheDocument();
  expect(screen.getAllByText(/headline placeholder/i).length).toBeGreaterThan(0);
});

test('renders the chat interface on the chat route', () => {
  window.location.hash = '#/chat';
  render(<App />);

  expect(screen.getByText(/chat with your people/i)).toBeInTheDocument();
  expect(screen.getByPlaceholderText(/type a message/i)).toBeInTheDocument();
});
