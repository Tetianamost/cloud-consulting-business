import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { FiBarChart2, FiPieChart, FiTrendingUp, FiClock } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import apiService, { SystemMetrics } from '../../services/api';
import Icon from '../ui/Icon';

const MetricsContainer = styled.div`
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

const MetricsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: ${theme.space[6]};
`;

const MetricCard = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  box-shadow: ${theme.shadows.md};
`;

const MetricHeader = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: ${theme.space[4]};
`;

const MetricIcon = styled.div<{ $bgColor: string }>`
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

const MetricTitle = styled.h3`
  font-size: ${theme.fontSizes.lg};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray700};
  margin: 0;
`;

const MetricValue = styled.div`
  font-size: ${theme.fontSizes['4xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-bottom: ${theme.space[2]};
`;

const MetricSubtext = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
`;

const ChartContainer = styled.div`
  height: 200px;
  margin-top: ${theme.space[4]};
  display: flex;
  align-items: center;
  justify-content: center;
`;

// Simple bar chart component
const BarChart = styled.div`
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 100%;
  width: 100%;
`;

const Bar = styled.div<{ $height: string; $color: string }>`
  width: 30px;
  height: ${props => props.$height};
  background-color: ${props => props.$color};
  border-radius: ${theme.borderRadius.sm} ${theme.borderRadius.sm} 0 0;
  transition: height 0.3s ease;
`;

// Simple pie chart component
const PieChart = styled.div<{ $percentage: number }>`
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background: conic-gradient(
    ${theme.colors.primary} ${props => props.$percentage}%,
    ${theme.colors.gray200} ${props => props.$percentage}% 100%
  );
  position: relative;
  
  &::after {
    content: '${props => props.$percentage}%';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: ${theme.fontSizes.xl};
    font-weight: ${theme.fontWeights.bold};
  }
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

const RefreshButton = styled.button`
  padding: ${theme.space[2]} ${theme.space[4]};
  background-color: ${theme.colors.primary};
  color: ${theme.colors.white};
  border: none;
  border-radius: ${theme.borderRadius.md};
  font-weight: ${theme.fontWeights.medium};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.primary};
    opacity: 0.9;
  }
`;

const MetricsDashboard: React.FC = () => {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  const fetchMetrics = async () => {
    try {
      setLoading(true);
      const response = await apiService.getSystemMetrics();
      setMetrics(response.data);
      setLoading(false);
    } catch (err) {
      console.error('Failed to fetch metrics:', err);
      setError('Failed to load metrics data. Please try again later.');
      setLoading(false);
    }
  };
  
  useEffect(() => {
    fetchMetrics();
    
    // Refresh metrics every 30 seconds
    const intervalId = setInterval(fetchMetrics, 30000);
    
    return () => clearInterval(intervalId);
  }, []);
  
  if (loading && !metrics) {
    return <LoadingState>Loading metrics data...</LoadingState>;
  }
  
  if (error) {
    return (
      <ErrorState>
        {error}
        <RefreshButton onClick={fetchMetrics} style={{ marginTop: theme.space[4] }}>
          Try Again
        </RefreshButton>
      </ErrorState>
    );
  }
  
  return (
    <MetricsContainer>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <PageTitle>System Metrics</PageTitle>
        <RefreshButton onClick={fetchMetrics}>Refresh Metrics</RefreshButton>
      </div>
      
      <MetricsGrid>
        <MetricCard>
          <MetricHeader>
            <MetricIcon $bgColor={theme.colors.accent}>
              <Icon icon={FiBarChart2} size={24} />
            </MetricIcon>
            <MetricTitle>Inquiry Statistics</MetricTitle>
          </MetricHeader>
          <MetricValue>{metrics?.total_inquiries || 0}</MetricValue>
          <MetricSubtext>Total inquiries processed</MetricSubtext>
          
          <ChartContainer>
            <BarChart>
              {/* Simulated data for visualization */}
              <Bar $height="40%" $color={theme.colors.accent} />
              <Bar $height="65%" $color={theme.colors.accent} />
              <Bar $height="85%" $color={theme.colors.accent} />
              <Bar $height="70%" $color={theme.colors.accent} />
              <Bar $height="90%" $color={theme.colors.accent} />
              <Bar $height="60%" $color={theme.colors.accent} />
              <Bar $height="75%" $color={theme.colors.accent} />
            </BarChart>
          </ChartContainer>
        </MetricCard>
        
        <MetricCard>
          <MetricHeader>
            <MetricIcon $bgColor={theme.colors.success}>
              <Icon icon={FiPieChart} size={24} />
            </MetricIcon>
            <MetricTitle>Email Delivery Rate</MetricTitle>
          </MetricHeader>
          <MetricValue>{metrics?.email_delivery_rate.toFixed(1)}%</MetricValue>
          <MetricSubtext>{metrics?.emails_sent || 0} emails sent successfully</MetricSubtext>
          
          <ChartContainer>
            <PieChart $percentage={metrics?.email_delivery_rate || 0} />
          </ChartContainer>
        </MetricCard>
        
        <MetricCard>
          <MetricHeader>
            <MetricIcon $bgColor={theme.colors.secondary}>
              <Icon icon={FiTrendingUp} size={24} />
            </MetricIcon>
            <MetricTitle>Report Generation</MetricTitle>
          </MetricHeader>
          <MetricValue>{metrics?.reports_generated || 0}</MetricValue>
          <MetricSubtext>Average generation time: {metrics?.avg_report_gen_time_ms.toFixed(0)}ms</MetricSubtext>
          
          <ChartContainer>
            <BarChart>
              {/* Simulated data for visualization */}
              <Bar $height="60%" $color={theme.colors.secondary} />
              <Bar $height="75%" $color={theme.colors.secondary} />
              <Bar $height="45%" $color={theme.colors.secondary} />
              <Bar $height="80%" $color={theme.colors.secondary} />
              <Bar $height="65%" $color={theme.colors.secondary} />
              <Bar $height="90%" $color={theme.colors.secondary} />
              <Bar $height="70%" $color={theme.colors.secondary} />
            </BarChart>
          </ChartContainer>
        </MetricCard>
        
        <MetricCard>
          <MetricHeader>
            <MetricIcon $bgColor={theme.colors.warning}>
              <Icon icon={FiClock} size={24} />
            </MetricIcon>
            <MetricTitle>System Performance</MetricTitle>
          </MetricHeader>
          <MetricValue>{metrics?.system_uptime || '0h'}</MetricValue>
          <MetricSubtext>
            Last activity: {metrics?.last_processed_at 
              ? new Date(metrics.last_processed_at).toLocaleString() 
              : 'No recent activity'}
          </MetricSubtext>
          
          <div style={{ marginTop: theme.space[4] }}>
            <div style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center',
              marginBottom: theme.space[2]
            }}>
              <span>CPU Usage</span>
              <span>45%</span>
            </div>
            <div style={{ 
              height: '8px', 
              background: theme.colors.gray200, 
              borderRadius: theme.borderRadius.full,
              overflow: 'hidden'
            }}>
              <div style={{ 
                width: '45%', 
                height: '100%', 
                background: theme.colors.warning,
                borderRadius: theme.borderRadius.full
              }} />
            </div>
            
            <div style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center',
              marginTop: theme.space[3],
              marginBottom: theme.space[2]
            }}>
              <span>Memory Usage</span>
              <span>68%</span>
            </div>
            <div style={{ 
              height: '8px', 
              background: theme.colors.gray200, 
              borderRadius: theme.borderRadius.full,
              overflow: 'hidden'
            }}>
              <div style={{ 
                width: '68%', 
                height: '100%', 
                background: theme.colors.warning,
                borderRadius: theme.borderRadius.full
              }} />
            </div>
          </div>
        </MetricCard>
      </MetricsGrid>
    </MetricsContainer>
  );
};

export default MetricsDashboard;