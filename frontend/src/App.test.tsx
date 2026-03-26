import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

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

beforeEach(() => {
  Object.defineProperty(global, 'fetch', {
    writable: true,
    value: jest.fn().mockResolvedValue({
      ok: true,
      json: async () => mockNewsResponse,
    }),
  });
});

afterEach(() => {
  window.location.hash = '#/';
  jest.resetAllMocks();
});

test('renders the live news feed by default', async () => {
  window.location.hash = '#/';
  render(<App />);

  expect(screen.getByText(/newshub/i)).toBeInTheDocument();
  expect(await screen.findByText(/tesla expands in canada/i)).toBeInTheDocument();
});

test('renders the chat interface on the chat route', () => {
  window.location.hash = '#/chat';
  render(<App />);

  expect(screen.getByText(/chat with your people/i)).toBeInTheDocument();
  expect(screen.getByPlaceholderText(/type a message/i)).toBeInTheDocument();
});
