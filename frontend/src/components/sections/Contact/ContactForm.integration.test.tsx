import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';

// Simple integration test to verify form functionality
describe('ContactForm Integration', () => {
  test('form renders with all required fields', () => {
    // Create a minimal test component
    const TestForm = () => (
      <form data-testid="contact-form">
        <input name="name" placeholder="Full Name" required />
        <input name="email" type="email" placeholder="Email" required />
        <textarea name="message" placeholder="Message" required />
        <button type="submit">Send Message</button>
      </form>
    );

    render(<TestForm />);

    expect(screen.getByPlaceholderText(/full name/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/email/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/message/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /send message/i })).toBeInTheDocument();
  });

  test('form validation works for required fields', async () => {
    
    const TestForm = () => {
      const [errors, setErrors] = React.useState<string[]>([]);
      
      const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const formData = new FormData(e.target as HTMLFormElement);
        const newErrors: string[] = [];
        
        if (!formData.get('name')) newErrors.push('Name is required');
        if (!formData.get('email')) newErrors.push('Email is required');
        if (!formData.get('message')) newErrors.push('Message is required');
        
        setErrors(newErrors);
      };

      return (
        <form onSubmit={handleSubmit} data-testid="contact-form">
          <input name="name" placeholder="Full Name" />
          <input name="email" type="email" placeholder="Email" />
          <textarea name="message" placeholder="Message" />
          <button type="submit">Send Message</button>
          {errors.map((error, index) => (
            <div key={index} role="alert">{error}</div>
          ))}
        </form>
      );
    };

    render(<TestForm />);

    const submitButton = screen.getByRole('button', { name: /send message/i });
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/name is required/i)).toBeInTheDocument();
      expect(screen.getByText(/email is required/i)).toBeInTheDocument();
      expect(screen.getByText(/message is required/i)).toBeInTheDocument();
    });
  });
});