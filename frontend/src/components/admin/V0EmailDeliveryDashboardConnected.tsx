import React, { useState, useEffect } from 'react';
import { V0EmailDeliveryDashboard, EmailMetrics } from './V0EmailDeliveryDashboard';
import { V0DataAdapter } from './V0DataAdapter';
import apiService, { SystemMetrics, EmailStatus, Inquiry } from '../../services/api';
import { AlertTriangle } from 'lucide-react';
import { Card, CardContent } from '../ui/card';

interface V0EmailDeliveryDashboardConnectedProps {
  className?: string;
}

/**
 * V0EmailDeliveryDashboardConnected - Connected version of email delivery dashboard
 * Integrates with backend APIs and uses V0DataAdapter for data transformation
 */
export const V0EmailDeliveryDashboardConnected: React.FC<V0EmailDeliveryDashboardConnectedProps> = ({
  className = '',
}) => {
  const [emailMetrics, setEmailMetrics] = useState<EmailMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState('24h');

  // Fetch email metrics data from backend
  const fetchEmailMetrics = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch system metrics for overall email statistics
      const systemMetricsResponse = await apiService.getSystemMetrics();
      const systemMetrics: SystemMetrics = systemMetricsResponse.data;

      // Fetch recent inquiries to get individual email statuses
      const inquiriesResponse = await apiService.listInquiries({ 
        limit: 100,
        // Filter by time range if needed
        ...(timeRange !== 'all' && {
          date_from: getDateFromTimeRange(timeRange)
        })
      });
      const inquiries: Inquiry[] = inquiriesResponse.data;

      // Fetch email status for each inquiry
      const emailStatusPromises = inquiries.map(async (inquiry) => {
        try {
          const statusResponse = await apiService.getEmailStatus(inquiry.id);
          return statusResponse.data;
        } catch (error) {
          // If email status is not available for this inquiry, return null
          console.warn(`Email status not available for inquiry ${inquiry.id}:`, error);
          return null;
        }
      });

      const emailStatusResults = await Promise.allSettled(emailStatusPromises);
      const emailStatuses: EmailStatus[] = emailStatusResults
        .filter((result): result is PromiseFulfilledResult<EmailStatus> => 
          result.status === 'fulfilled' && result.value !== null
        )
        .map(result => result.value);

      // Transform data using V0DataAdapter
      const adaptedMetrics = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, emailStatuses);
      setEmailMetrics(adaptedMetrics);

    } catch (err) {
      console.error('Failed to fetch email metrics:', err);
      setError(err instanceof Error ? err.message : 'Failed to load email metrics');
      
      // Generate mock data as fallback for demo purposes
      const mockMetrics = V0DataAdapter.generateMockEmailMetrics(30);
      setEmailMetrics(mockMetrics);
    } finally {
      setLoading(false);
    }
  };

  // Convert time range to date string for API filtering
  const getDateFromTimeRange = (range: string): string => {
    const now = new Date();
    let daysAgo = 1;

    switch (range) {
      case '1h':
        return new Date(now.getTime() - 60 * 60 * 1000).toISOString();
      case '24h':
        daysAgo = 1;
        break;
      case '7d':
        daysAgo = 7;
        break;
      case '30d':
        daysAgo = 30;
        break;
      case '90d':
        daysAgo = 90;
        break;
      default:
        daysAgo = 1;
    }

    return new Date(now.getTime() - daysAgo * 24 * 60 * 60 * 1000).toISOString();
  };

  // Handle time range change
  const handleTimeRangeChange = (newRange: string) => {
    setTimeRange(newRange);
  };

  // Fetch data on component mount and when time range changes
  useEffect(() => {
    fetchEmailMetrics();
  }, [timeRange]);

  // Auto-refresh data every 5 minutes
  useEffect(() => {
    const interval = setInterval(() => {
      fetchEmailMetrics();
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(interval);
  }, [timeRange]);

  // Loading state
  if (loading && !emailMetrics) {
    return (
      <div className={`space-y-6 ${className}`}>
        <div className="flex items-center justify-between">
          <div>
            <div className="h-8 w-64 bg-gray-200 rounded animate-pulse"></div>
            <div className="h-4 w-96 bg-gray-100 rounded animate-pulse mt-2"></div>
          </div>
          <div className="h-10 w-40 bg-gray-200 rounded animate-pulse"></div>
        </div>
        
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <Card key={i} className="animate-pulse">
              <CardContent className="pt-6">
                <div className="h-4 w-20 bg-gray-200 rounded mb-2"></div>
                <div className="h-8 w-16 bg-gray-200 rounded mb-2"></div>
                <div className="h-3 w-24 bg-gray-100 rounded mb-3"></div>
                <div className="h-1 w-full bg-gray-100 rounded"></div>
              </CardContent>
            </Card>
          ))}
        </div>
        
        <Card className="animate-pulse">
          <CardContent className="pt-6">
            <div className="h-6 w-48 bg-gray-200 rounded mb-2"></div>
            <div className="h-4 w-96 bg-gray-100 rounded mb-6"></div>
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <div key={i} className="space-y-2">
                  <div className="flex justify-between">
                    <div className="h-4 w-20 bg-gray-200 rounded"></div>
                    <div className="h-4 w-16 bg-gray-100 rounded"></div>
                  </div>
                  <div className="h-2 w-full bg-gray-100 rounded"></div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Error state (but still show dashboard with mock data)
  const errorBanner = error && (
    <Card className="border-amber-200 bg-amber-50 mb-6">
      <CardContent className="pt-6">
        <div className="flex items-center space-x-2 text-amber-600">
          <AlertTriangle className="h-5 w-5" />
          <span className="font-medium">Using demo data</span>
        </div>
        <p className="mt-2 text-sm text-amber-600">
          Unable to load live email metrics. Displaying sample data for demonstration.
        </p>
        <p className="text-xs text-amber-500 mt-1">
          Error: {error}
        </p>
      </CardContent>
    </Card>
  );

  return (
    <div className={className}>
      {errorBanner}
      <V0EmailDeliveryDashboard
        metrics={emailMetrics || V0DataAdapter.generateMockEmailMetrics(50)}
        timeRange={timeRange}
        onTimeRangeChange={handleTimeRangeChange}
      />
    </div>
  );
};

export default V0EmailDeliveryDashboardConnected;