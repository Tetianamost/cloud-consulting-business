import { V0DataAdapter } from './V0DataAdapter';
import { SystemMetrics } from '../../services/api';

describe('V0DataAdapter', () => {
  const mockMetrics: SystemMetrics = {
    total_inquiries: 15,
    reports_generated: 8,
    emails_sent: 30,
    email_delivery_rate: 94.5,
    avg_report_gen_time_ms: 1250,
    system_uptime: '3d 7h 22m',
    last_processed_at: new Date().toISOString(),
  };

  describe('adaptSystemMetrics', () => {
    it('should transform backend metrics to V0 format correctly', () => {
      const result = V0DataAdapter.adaptSystemMetrics(mockMetrics);
      
      expect(result).toHaveLength(4);
      
      // Check AI Reports metric
      expect(result[0].title).toBe('AI Reports Generated');
      expect(result[0].value).toBe(8);
      expect(result[0].trend).toBe('up');
      
      // Check Confidence Score metric
      expect(result[1].title).toBe('Avg Confidence Score');
      expect(result[1].value).toMatch(/\d+\.\d%/);
      expect(result[1].trend).toBe('up');
      
      // Check Processing Time metric
      expect(result[2].title).toBe('Avg Processing Time');
      expect(result[2].value).toBe('1.3s');
      expect(result[2].trend).toBe('up');
      
      // Check Opportunities metric
      expect(result[3].title).toBe('High-Value Opportunities');
      expect(result[3].value).toBe(4); // 30% of 15 inquiries
      expect(result[3].trend).toBe('up');
    });

    it('should handle null metrics gracefully', () => {
      const result = V0DataAdapter.adaptSystemMetrics(null);
      
      expect(result).toHaveLength(4);
      
      result.forEach(metric => {
        expect(metric.value).toBeDefined();
        expect(metric.change).toBeDefined();
        expect(metric.trend).toBeDefined();
        expect(metric.icon).toBeDefined();
      });
    });

    it('should handle undefined metrics gracefully', () => {
      const result = V0DataAdapter.adaptSystemMetrics(undefined as any);
      
      expect(result).toHaveLength(4);
      
      result.forEach(metric => {
        expect(metric.value).toBeDefined();
        expect(metric.change).toBeDefined();
        expect(metric.trend).toBeDefined();
        expect(metric.icon).toBeDefined();
      });
    });
  });

  describe('safeAdaptSystemMetrics', () => {
    it('should handle errors gracefully and return fallback metrics', () => {
      // Mock console.error to avoid noise in tests
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
      
      // Force an error by passing invalid data
      const result = V0DataAdapter.safeAdaptSystemMetrics(null);
      
      expect(result).toHaveLength(4);
      
      result.forEach(metric => {
        expect(metric.title).toBeDefined();
        expect(metric.value).toBeDefined();
        expect(metric.change).toBeDefined();
        expect(metric.trend).toBe('neutral');
        expect(metric.icon).toBeDefined();
      });
      
      consoleSpy.mockRestore();
    });

    it('should work with valid data', () => {
      const result = V0DataAdapter.safeAdaptSystemMetrics(mockMetrics);
      
      expect(result).toHaveLength(4);
      expect(result[0].value).toBe(8);
    });
  });

  describe('helper methods', () => {
    it('should format processing time correctly', () => {
      // Test seconds
      const result1 = V0DataAdapter.adaptSystemMetrics({
        ...mockMetrics,
        avg_report_gen_time_ms: 500
      });
      expect(result1[2].value).toBe('1s');

      // Test minutes
      const result2 = V0DataAdapter.adaptSystemMetrics({
        ...mockMetrics,
        avg_report_gen_time_ms: 90000 // 1.5 minutes
      });
      expect(result2[2].value).toBe('1.5min');

      // Test hours
      const result3 = V0DataAdapter.adaptSystemMetrics({
        ...mockMetrics,
        avg_report_gen_time_ms: 3900000 // 65 minutes
      });
      expect(result3[2].value).toBe('1h 5m');
    });

    it('should calculate confidence score correctly', () => {
      const result = V0DataAdapter.adaptSystemMetrics({
        ...mockMetrics,
        email_delivery_rate: 90,
        reports_generated: 12
      });
      
      // Should be around 90 * 0.89 + 2 = 82.1%
      expect(result[1].value).toMatch(/8[0-9]\.\d%/);
    });

    it('should calculate high-value opportunities correctly', () => {
      const result = V0DataAdapter.adaptSystemMetrics({
        ...mockMetrics,
        total_inquiries: 20
      });
      
      // Should be 30% of 20 = 6
      expect(result[3].value).toBe(6);
    });
  });

  describe('testDataTransformation', () => {
    it('should run without errors', () => {
      const consoleSpy = jest.spyOn(console, 'log').mockImplementation();
      
      expect(() => {
        V0DataAdapter.testDataTransformation();
      }).not.toThrow();
      
      expect(consoleSpy).toHaveBeenCalled();
      consoleSpy.mockRestore();
    });
  });
});