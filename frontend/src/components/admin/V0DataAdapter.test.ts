import { V0DataAdapter } from './V0DataAdapter';
import { SystemMetrics, EmailStatus } from '../../services/api';

describe('V0DataAdapter', () => {
  describe('safeAdaptEmailMetrics', () => {
    it('should return null when no real data is available', () => {
      const result = V0DataAdapter.safeAdaptEmailMetrics(null, []);
      expect(result).toBeNull();
    });

    it('should return null when system metrics are empty and no email statuses', () => {
      const emptyMetrics: SystemMetrics = {
        total_inquiries: 0,
        reports_generated: 0,
        emails_sent: 0,
        email_delivery_rate: 0,
        avg_report_gen_time_ms: 0,
        system_uptime: '0h',
      };

      const result = V0DataAdapter.safeAdaptEmailMetrics(emptyMetrics, []);
      expect(result).toBeNull();
    });

    it('should adapt system metrics with real email data', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 25,
        reports_generated: 15,
        emails_sent: 50,
        email_delivery_rate: 92.5,
        avg_report_gen_time_ms: 1250,
        system_uptime: '3d 7h 22m',
      };

      const result = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, []);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(50);
      expect(result!.deliveryRate).toBe(92.5);
      expect(result!.delivered).toBe(46); // 92.5% of 50
      expect(result!.failedEmails).toBe(4); // 50 - 46
      expect(result!.openRate).toBeGreaterThan(0);
      expect(result!.clickRate).toBeGreaterThan(0);
    });

    it('should adapt email statuses when system metrics are null', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'customer@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: '2024-01-01T10:00:00Z',
          delivered_at: '2024-01-01T10:05:00Z',
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'customer2@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T11:00:00Z',
          delivered_at: '2024-01-01T11:00:00Z',
          error_message: 'SMTP connection failed',
        },
        {
          inquiry_id: 'inquiry-3',
          customer_email: 'customer3@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: '2024-01-01T12:00:00Z',
          delivered_at: '2024-01-01T12:03:00Z',
          error_message: '',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(null, emailStatuses);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(3);
      expect(result!.delivered).toBe(2);
      expect(result!.failedEmails).toBe(1);
      expect(result!.deliveryRate).toBeCloseTo(66.67, 1); // 2/3 * 100
      expect(result!.bounced).toBe(0);
      expect(result!.spam).toBe(0);
    });

    it('should categorize bounced emails correctly', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'invalid@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T10:00:00Z',
          delivered_at: '2024-01-01T10:00:00Z',
          error_message: 'Email address does not exist',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'bounced@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T11:00:00Z',
          delivered_at: '2024-01-01T11:00:00Z',
          error_message: 'Permanent bounce - invalid recipient',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(null, emailStatuses);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(2);
      expect(result!.bounced).toBe(2); // Both should be categorized as bounced
      expect(result!.failedEmails).toBe(2);
      expect(result!.delivered).toBe(0);
    });

    it('should categorize spam emails correctly', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'spam@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T10:00:00Z',
          delivered_at: '2024-01-01T10:00:00Z',
          error_message: 'Message blocked as spam',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'rejected@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T11:00:00Z',
          delivered_at: '2024-01-01T11:00:00Z',
          error_message: 'Rejected by spam filter',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(null, emailStatuses);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(2);
      expect(result!.spam).toBe(2); // Both should be categorized as spam
      expect(result!.failedEmails).toBe(2);
      expect(result!.delivered).toBe(0);
    });

    it('should handle mixed email statuses correctly', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'delivered@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: '2024-01-01T10:00:00Z',
          delivered_at: '2024-01-01T10:05:00Z',
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'failed@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T11:00:00Z',
          delivered_at: '2024-01-01T11:00:00Z',
          error_message: 'SMTP connection failed',
        },
        {
          inquiry_id: 'inquiry-3',
          customer_email: 'bounced@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T12:00:00Z',
          delivered_at: '2024-01-01T12:00:00Z',
          error_message: 'Email address not found',
        },
        {
          inquiry_id: 'inquiry-4',
          customer_email: 'spam@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T13:00:00Z',
          delivered_at: '2024-01-01T13:00:00Z',
          error_message: 'Blocked as spam',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(null, emailStatuses);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(4);
      expect(result!.delivered).toBe(1);
      expect(result!.failedEmails).toBe(3);
      expect(result!.bounced).toBe(1);
      expect(result!.spam).toBe(1);
      expect(result!.deliveryRate).toBe(25); // 1/4 * 100
    });

    it('should estimate open and click rates based on delivery performance', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 100,
        email_delivery_rate: 95, // High delivery rate
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, []);
      
      expect(result).not.toBeNull();
      expect(result!.openRate).toBeGreaterThan(20); // Should be higher for good delivery rate
      expect(result!.clickRate).toBeGreaterThan(3); // Should be reasonable percentage of opens
      expect(result!.opened).toBeGreaterThan(0);
      expect(result!.clicked).toBeGreaterThan(0);
      expect(result!.clicked).toBeLessThan(result!.opened); // Clicks should be less than opens
    });

    it('should handle error gracefully and return null', () => {
      // Mock console.error to avoid noise in tests
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

      // Pass invalid data that might cause an error
      const invalidMetrics = { invalid: 'data' } as any;
      
      const result = V0DataAdapter.safeAdaptEmailMetrics(invalidMetrics, []);
      
      expect(result).toBeNull();
      expect(consoleSpy).toHaveBeenCalledWith('Error adapting email metrics:', expect.any(Error));
      
      consoleSpy.mockRestore();
    });

    it('should prioritize email statuses over system metrics when both are available', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 20, // System says 20 emails
        email_delivery_rate: 80,
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'delivered@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: '2024-01-01T10:00:00Z',
          delivered_at: '2024-01-01T10:05:00Z',
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'failed@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: '2024-01-01T11:00:00Z',
          delivered_at: '2024-01-01T11:00:00Z',
          error_message: 'SMTP failed',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, emailStatuses);
      
      expect(result).not.toBeNull();
      // Should use the higher of system metrics or email statuses count
      expect(result!.totalEmails).toBe(20); // Uses system metrics (higher)
      // But delivery rate should be calculated from actual statuses when available
      expect(result!.delivered).toBe(1); // From email statuses
      expect(result!.failedEmails).toBe(19); // 20 - 1 delivered
    });
  });

  describe('hasRealEmailData', () => {
    it('should return false when no data is available', () => {
      const result = V0DataAdapter.hasRealEmailData(null, []);
      expect(result).toBe(false);
    });

    it('should return false when system metrics have no email data', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 0,
        email_delivery_rate: 0,
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.hasRealEmailData(systemMetrics, []);
      expect(result).toBe(false);
    });

    it('should return true when system metrics have email data', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 25,
        email_delivery_rate: 90,
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.hasRealEmailData(systemMetrics, []);
      expect(result).toBe(true);
    });

    it('should return true when email statuses are available', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'test@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: '',
        },
      ];

      const result = V0DataAdapter.hasRealEmailData(null, emailStatuses);
      expect(result).toBe(true);
    });

    it('should return true when system metrics have only emails_sent', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 15, // Has email data
        email_delivery_rate: 0,
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.hasRealEmailData(systemMetrics, []);
      expect(result).toBe(true);
    });

    it('should return true when system metrics have only email_delivery_rate', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 0,
        email_delivery_rate: 85.5, // Has email data
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.hasRealEmailData(systemMetrics, []);
      expect(result).toBe(true);
    });
  });

  describe('getEmailDataErrorMessage', () => {
    it('should return specific message for EMAIL_MONITORING_UNAVAILABLE', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(
        null,
        [],
        'EMAIL_MONITORING_UNAVAILABLE: Service not configured'
      );
      
      expect(result).toContain('Email monitoring is not configured');
      expect(result).toContain('Contact your administrator');
    });

    it('should return specific message for EMAIL_MONITORING_UNHEALTHY', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(
        null,
        [],
        'EMAIL_MONITORING_UNHEALTHY: Database connection failed'
      );
      
      expect(result).toContain('Email monitoring system is experiencing issues');
      expect(result).toContain('temporarily unavailable');
    });

    it('should return specific message for EMAIL_STATUS_RETRIEVAL_ERROR', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(
        null,
        [],
        'EMAIL_STATUS_RETRIEVAL_ERROR: Query timeout'
      );
      
      expect(result).toContain('Unable to retrieve email status data');
      expect(result).toContain('monitoring system may be overloaded');
    });

    it('should return specific message for NO_EMAIL_EVENTS', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(
        null,
        [],
        'NO_EMAIL_EVENTS: No events found'
      );
      
      expect(result).toContain('No email events have been recorded yet');
      expect(result).toContain('Email metrics will appear once emails are sent');
    });

    it('should return generic API error message for unknown errors', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(
        null,
        [],
        'UNKNOWN_ERROR: Something went wrong'
      );
      
      expect(result).toContain('Unable to load email metrics');
      expect(result).toContain('UNKNOWN_ERROR: Something went wrong');
    });

    it('should return appropriate message when no API error but no data', () => {
      const result = V0DataAdapter.getEmailDataErrorMessage(null, []);
      
      expect(result).toContain('Email monitoring data is not available');
      expect(result).toContain('email tracking is not configured');
    });

    it('should handle partial data availability', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 10,
        reports_generated: 5,
        emails_sent: 0, // No email data
        email_delivery_rate: 0,
        avg_report_gen_time_ms: 1000,
        system_uptime: '1d',
      };

      const result = V0DataAdapter.getEmailDataErrorMessage(systemMetrics, []);
      
      expect(result).toContain('Email monitoring data is not available');
    });
  });

  describe('adaptEmailStatusesToMetrics', () => {
    it('should return empty metrics for empty array', () => {
      const result = V0DataAdapter.adaptEmailStatusesToMetrics([]);
      
      expect(result.totalEmails).toBe(0);
      expect(result.delivered).toBe(0);
      expect(result.failedEmails).toBe(0);
      expect(result.deliveryRate).toBe(0);
      expect(result.openRate).toBe(0);
      expect(result.clickRate).toBe(0);
    });

    it('should return empty metrics for null input', () => {
      const result = V0DataAdapter.adaptEmailStatusesToMetrics(null as any);
      
      expect(result.totalEmails).toBe(0);
      expect(result.delivered).toBe(0);
      expect(result.failedEmails).toBe(0);
      expect(result.deliveryRate).toBe(0);
      expect(result.openRate).toBe(0);
      expect(result.clickRate).toBe(0);
    });

    it('should correctly calculate metrics from email statuses', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'delivered1@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'delivered2@example.com',
          consultant_email: 'consultant@example.com',
          status: 'delivered',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-3',
          customer_email: 'failed@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'SMTP failed',
        },
      ];

      const result = V0DataAdapter.adaptEmailStatusesToMetrics(emailStatuses);
      
      expect(result.totalEmails).toBe(3);
      expect(result.delivered).toBe(2);
      expect(result.failedEmails).toBe(1);
      expect(result.deliveryRate).toBeCloseTo(66.67, 1); // 2/3 * 100
      expect(result.openRate).toBeGreaterThan(0);
      expect(result.clickRate).toBeGreaterThan(0);
      expect(result.opened).toBeGreaterThan(0);
      expect(result.clicked).toBeGreaterThan(0);
    });

    it('should handle complex error message categorization', () => {
      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-1',
          customer_email: 'bounce1@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'Email address does not exist',
        },
        {
          inquiry_id: 'inquiry-2',
          customer_email: 'bounce2@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'Invalid recipient address',
        },
        {
          inquiry_id: 'inquiry-3',
          customer_email: 'spam1@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'Message blocked by spam filter',
        },
        {
          inquiry_id: 'inquiry-4',
          customer_email: 'spam2@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'Rejected as spam',
        },
        {
          inquiry_id: 'inquiry-5',
          customer_email: 'other@example.com',
          consultant_email: 'consultant@example.com',
          status: 'failed',
          sent_at: new Date().toISOString(),
          delivered_at: new Date().toISOString(),
          error_message: 'Connection timeout',
        },
      ];

      const result = V0DataAdapter.adaptEmailStatusesToMetrics(emailStatuses);
      
      expect(result.totalEmails).toBe(5);
      expect(result.delivered).toBe(0);
      expect(result.failedEmails).toBe(5);
      expect(result.bounced).toBe(2); // 'does not exist' and 'invalid'
      expect(result.spam).toBe(2); // 'blocked' and 'rejected'
      expect(result.deliveryRate).toBe(0);
    });
  });

  describe('integration with real API responses', () => {
    it('should handle typical successful API response', () => {
      const systemMetrics: SystemMetrics = {
        total_inquiries: 42,
        reports_generated: 28,
        emails_sent: 84,
        email_delivery_rate: 94.2,
        avg_report_gen_time_ms: 1850,
        system_uptime: '5d 12h 30m',
        last_processed_at: '2024-01-15T14:30:00Z',
      };

      const emailStatuses: EmailStatus[] = [
        {
          inquiry_id: 'inquiry-real-1',
          customer_email: 'real.customer@company.com',
          consultant_email: 'info@cloudpartner.pro',
          status: 'delivered',
          sent_at: '2024-01-15T10:00:00Z',
          delivered_at: '2024-01-15T10:02:30Z',
          error_message: '',
        },
        {
          inquiry_id: 'inquiry-real-2',
          customer_email: 'another@business.org',
          consultant_email: 'info@cloudpartner.pro',
          status: 'delivered',
          sent_at: '2024-01-15T11:15:00Z',
          delivered_at: '2024-01-15T11:16:45Z',
          error_message: '',
        },
      ];

      const result = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, emailStatuses);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(84); // From system metrics
      expect(result!.deliveryRate).toBe(94.2); // From system metrics
      expect(result!.delivered).toBe(79); // Calculated from delivery rate
      expect(result!.failedEmails).toBe(5); // 84 - 79
      expect(result!.openRate).toBeGreaterThan(20); // Should be reasonable for good delivery
      expect(result!.clickRate).toBeGreaterThan(4); // Should be reasonable percentage
    });

    it('should handle API error scenarios gracefully', () => {
      // Simulate what happens when API returns error
      const result = V0DataAdapter.safeAdaptEmailMetrics(undefined, undefined as any);
      
      expect(result).toBeNull();
    });

    it('should handle partial API failures', () => {
      // System metrics available but email statuses failed to load
      const systemMetrics: SystemMetrics = {
        total_inquiries: 15,
        reports_generated: 8,
        emails_sent: 30,
        email_delivery_rate: 86.7,
        avg_report_gen_time_ms: 2100,
        system_uptime: '2d 4h',
      };

      const result = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, []);
      
      expect(result).not.toBeNull();
      expect(result!.totalEmails).toBe(30);
      expect(result!.deliveryRate).toBe(86.7);
      expect(result!.delivered).toBe(26); // 86.7% of 30
      expect(result!.failedEmails).toBe(4);
    });
  });
});