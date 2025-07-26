import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { AdminSidebar } from './sidebar';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';
import { InquiryAnalysisDashboard } from './inquiry-analysis-dashboard';

interface IntegratedAdminDashboardProps {
  children?: React.ReactNode;
}

export const IntegratedAdminDashboard: React.FC<IntegratedAdminDashboardProps> = ({ children }) => {
  return (
    <div className="flex min-h-screen bg-gray-100">
      <AdminSidebar />
      <main className="flex-1 p-6 overflow-y-auto">
        {children || (
          <Routes>
            <Route index element={<Navigate to="dashboard" replace />} />
            <Route path="dashboard" element={<InquiryAnalysisDashboard />} />
            <Route path="inquiries" element={<InquiryList />} />
            <Route path="metrics" element={<MetricsDashboard />} />
            <Route path="email-status" element={<EmailMonitor />} />
          </Routes>
        )}
      </main>
    </div>
  );
};

export default IntegratedAdminDashboard;