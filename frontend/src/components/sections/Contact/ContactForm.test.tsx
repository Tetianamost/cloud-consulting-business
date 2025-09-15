import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import ContactForm from './ContactForm';
import { apiService } from '../../../services/api';

// Mock the API service
jest.mock('../../../services/api', () => ({
  apiService: {
    createInquiry: jest.fn(),
  },
}));

// Mock framer-motion to avoid animation issues in tests
jest.mock('framer-motion', () => ({
  motion: {
    div: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
  },
  AnimatePresence: ({ children }: any) => <>{children}</>,
}));

// Mock styled-components keyframes
jest.mock('styled-components', () => {
  const actual = jest.requireActual('styled-components');
  return {
    ...actual,
    keyframes: () => 'mocked-keyframes',
  };
});

const mockApiService = apiService as jest.Mocked<typeof apiService>;

describe('ContactForm Enhanced Features', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('shows real-time validation feedback', async () => {
    render(<ContactForm />);

    const nameInput = screen.getByLabelText(/full name/i);
    const emailInput = screen.getByLabelText(/email address/i);

    // Test invalid input shows error
    await userEvent.type(nameInput, 'A'); // Too short
    await userEvent.tab(); // Trigger blur

    await waitFor(() => {
      expect(screen.getByText(/name is too short/i)).toBeInTheDocument();
    });

    // Test valid input shows success indicator
    await userEvent.clear(nameInput);
    await userEvent.type(nameInput, 'John Doe');
    await userEvent.type(emailInput, 'john@example.com');

    // The form should show visual success indicators for valid fields
    expect(nameInput).toHaveStyle('border-color: #10b981'); // success color
  });

  test('displays loading states during form submission', async () => {

    // Mock successful API response
    mockApiService.createInquiry.mockResolvedValue({
      success: true,
      data: {
        id: 'test-inquiry-123',
        name: 'John Doe',
        email: 'john@example.com',
        company: 'Test Company',
        phone: '+1234567890',
        services: ['assessment'],
        message: 'Test message',
        status: 'new',
        priority: 'medium',
        source: 'contact_form',
        created_at: '2025-01-01T00:00:00Z',
        updated_at: '2025-01-01T00:00:00Z',
      },
    });

    render(<ContactForm />);

    // Fill out the form
    await userEvent.type(screen.getByLabelText(/full name/i), 'John Doe');
    await userEvent.type(screen.getByLabelText(/email address/i), 'john@example.com');
    await userEvent.type(screen.getByLabelText(/message/i), 'Test message');
    await userEvent.click(screen.getByLabelText(/cloud assessment/i));

    const submitButton = screen.getByRole('button', { name: /send message/i });
    await userEvent.click(submitButton);

    // Should show loading state
    await waitFor(() => {
      expect(screen.getByText(/processing/i)).toBeInTheDocument();
    });

    // Should show progress indicator
    expect(screen.getByText(/validating information/i)).toBeInTheDocument();
  });

  test('shows professional success message with next steps', async () => {

    // Mock successful API response
    mockApiService.createInquiry.mockResolvedValue({
      success: true,
      data: {
        id: 'test-inquiry-123',
        name: 'John Doe',
        email: 'john@example.com',
        company: 'Test Company',
        phone: '+1234567890',
        services: ['assessment'],
        message: 'Test message',
        status: 'new',
        priority: 'medium',
        source: 'contact_form',
        created_at: '2025-01-01T00:00:00Z',
        updated_at: '2025-01-01T00:00:00Z',
      },
    });

    render(<ContactForm />);

    // Fill out and submit form
    await userEvent.type(screen.getByLabelText(/full name/i), 'John Doe');
    await userEvent.type(screen.getByLabelText(/email address/i), 'john@example.com');
    await userEvent.type(screen.getByLabelText(/message/i), 'Test message');
    await userEvent.click(screen.getByLabelText(/cloud assessment/i));
    await userEvent.click(screen.getByRole('button', { name: /send message/i }));

    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText(/inquiry submitted successfully/i)).toBeInTheDocument();
    });

    // Check for professional success elements
    expect(screen.getByText(/reference id/i)).toBeInTheDocument();
    expect(screen.getByText(/test-inquiry-123/)).toBeInTheDocument();
    expect(screen.getByText(/what happens next/i)).toBeInTheDocument();
    expect(screen.getByText(/confirmation email within 30 seconds/i)).toBeInTheDocument();
    expect(screen.getByText(/consultant will review and respond within 24 hours/i)).toBeInTheDocument();
  });

  test('shows clear error messages on submission failure', async () => {
    // Mock API error
    mockApiService.createInquiry.mockRejectedValue(new Error('Network error'));

    render(<ContactForm />);

    // Fill out and submit form
    await userEvent.type(screen.getByLabelText(/full name/i), 'John Doe');
    await userEvent.type(screen.getByLabelText(/email address/i), 'john@example.com');
    await userEvent.type(screen.getByLabelText(/message/i), 'Test message');
    await userEvent.click(screen.getByLabelText(/cloud assessment/i));
    await userEvent.click(screen.getByRole('button', { name: /send message/i }));

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText(/submission failed/i)).toBeInTheDocument();
    });

    // Check for helpful error information
    expect(screen.getByText(/network error/i)).toBeInTheDocument();
    expect(screen.getByText(/email us directly/i)).toBeInTheDocument();
    expect(screen.getByText(/info@cloudpartner.pro/)).toBeInTheDocument();
  });

  test('validates required fields with clear messages', async () => {
    render(<ContactForm />);

    const submitButton = screen.getByRole('button', { name: /send message/i });
    await userEvent.click(submitButton);

    // Should show validation errors for required fields
    await waitFor(() => {
      expect(screen.getByText(/name is required/i)).toBeInTheDocument();
      expect(screen.getByText(/email is required/i)).toBeInTheDocument();
      expect(screen.getByText(/message is required/i)).toBeInTheDocument();
      expect(screen.getByText(/please select at least one service/i)).toBeInTheDocument();
    });
  });

  test('disables submit button during submission', async () => {

    // Mock slow API response
    mockApiService.createInquiry.mockImplementation(
      () => new Promise(resolve => setTimeout(() => resolve({
        success: true,
        data: {
          id: 'test-inquiry-123',
          name: 'John Doe',
          email: 'john@example.com',
          company: 'Test Company',
          phone: '+1234567890',
          services: ['assessment'],
          message: 'Test message',
          status: 'new',
          priority: 'medium',
          source: 'contact_form',
          created_at: '2025-01-01T00:00:00Z',
          updated_at: '2025-01-01T00:00:00Z',
        },
      }), 1000))
    );

    render(<ContactForm />);

    // Fill out form
    await userEvent.type(screen.getByLabelText(/full name/i), 'John Doe');
    await userEvent.type(screen.getByLabelText(/email address/i), 'john@example.com');
    await userEvent.type(screen.getByLabelText(/message/i), 'Test message');
    await userEvent.click(screen.getByLabelText(/cloud assessment/i));

    const submitButton = screen.getByRole('button', { name: /send message/i });
    await userEvent.click(submitButton);

    // Button should be disabled during submission
    await waitFor(() => {
      expect(submitButton).toBeDisabled();
    });
  });
});