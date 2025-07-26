import React, { useState, useEffect } from 'react';
import { useLocation, Link } from 'react-router-dom';
import styled from 'styled-components';
import {
  FiBarChart,
  FiMail,
  FiMessageSquare,
  FiSettings,
  FiUsers,
  FiArrowUp,
  FiArrowDown,
  FiSearch,
  FiFilter,
  FiDownload,
  FiMoreHorizontal,
  FiCheckCircle,
  FiXCircle,
  FiClock,
  FiAlertCircle
} from 'react-icons/fi';
import { theme } from '../../styles/theme';
import apiService, { Inquiry, SystemMetrics, EmailStatus } from '../../services/api';
import Icon from '../ui/Icon';
import { AdminSidebar } from './sidebar';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';

// Styled Components
const DashboardContainer = styled.div`
  display: flex;
  min-height: 100vh;
  background-color: ${theme.colors.gray100};
`;

const Sidebar = styled.aside`
  width: 250px;
  background-color: ${theme.colors.white};
  border-right: 1px solid ${theme.colors.gray200};
  display: flex;
  flex-direction: column;
`;

const SidebarHeader = styled.div`
  padding: ${theme.space[6]};
  border-bottom: 1px solid ${theme.colors.gray200};
`;

const SidebarTitle = styled.h1`
  font-size: ${theme.fontSizes.xl};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  display: flex;
  align-items: center;
  gap: ${theme.space[2]};
`;

const SidebarNav = styled.nav`
  flex: 1;
  padding: ${theme.space[4]};
`;

const NavItem = styled(Link) <{ $active?: boolean }>`
  display: flex;
  align-items: center;
  gap: ${theme.space[3]};
  padding: ${theme.space[3]};
  margin-bottom: ${theme.space[2]};
  border-radius: ${theme.borderRadius.md};
  color: ${props => props.$active ? theme.colors.primary : theme.colors.gray700};
  background-color: ${props => props.$active ? `${theme.colors.primary}10` : 'transparent'};
  text-decoration: none;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${props => props.$active ? `${theme.colors.primary}15` : theme.colors.gray100};
    color: ${props => props.$active ? theme.colors.primary : theme.colors.gray900};
  }
`;

const MainContent = styled.main`
  flex: 1;
  padding: ${theme.space[6]};
  overflow-y: auto;
`;

const PageHeader = styled.div`
  margin-bottom: ${theme.space[6]};
`;

const PageTitle = styled.h1`
  font-size: ${theme.fontSizes['2xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-bottom: ${theme.space[2]};
`;

const PageDescription = styled.p`
  color: ${theme.colors.gray600};
  font-size: ${theme.fontSizes.md};
`;

const MetricsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: ${theme.space[4]};
  margin-bottom: ${theme.space[6]};
`;

const MetricCard = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[6]};
  box-shadow: ${theme.shadows.sm};
  border: 1px solid ${theme.colors.gray200};
`;

const MetricHeader = styled.div`
  display: flex;
  justify-content: between;
  align-items: center;
  margin-bottom: ${theme.space[4]};
`;

const MetricTitle = styled.h3`
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray600};
  margin: 0;
`;

const MetricValue = styled.div`
  font-size: ${theme.fontSizes['2xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-bottom: ${theme.space[2]};
`;

const MetricChange = styled.div<{ $positive?: boolean }>`
  display: flex;
  align-items: center;
  gap: ${theme.space[1]};
  font-size: ${theme.fontSizes.xs};
  color: ${props => props.$positive ? theme.colors.success : theme.colors.danger};
`;

const TabsContainer = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  box-shadow: ${theme.shadows.sm};
  border: 1px solid ${theme.colors.gray200};
`;

const TabsList = styled.div`
  display: flex;
  border-bottom: 1px solid ${theme.colors.gray200};
  padding: ${theme.space[4]} ${theme.space[6]} 0;
`;

const TabButton = styled.button<{ $active?: boolean }>`
  padding: ${theme.space[3]} ${theme.space[4]};
  border: none;
  background: none;
  color: ${props => props.$active ? theme.colors.primary : theme.colors.gray600};
  font-weight: ${props => props.$active ? theme.fontWeights.medium : theme.fontWeights.regular};
  border-bottom: 2px solid ${props => props.$active ? theme.colors.primary : 'transparent'};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    color: ${theme.colors.primary};
  }
`;

const TabContent = styled.div`
  padding: ${theme.space[6]};
`;

const LoadingState = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: ${theme.space[8]};
  color: ${theme.colors.gray600};
`;

const ErrorState = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: ${theme.space[8]};
  color: ${theme.colors.danger};
  text-align: center;
`;

// Navigation items
const navItems = [
  {
    title: "Dashboard",
    href: "/admin/dashboard",
    icon: FiBarChart,
  },
  {
    title: "Inquiries",
    href: "/admin/inquiries",
    icon: FiMessageSquare,
  },
  {
    title: "Email Status",
    href: "/admin/email-status",
    icon: FiMail,
  },
  {
    title: "Metrics",
    href: "/admin/metrics",
    icon: FiBarChart,
  },
];

interface V0DashboardProps {
  children?: React.ReactNode;
}

const V0Dashboard: React.FC<V0DashboardProps> = ({ children }) => {
  const location = useLocation();
  const [activeTab, setActiveTab] = useState('overview');
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);

      console.log('Fetching dashboard data...');

      // Fetch metrics and inquiries in parallel
      const [metricsResponse, inquiriesResponse] = await Promise.all([
        apiService.getSystemMetrics().catch(err => {
          console.error('Failed to fetch metrics:', err);
          return { success: false, data: null };
        }),
        apiService.listInquiries({ limit: 10 }).catch(err => {
          console.error('Failed to fetch inquiries:', err);
          return { success: false, data: [], count: 0, total: 0, page: 1, pages: 1 };
        })
      ]);

      console.log('Metrics response:', metricsResponse);
      console.log('Inquiries response:', inquiriesResponse);

      if (metricsResponse.success && metricsResponse.data) {
        setMetrics(metricsResponse.data);
      }

      if (inquiriesResponse.success && inquiriesResponse.data) {
        setInquiries(inquiriesResponse.data);
      } else {
        setInquiries([]);
      }
    } catch (err: any) {
      console.error('Failed to fetch dashboard data:', err);
      setError(err.message || 'Failed to load dashboard data');
      // Set default values to prevent null reference errors
      setMetrics(null);
      setInquiries([]);
    } finally {
      setLoading(false);
    }
  };

  const renderOverviewTab = () => (
    <div>
      <MetricsGrid>
        <MetricCard>
          <MetricHeader>
            <MetricTitle>Total Inquiries</MetricTitle>
            <Icon icon={FiMessageSquare} size={16} color={theme.colors.gray400} />
          </MetricHeader>
          <MetricValue>{metrics?.total_inquiries || 0}</MetricValue>
          <MetricChange $positive={true}>
            <Icon icon={FiArrowUp} size={12} />
            12% from last month
          </MetricChange>
        </MetricCard>

        <MetricCard>
          <MetricHeader>
            <MetricTitle>Emails Sent</MetricTitle>
            <Icon icon={FiMail} size={16} color={theme.colors.gray400} />
          </MetricHeader>
          <MetricValue>{metrics?.emails_sent || 0}</MetricValue>
          <MetricChange $positive={true}>
            <Icon icon={FiArrowUp} size={12} />
            8% from last month
          </MetricChange>
        </MetricCard>

        <MetricCard>
          <MetricHeader>
            <MetricTitle>Reports Generated</MetricTitle>
            <Icon icon={FiBarChart} size={16} color={theme.colors.gray400} />
          </MetricHeader>
          <MetricValue>{metrics?.reports_generated || 0}</MetricValue>
          <MetricChange $positive={true}>
            <Icon icon={FiArrowUp} size={12} />
            20% from last month
          </MetricChange>
        </MetricCard>

        <MetricCard>
          <MetricHeader>
            <MetricTitle>Email Delivery Rate</MetricTitle>
            <Icon icon={FiCheckCircle} size={16} color={theme.colors.gray400} />
          </MetricHeader>
          <MetricValue>{metrics?.email_delivery_rate?.toFixed(1) || 0}%</MetricValue>
          <MetricChange $positive={false}>
            <Icon icon={FiArrowDown} size={12} />
            3% from last month
          </MetricChange>
        </MetricCard>
      </MetricsGrid>

      {/* Recent Inquiries */}
      <div style={{ marginTop: theme.space[6] }}>
        <h3 style={{
          fontSize: theme.fontSizes.lg,
          fontWeight: theme.fontWeights.medium,
          marginBottom: theme.space[4],
          color: theme.colors.gray900
        }}>
          Recent Inquiries
        </h3>

        {inquiries && inquiries.length > 0 ? (
          <div style={{
            backgroundColor: theme.colors.white,
            borderRadius: theme.borderRadius.lg,
            border: `1px solid ${theme.colors.gray200}`,
            overflow: 'hidden'
          }}>
            {(inquiries || []).slice(0, 5).map((inquiry, index) => (
              <div key={inquiry.id} style={{
                padding: theme.space[4],
                borderBottom: index < 4 ? `1px solid ${theme.colors.gray200}` : 'none',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center'
              }}>
                <div>
                  <div style={{
                    fontWeight: theme.fontWeights.medium,
                    color: theme.colors.gray900,
                    marginBottom: theme.space[1]
                  }}>
                    {inquiry.name}
                  </div>
                  <div style={{
                    fontSize: theme.fontSizes.sm,
                    color: theme.colors.gray600
                  }}>
                    {inquiry.email} • {inquiry.services.join(', ')}
                  </div>
                </div>
                <div style={{
                  fontSize: theme.fontSizes.sm,
                  color: theme.colors.gray500
                }}>
                  {new Date(inquiry.created_at).toLocaleDateString()}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div style={{
            textAlign: 'center',
            padding: theme.space[8],
            color: theme.colors.gray500
          }}>
            No inquiries found
          </div>
        )}
      </div>
    </div>
  );

  // Route-based content rendering
  const renderContent = () => {
    switch (location.pathname) {
      case '/admin/inquiries':
        return <InquiryList />;
      case '/admin/metrics':
        return <MetricsDashboard />;
      case '/admin/email-status':
        return <EmailMonitor />;
      default:
        return renderOverviewTab();
    }
  };

  if (children) {
    return (
      <div className="admin-layout flex min-h-screen bg-gray-100">
        <AdminSidebar />
        <main className="flex-1 p-6 overflow-y-auto">
          {children}
        </main>
      </div>
    );
  }

  return (
    <div className="admin-layout flex min-h-screen bg-gray-100">
      <AdminSidebar />

      <main className="flex-1 p-6 overflow-y-auto">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            {location.pathname === '/admin/dashboard' ? 'Dashboard' :
              location.pathname === '/admin/inquiries' ? 'Inquiries' :
                location.pathname === '/admin/metrics' ? 'Metrics' :
                  location.pathname === '/admin/email-status' ? 'Email Status' : 'Admin Portal'}
          </h1>
          <p className="text-gray-600">
            {location.pathname === '/admin/dashboard' ? 'Overview of your cloud consulting business metrics and recent activity.' :
              location.pathname === '/admin/inquiries' ? 'Manage and track customer inquiries and reports.' :
                location.pathname === '/admin/metrics' ? 'System performance metrics and analytics.' :
                  location.pathname === '/admin/email-status' ? 'Monitor email delivery status and performance.' : ''}
          </p>
        </div>

        {loading && <div className="flex justify-center items-center py-8 text-gray-600">Loading dashboard data...</div>}

        {error && (
          <div className="flex flex-col items-center py-8 text-center">
            <p className="text-red-600 mb-4">Failed to load dashboard data</p>
            <p className="text-gray-600 mb-4">{error}</p>
            <button
              onClick={fetchDashboardData}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Try Again
            </button>
          </div>
        )}

        {!loading && renderContent()}
      </main>
    </div>
  );
};

export default V0Dashboard;