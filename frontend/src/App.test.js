import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import App from './App';
import { shortenUrl, getStats } from './api';

jest.mock('./api', () => ({ shortenUrl: jest.fn(), getStats: jest.fn() }));
beforeEach(() => jest.resetAllMocks());

test('does not submit an empty URL', () => {
  render(<App />);
  expect(screen.getByRole('button', { name: 'Shorten' })).toBeDisabled();
});

test('shortens a URL and displays its click count', async () => {
  shortenUrl.mockResolvedValue({ short_url: 'http://localhost:8080/g/abc123' });
  getStats.mockResolvedValue({ clicks: '5' });
  render(<App />);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: 'https://example.com' } });
  fireEvent.click(screen.getByRole('button', { name: 'Shorten' }));
  expect(await screen.findByRole('link')).toHaveAttribute('href', 'http://localhost:8080/g/abc123');
  await waitFor(() => expect(getStats).toHaveBeenCalledWith('abc123'));
  expect(shortenUrl).toHaveBeenCalledWith('https://example.com');
  expect(await screen.findByText('5', { exact: false })).toBeInTheDocument();
});

test('shows a useful error when the API rejects the request', async () => {
  shortenUrl.mockRejectedValue(new Error('Network error'));
  render(<App />);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: 'https://example.com' } });
  fireEvent.click(screen.getByRole('button', { name: 'Shorten' }));
  expect(await screen.findByRole('alert')).toHaveTextContent('URL kısaltılamadı');
  expect(getStats).not.toHaveBeenCalled();
  expect(screen.getByRole('button', { name: 'Shorten' })).toBeEnabled();
});
