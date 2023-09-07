import { render, screen } from '@testing-library/react';
import App from './App';

test('renders wasm loading', () => {
  render(<App />);
  const linkElement = screen.getByText(/Loading WebAssembly module/i);
  expect(linkElement).toBeInTheDocument();
});
