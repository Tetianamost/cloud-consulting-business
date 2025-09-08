import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { V0EmailDeliveryDashboardConnected } from './V0EmailDeliveryDashboardConnected';
import apiService from '../../services/api';

// Mock the API service
jest.mock('../../services/api');
const mockApiService = apiService as jest.Mocked<typeof apiService>;

describe('V0EmailDeliveryDashboardConnected', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should display error state when no real data is available', async () => {
    // Mock API calls to return empty/error responses
    mockApiService.getSystemMetrics.mockRejectedValue(new Error('API not available'));
    mockApiService.getEmailMetrics.mockRejectedValue(new Error('Email metrics not available'));
    mockApiService.listInquiries.mockResolvedValue({
      success: true,
      data: [],
      count: 0,
      total: 0,
      page: 1,
      pages: 1
    });

    render(<V0EmailDeliveryDashboardConnected />);

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.getByText('Email Metrics Unavailable')).toBeInTheDocument();
    });

    expect(screen.getByText(/No email monitoring data available/)).toBeInTheDocument();
    expect(screen.getByText('Retry')).toBeInTheDocument();
  });

  it('should display real data when system metrics are available', async () => {
    // Mock API calls to return real data
    mockApiService.getSystemMetrics.mockResolvedValue({
      success: true,
      data: {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 20,
        email_delivery_rate: 95.0,
        avg_report_gen_time_ms: 1500,
        system_uptime: '2d 5h',
        last_processed_at: new Date().toISOString()
      }
    });

    mockApiService.getEmailMetrics.mockRejectedValue(new Error('Not implemented'));
    mockApiService.listInquiries.mockResolvedValue({
      success: true,
      data: [],
      count: 0,
      total: 0,
      page: 1,
      pages: 1
    });

    render(<V0EmailDeliveryDashboardConnected />);

    // Wait for loading to complete and data to be displayed
    await waitFor(() => {
      expect(screen.getByText('Email Delivery Monitoring')).toBeInTheDocument();
    });

    // Should show real metrics, not error state
    expect(screen.queryByText('Email Metrics Unavailable')).not.toBeInTheDocument();
    expect(screen.queryByText('Using demo data')).not.toBeInTheDocument();
    
    // Should display delivery rate from real data
    await waitFor(() => {
      expect(screen.getAllByText('95.0%')[0]).toBeInTheDocument(); // Use getAllByText to handle multiple matches
    });
  });

  it('should display real data when system metrics are available even with some API failures', async () => {
    // Mock system metrics to succeed but other calls to fail
    mockApiService.getSystemMetrics.mockResolvedValue({
      success: true,
      data: {
        total_inquiries: 5,
        reports_generated: 3,
        emails_sent: 10,
        email_delivery_rate: 90.0,
        avg_report_gen_time_ms: 1200,
        system_uptime: '1d 3h',
        last_processed_at: new Date().toISOString()
      }
    });

    mockApiService.getEmailMetrics.mockRejectedValue(new Error('Email events not available'));
    mockApiService.listInquiries.mockRejectedValue(new Error('Inquiries not available'));

    render(<V0EmailDeliveryDashboardConnected />);

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.getByText('Email Delivery Monitoring')).toBeInTheDocument();
    });

    // Should not show error state since we have real system metrics
    expect(screen.queryByText('Email Metrics Unavailable')).not.toBeInTheDocument();
    
    // Should display the available data
    await waitFor(() => {
      expect(screen.getAllByText('90.0%')[0]).toBeInTheDocument(); // Use getAllByText to handle multiple matches
    });
  });

  it('should display error state when no emails have been sent', async () => {
    // Mock API to return metrics with no email activity
    mockApiService.getSystemMetrics.mockResolvedValue({
      success: true,
      data: {
        total_inquiries: 0,
        reports_generated: 0,
        emails_sent: 0,
        email_delivery_rate: 0,
        avg_report_gen_time_ms: 0,
        system_uptime: '1h',
        last_processed_at: new Date().toISOString()
      }
    });

    mockApiService.getEmailMetrics.mockRejectedValue(new Error('No data'));
    mockApiService.listInquiries.mockResolvedValue({
      success: true,
      data: [],
      count: 0,
      total: 0,
      page: 1,
      pages: 1
    });

    render(<V0EmailDeliveryDashboardConnected />);

    // Should show error state since no emails have been sent
    await waitFor(() => {
      expect(screen.getByText('Email Metrics Unavailable')).toBeInTheDocument();
    });

    expect(screen.getByText(/No emails have been sent yet/)).toBeInTheDocument();
  });
});