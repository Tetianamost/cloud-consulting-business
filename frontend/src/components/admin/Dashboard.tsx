import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { Link } from 'react-router-dom';
import { FiUsers, FiFileText, FiMail, FiClock } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import apiService, { SystemMetrics, Inquiry } from '../../services/api';
import Icon from '../ui/Icon';

const DashboardContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${theme.space[6]};
`;

const PageTitle = styled.h1`
  font-size: ${theme.fontSizes['3xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-bottom: ${theme.space[6]};
`;

const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: ${theme.space[4]};
`;

const StatCard = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  box-shadow: ${theme.shadows.md};
  display: flex;
  flex-direction: column;
`;

const StatHeader = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: ${theme.space[3]};
`;

const StatIcon = styled.div<{ $bgColor: string }>`
  width: 48px;
  height: 48px;
  border-radius: ${theme.borderRadius.full};
  background-color: ${props => props.$bgColor};
  display: flex;
  align-items: center;
  justify-content: center;
  color: ${theme.colors.white};
  margin-right: ${theme.space[3]};
`;

const StatTitle = styled.h3`
  font-size: ${theme.fontSizes.md};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray600};
  margin: 0;
`;

const StatValue = styled.div`
  font-size: ${theme.fontSizes['3xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-top: ${theme.space[2]};
`;

const RecentInquiriesSection = styled.div`
  margin-top: ${theme.space[6]};
`;

const SectionHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: ${theme.space[4]};
`;

const SectionTitle = styled.h2`
  font-size: ${theme.fontSizes['2xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin: 0;
`;

const ViewAllLink = styled(Link)`
  color: ${theme.colors.primary};
  font-weight: ${theme.fontWeights.medium};
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  box-shadow: ${theme.shadows.md};
  overflow: hidden;
`;

const TableHead = styled.thead`
  background-color: ${theme.colors.gray100};
`;

const TableRow = styled.tr`
  &:not(:last-child) {
    border-bottom: 1px solid ${theme.colors.gray200};
  }
  
  &:hover {
    background-color: ${theme.colors.gray100};
  }
`;

const TableHeader = styled.th`
  padding: ${theme.space[4]};
  text-align: left;
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray700};
`;

const TableCell = styled.td`
  padding: ${theme.space[4]};
  color: ${theme.colors.gray800};
`;

const StatusBadge = styled.span<{ $status: string }>`
  display: inline-block;
  padding: ${theme.space[1]} ${theme.space[3]};
  border-radius: ${theme.borderRadius.full};
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  
  ${props => {
    switch (props.$status) {
      case 'pending':
        return `
          background-color: ${theme.colors.warning};
          color: ${theme.colors.dark};
        `;
      case 'processing':
        return `
          background-color: ${theme.colors.info};
          color: ${theme.colors.white};
        `;
      case 'reviewed':
        return `
          background-color: ${theme.colors.success};
          color: ${theme.colors.white};
        `;
      case 'responded':
        return `
          background-color: ${theme.colors.accent};
          color: ${theme.colors.white};
        `;
      case 'closed':
        return `
          background-color: ${theme.colors.gray200};
          color: ${theme.colors.gray800};
        `;
      default:
        return `
          background-color: ${theme.colors.gray100};
          color: ${theme.colors.gray800};
        `;
    }
  }}
`;

const LoadingState = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: ${theme.space[6]};
  color: ${theme.colors.gray600};
  font-size: ${theme.fontSizes.lg};
`;

const ErrorState = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: ${theme.space[6]};
  color: ${theme.colors.danger};
  font-size: ${theme.fontSizes.lg};
`;

const Dashboard: React.FC = () => {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [recentInquiries, setRecentInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        
        // Fetch metrics
        const metricsResponse = await apiService.getSystemMetrics();
        setMetrics(metricsResponse.data);
        
        // Fetch recent inquiries (limit to 5)
        const inquiriesResponse = await apiService.listInquiries({ limit: 5 });
        setRecentInquiries(inquiriesResponse.data);
        
        setLoading(false);
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
        setError('Failed to load dashboard data. Please try again later.');
        setLoading(false);
      }
    };
    
    fetchDashboardData();
  }, []);
  
  if (loading) {
    return <LoadingState>Loading dashboard data...</LoadingState>;
  }
  
  if (error) {
    return <ErrorState>{error}</ErrorState>;
  }
  
  return (
    <DashboardContainer>
      <PageTitle>Admin Dashboard</PageTitle>
      
      <StatsGrid>
        <StatCard>
          <StatHeader>
            <StatIcon $bgColor={theme.colors.accent}>
              <Icon icon={FiUsers} size={24} />
            </StatIcon>
            <StatTitle>Total Inquiries</StatTitle>
          </StatHeader>
          <StatValue>{metrics?.total_inquiries || 0}</StatValue>
        </StatCard>
        
        <StatCard>
          <StatHeader>
            <StatIcon $bgColor={theme.colors.success}>
              <Icon icon={FiFileText} size={24} />
            </StatIcon>
            <StatTitle>Reports Generated</StatTitle>
          </StatHeader>
          <StatValue>{metrics?.reports_generated || 0}</StatValue>
        </StatCard>
        
        <StatCard>
          <StatHeader>
            <StatIcon $bgColor={theme.colors.secondary}>
              <Icon icon={FiMail} size={24} />
            </StatIcon>
            <StatTitle>Emails Sent</StatTitle>
          </StatHeader>
          <StatValue>{metrics?.emails_sent || 0}</StatValue>
        </StatCard>
        
        <StatCard>
          <StatHeader>
            <StatIcon $bgColor={theme.colors.warning}>
              <Icon icon={FiClock} size={24} />
            </StatIcon>
            <StatTitle>System Uptime</StatTitle>
          </StatHeader>
          <StatValue>{metrics?.system_uptime || '0h'}</StatValue>
        </StatCard>
      </StatsGrid>
      
      <RecentInquiriesSection>
        <SectionHeader>
          <SectionTitle>Recent Inquiries</SectionTitle>
          <ViewAllLink to="/admin/inquiries">View All</ViewAllLink>
        </SectionHeader>
        
        {recentInquiries.length > 0 ? (
          <Table>
            <TableHead>
              <TableRow>
                <TableHeader>Name</TableHeader>
                <TableHeader>Email</TableHeader>
                <TableHeader>Service</TableHeader>
                <TableHeader>Status</TableHeader>
                <TableHeader>Date</TableHeader>
              </TableRow>
            </TableHead>
            <tbody>
              {recentInquiries.map((inquiry) => (
                <TableRow key={inquiry.id}>
                  <TableCell>{inquiry.name}</TableCell>
                  <TableCell>{inquiry.email}</TableCell>
                  <TableCell>{inquiry.services.join(', ')}</TableCell>
                  <TableCell>
                    <StatusBadge $status={inquiry.status}>
                      {inquiry.status.charAt(0).toUpperCase() + inquiry.status.slice(1)}
                    </StatusBadge>
                  </TableCell>
                  <TableCell>
                    {new Date(inquiry.created_at).toLocaleDateString()}
                  </TableCell>
                </TableRow>
              ))}
            </tbody>
          </Table>
        ) : (
          <div>No recent inquiries found.</div>
        )}
      </RecentInquiriesSection>
    </DashboardContainer>
  );
};

export default Dashboard;