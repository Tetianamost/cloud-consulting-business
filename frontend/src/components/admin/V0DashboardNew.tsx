import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import V0MetricsCards from './V0MetricsCards';
import { V0DataAdapter } from './V0DataAdapter';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';
import { 
  V0SkeletonDashboard, 
  V0InlineLoader, 
  V0RefreshButton,
  V0LoadingSpinner 
} from './V0LoadingStates';
import { V0ApiErrorFallback } from './V0ErrorFallbacks';
import { useV0ApiErrorHandler } from './useV0ErrorHandler';
import apiService, { Inquiry, SystemMetrics } from '../../services/api';

interface V0DashboardNewProps {
  children?: React.ReactNode;
}

const V0DashboardNew: React.FC<V0DashboardNewProps> = ({ children }) => {
  const location = useLocation();
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState(true);
  const [metricsLoading, setMetricsLoading] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  
  const { handleApiCall, isError, error, retry, clearError } = useV0ApiErrorHandler({
    maxRetries: 3,
    onError: (error) => console.error('Dashboard API error:', error)
  });

  useEffect(() => {
    fetchDashboardData();
    const intervalId = setInterval(() => {
      fetchDashboardData();
    }, 30000);
    return () => clearInterval(intervalId);
  }, []);

  const refreshMetrics = async () => {
    setMetricsLoading(true);
    const result = await handleApiCall(
      () => apiService.getSystemMetrics(),
      'Failed to refresh metrics'
    );
    
    if (result) {
      setMetrics(result);
      setLastUpdated(new Date());
    }
    setMetricsLoading(false);
  };

  const fetchDashboardData = async () => {
    if (!metrics) {
      setLoading(true);
    }
    clearError();
    
    const [metricsResult, inquiriesResult] = await Promise.all([
      handleApiCall(
        () => apiService.getSystemMetrics(),
        'Failed to load metrics'
      ),
      handleApiCall(
        () => apiService.listInquiries({ limit: 10 }),
        'Failed to load inquiries'
      )
    ]);
    
    if (metricsResult) {
      setMetrics(metricsResult);
      setLastUpdated(new Date());
    }
    
    if (inquiriesResult) {
      setInquiries(inquiriesResult);
    }
    
    setLoading(false);
  };

  const renderOverviewContent = () => {
    const v0Metrics = V0DataAdapter.safeAdaptSystemMetrics(metrics);
    return (
      <div className="space-y-6">
        <div className="relative bg-white rounded-lg p-6 shadow-sm border border-gray-200 hover:shadow-md transition-shadow duration-200">
          <V0MetricsCards metrics={v0Metrics} loading={loading && !metrics} />
          {metricsLoading && (
            <div className="absolute top-0 right-0 bg-blue-50 border border-blue-200 rounded-md px-3 py-1 shadow-sm">
              <V0InlineLoader message="Updating..." />
            </div>
          )}
        </div>
      </div>
    );
  };

  const renderContent = () => {
    if (children) return children;
    switch (location.pathname) {
      case '/admin/inquiries': return <InquiryList />;
      case '/admin/metrics': return <MetricsDashboard />;
      case '/admin/email-status': return <EmailMonitor />;
      default: return renderOverviewContent();
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between mb-6 bg-white rounded-lg p-6 shadow-sm border border-gray-200">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">AI Inquiry Analysis Dashboard</h1>
          <div className="flex items-center space-x-4 mt-1">
            <p className="text-gray-600">Real-time insights and analytics</p>
            {lastUpdated && (
              <span className="text-sm text-gray-500">
                Last updated: {lastUpdated.toLocaleTimeString()}
              </span>
            )}
          </div>
        </div>
        <div className="flex space-x-2">
          <V0RefreshButton
            onClick={refreshMetrics}
            loading={metricsLoading}
            size="md"
          >
            Refresh Metrics
          </V0RefreshButton>
        </div>
      </div>
      {loading && !metrics && <V0SkeletonDashboard />}
      {isError && (
        <V0ApiErrorFallback 
          message={error?.message || 'Failed to load dashboard data'}
          onRetry={() => retry(fetchDashboardData)}
        />
      )}
      {!loading && !isError && renderContent()}
    </div>
  );
};

export default V0DashboardNew;
