import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import V0MetricsCards from './V0MetricsCards';
import { V0DataAdapter } from './V0DataAdapter';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';
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
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  useEffect(() => {
    fetchDashboardData();
    const intervalId = setInterval(() => {
      fetchDashboardData();
    }, 30000);
    return () => clearInterval(intervalId);
  }, []);

  const refreshMetrics = async () => {
    try {
      setMetricsLoading(true);
      const metricsResponse = await apiService.getSystemMetrics();
      if (metricsResponse.success && metricsResponse.data) {
        setMetrics(metricsResponse.data);
        setLastUpdated(new Date());
      }
    } catch (err: any) {
      console.error('Error refreshing metrics:', err);
    } finally {
      setMetricsLoading(false);
    }
  };

  const fetchDashboardData = async () => {
    try {
      if (!metrics) {
        setLoading(true);
      }
      setError(null);
      
      const [metricsResponse, inquiriesResponse] = await Promise.all([
        apiService.getSystemMetrics().catch(err => ({ success: false, data: null, error: err.message })),
        apiService.listInquiries({ limit: 10 }).catch(err => ({ success: false, data: [], error: err.message }))
      ]);
      
      if (metricsResponse.success && metricsResponse.data) {
        setMetrics(metricsResponse.data);
        setLastUpdated(new Date());
      }
      
      if (inquiriesResponse.success && inquiriesResponse.data) {
        setInquiries(inquiriesResponse.data);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const renderOverviewContent = () => {
    const v0Metrics = V0DataAdapter.safeAdaptSystemMetrics(metrics);
    return (
      <div className="space-y-6">
        <div className="relative">
          <V0MetricsCards metrics={v0Metrics} loading={loading && !metrics} />
          {metricsLoading && (
            <div className="absolute top-0 right-0 bg-white border border-gray-200 rounded-md px-3 py-1 shadow-sm">
              <div className="flex items-center text-sm text-gray-600">
                <span className="animate-spin mr-2">🔄</span>
                Updating...
              </div>
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
    <>
      <div className="flex items-center justify-between mb-6">
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
          <button
            onClick={refreshMetrics}
            disabled={metricsLoading}
            className="bg-white border border-gray-300 text-gray-700 px-3 py-2 rounded-md hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            <span className={`mr-2 ${metricsLoading ? 'animate-spin' : ''}`}>📊</span>
            {metricsLoading ? 'Updating...' : 'Refresh Metrics'}
          </button>
        </div>
      </div>
      {loading && !metrics && (
        <div className="flex justify-center items-center py-8 text-gray-600">
          Loading dashboard data...
        </div>
      )}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <p className="text-red-800">{error}</p>
        </div>
      )}
      {!loading && renderContent()}
    </>
  );
};

export default V0DashboardNew;
