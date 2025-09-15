import React, { useState, useEffect } from 'react';
import { V0EmailDeliveryDashboard, EmailMetrics } from './V0EmailDeliveryDashboard';
import { V0DataAdapter } from './V0DataAdapter';
import apiService, { SystemMetrics, EmailStatus, Inquiry } from '../../services/api';
import { AlertTriangle } from 'lucide-react';
import { Card, CardContent } from '../ui/Card';

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

      // Try to fetch real email metrics first
      let systemMetrics: SystemMetrics | null = null;
      let emailStatuses: EmailStatus[] = [];
      let hasRealData = false;

      try {
        // Fetch system metrics for overall email statistics
        const systemMetricsResponse = await apiService.getSystemMetrics();
        systemMetrics = systemMetricsResponse.data;
        
        // Check if system metrics contain real email data
        if (systemMetrics && (systemMetrics.emails_sent > 0 || systemMetrics.email_delivery_rate > 0)) {
          hasRealData = true;
        }
      } catch (systemMetricsError) {
        console.warn('Failed to fetch system metrics:', systemMetricsError);
      }

      try {
        // Try to fetch detailed email events/metrics
        const emailMetricsResponse = await apiService.getEmailMetrics({
          start: getDateFromTimeRange(timeRange),
          end: new Date().toISOString()
        });
        
        if (emailMetricsResponse.success && emailMetricsResponse.data) {
          // If we get email events data, we have real data
          hasRealData = true;
          // Convert email events to email statuses format if needed
          // This would depend on the actual API response structure
        }
      } catch (emailMetricsError) {
        console.warn('Failed to fetch email metrics:', emailMetricsError);
      }

      try {
        // Fetch recent inquiries to get individual email statuses as fallback
        const inquiriesResponse = await apiService.listInquiries({ 
          limit: 50, // Reduced limit for better performance
          ...(timeRange !== 'all' && {
            date_from: getDateFromTimeRange(timeRange)
          })
        });
        const inquiries: Inquiry[] = inquiriesResponse.data;

        // Fetch email status for each inquiry (only if we don't have real data yet)
        if (!hasRealData && inquiries.length > 0) {
          const emailStatusPromises = inquiries.slice(0, 20).map(async (inquiry) => {
            try {
              const statusResponse = await apiService.getEmailStatus(inquiry.id);
              return statusResponse.data;
            } catch (error) {
              return null;
            }
          });

          const emailStatusResults = await Promise.allSettled(emailStatusPromises);
          emailStatuses = emailStatusResults
            .filter((result): result is PromiseFulfilledResult<EmailStatus> => 
              result.status === 'fulfilled' && result.value !== null
            )
            .map(result => result.value);
            
          if (emailStatuses.length > 0) {
            hasRealData = true;
          }
        }
      } catch (inquiriesError) {
        console.warn('Failed to fetch inquiries for email status:', inquiriesError);
      }

      // Transform data if we have any real data
      if (V0DataAdapter.hasRealEmailData(systemMetrics, emailStatuses)) {
        const adaptedMetrics = V0DataAdapter.safeAdaptEmailMetrics(systemMetrics, emailStatuses);
        if (adaptedMetrics) {
          setEmailMetrics(adaptedMetrics);
          setError(null); // Clear any previous errors since we have real data
          return; // Successfully loaded real data
        }
      }

      // If we reach here, we don't have real data
      const errorMessage = V0DataAdapter.getEmailDataErrorMessage(systemMetrics, emailStatuses);
      const errorDetails = V0DataAdapter.getEmailDataErrorDetails(systemMetrics, emailStatuses);
      
      setError(errorMessage);
      setEmailMetrics(null); // Don't show mock data, show error state instead
      
      // Log detailed error information for debugging
      console.warn('Email metrics unavailable:', {
        errorMessage,
        errorDetails,
        hasSystemMetrics: !!systemMetrics,
        emailStatusCount: emailStatuses.length,
        systemMetricsData: systemMetrics
      });

    } catch (err) {
      console.error('Failed to fetch email metrics:', err);
      
      // Extract more specific error information
      let apiError = '';
      if (err instanceof Error) {
        apiError = err.message;
      } else if (typeof err === 'object' && err !== null && 'message' in err) {
        apiError = String(err.message);
      } else {
        apiError = 'Failed to load email metrics';
      }
      
      const errorMessage = V0DataAdapter.getEmailDataErrorMessage(null, [], apiError);
      const errorDetails = V0DataAdapter.getEmailDataErrorDetails(null, [], apiError);
      
      setError(errorMessage);
      setEmailMetrics(null); // Don't fallback to mock data
      
      // Log detailed error information for debugging
      console.error('Email metrics fetch failed:', {
        originalError: err,
        errorMessage,
        errorDetails,
        timeRange
      });
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

  // Show error state when no real data is available
  if (error && !emailMetrics) {
    // Get error details for better user experience
    const errorDetails = V0DataAdapter.getEmailDataErrorDetails(null, [], error);
    
    return (
      <div className={className}>
        <Card className="border-red-200 bg-red-50">
          <CardContent className="pt-6">
            <div className="flex items-center space-x-2 text-red-600">
              <AlertTriangle className="h-5 w-5" />
              <span className="font-medium">Email Metrics Unavailable</span>
            </div>
            <p className="mt-2 text-sm text-red-600">
              {error}
            </p>
            
            {/* Show helpful suggestions based on error type */}
            {errorDetails.suggestions.length > 0 && (
              <div className="mt-3">
                <p className="text-xs font-medium text-red-700 mb-1">Suggestions:</p>
                <ul className="text-xs text-red-600 space-y-1">
                  {errorDetails.suggestions.slice(0, 2).map((suggestion, index) => (
                    <li key={index} className="flex items-start">
                      <span className="mr-1">•</span>
                      <span>{suggestion}</span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
            
            <div className="mt-4 flex space-x-2">
              <button
                onClick={fetchEmailMetrics}
                className="px-3 py-1 text-xs bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors"
              >
                Retry
              </button>
              
              {/* Show different actions based on error category */}
              {errorDetails.category === 'configuration' && (
                <button
                  onClick={() => window.open('/health/email-monitoring', '_blank')}
                  className="px-3 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
                >
                  Check System Health
                </button>
              )}
              
              {errorDetails.severity === 'high' && (
                <span className="px-2 py-1 text-xs bg-red-200 text-red-800 rounded">
                  High Priority Issue
                </span>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Show partial data warning if we have some data but encountered errors
  const partialDataWarning = error && emailMetrics && (
    <Card className="border-amber-200 bg-amber-50 mb-6">
      <CardContent className="pt-6">
        <div className="flex items-center space-x-2 text-amber-600">
          <AlertTriangle className="h-5 w-5" />
          <span className="font-medium">Partial Data Available</span>
        </div>
        <p className="mt-2 text-sm text-amber-600">
          Some email metrics could not be loaded. The displayed data may be incomplete.
        </p>
        <p className="text-xs text-amber-500 mt-1">
          {error}
        </p>
      </CardContent>
    </Card>
  );

  return (
    <div className={className}>
      {partialDataWarning}
      <V0EmailDeliveryDashboard
        metrics={emailMetrics || {
          deliveryRate: 0,
          openRate: 0,
          clickRate: 0,
          failedEmails: 0,
          totalEmails: 0,
          bounced: 0,
          spam: 0,
          delivered: 0,
          opened: 0,
          clicked: 0,
        }}
        timeRange={timeRange}
        onTimeRangeChange={handleTimeRangeChange}
      />
    </div>
  );
};

export default V0EmailDeliveryDashboardConnected;