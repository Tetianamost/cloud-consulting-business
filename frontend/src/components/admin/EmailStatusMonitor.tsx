import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { FiMail, FiCheck, FiX, FiClock, FiRefreshCw } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import apiService, { EmailStatus, Inquiry } from '../../services/api';
import Icon from '../ui/icon';

const EmailStatusContainer = styled.div`
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

const EmailStatusCard = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  padding: ${theme.space[5]};
  box-shadow: ${theme.shadows.md};
  margin-bottom: ${theme.space[4]};
`;

const EmailStatusHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: ${theme.space[4]};
  padding-bottom: ${theme.space[4]};
  border-bottom: 1px solid ${theme.colors.gray200};
`;

const EmailStatusTitle = styled.h3`
  font-size: ${theme.fontSizes.xl};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray800};
  margin: 0;
  display: flex;
  align-items: center;
  
  svg {
    margin-right: ${theme.space[2]};
  }
`;

const EmailStatusBadge = styled.span<{ $status: string }>`
  display: inline-flex;
  align-items: center;
  padding: ${theme.space[1]} ${theme.space[3]};
  border-radius: ${theme.borderRadius.full};
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  
  svg {
    margin-right: ${theme.space[1]};
  }
  
  ${props => {
        switch (props.$status) {
            case 'delivered':
                return `
          background-color: ${theme.colors.success};
          color: ${theme.colors.white};
        `;
            case 'sending':
                return `
          background-color: ${theme.colors.info};
          color: ${theme.colors.white};
        `;
            case 'failed':
                return `
          background-color: ${theme.colors.danger};
          color: ${theme.colors.white};
        `;
            default:
                return `
          background-color: ${theme.colors.gray100};
          color: ${theme.colors.gray800};
        `;
        }
    }}
`;

const EmailStatusGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: ${theme.space[4]};
`;

const EmailStatusItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const EmailStatusLabel = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
  margin-bottom: ${theme.space[1]};
`;

const EmailStatusValue = styled.div`
  font-size: ${theme.fontSizes.md};
  color: ${theme.colors.gray900};
`;

const EmailStatusTimeline = styled.div`
  margin-top: ${theme.space[4]};
  padding-top: ${theme.space[4]};
  border-top: 1px solid ${theme.colors.gray200};
`;

const TimelineTitle = styled.h4`
  font-size: ${theme.fontSizes.md};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray800};
  margin-bottom: ${theme.space[3]};
`;

const TimelineItem = styled.div`
  display: flex;
  margin-bottom: ${theme.space[3]};
`;

const TimelineIcon = styled.div<{ $status: string }>`
  width: 24px;
  height: 24px;
  border-radius: ${theme.borderRadius.full};
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: ${theme.space[3]};
  
  ${props => {
        switch (props.$status) {
            case 'complete':
                return `
          background-color: ${theme.colors.success};
          color: ${theme.colors.white};
        `;
            case 'pending':
                return `
          background-color: ${theme.colors.info};
          color: ${theme.colors.white};
        `;
            case 'error':
                return `
          background-color: ${theme.colors.danger};
          color: ${theme.colors.white};
        `;
            default:
                return `
          background-color: ${theme.colors.gray500};
          color: ${theme.colors.white};
        `;
        }
    }}
`;

const TimelineContent = styled.div`
  flex: 1;
`;

const TimelineText = styled.div`
  font-size: ${theme.fontSizes.md};
  color: ${theme.colors.gray800};
`;

const TimelineTime = styled.div`
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray600};
`;

const InquirySelector = styled.div`
  margin-bottom: ${theme.space[4]};
  display: flex;
  gap: ${theme.space[4]};
`;

const InquirySelect = styled.select`
  padding: ${theme.space[3]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  background-color: ${theme.colors.white};
  flex: 1;
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
  }
`;

const RefreshButton = styled.button`
  padding: ${theme.space[3]};
  background-color: ${theme.colors.primary};
  color: ${theme.colors.white};
  border: none;
  border-radius: ${theme.borderRadius.md};
  font-weight: ${theme.fontWeights.medium};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  display: flex;
  align-items: center;
  justify-content: center;
  
  svg {
    margin-right: ${theme.space[2]};
  }
  
  &:hover {
    background-color: ${theme.colors.primary};
    opacity: 0.9;
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

const NoDataState = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: ${theme.space[6]};
  color: ${theme.colors.gray600};
  text-align: center;
`;

const IconWrapper = styled.div<{ $marginBottom?: string; $opacity?: number }>`
  margin-bottom: ${props => props.$marginBottom || '0'};
  opacity: ${props => props.$opacity || 1};
`;

const EmailStatusMonitor: React.FC = () => {
    const [inquiries, setInquiries] = useState<Inquiry[]>([]);
    const [selectedInquiryId, setSelectedInquiryId] = useState<string>('');
    const [emailStatus, setEmailStatus] = useState<EmailStatus | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    const fetchInquiries = async () => {
        try {
            const response = await apiService.listInquiries({ limit: 100 });
            setInquiries(response.data);

            // Select the first inquiry by default
            if (response.data.length > 0 && !selectedInquiryId) {
                setSelectedInquiryId(response.data[0].id);
            }

            setLoading(false);
        } catch (err) {
            console.error('Failed to fetch inquiries:', err);
            setError('Failed to load inquiries. Please try again later.');
            setLoading(false);
        }
    };

    const fetchEmailStatus = async (inquiryId: string) => {
        if (!inquiryId) return;

        try {
            setLoading(true);
            const response = await apiService.getEmailStatus(inquiryId);
            setEmailStatus(response.data);
            setLoading(false);
        } catch (err) {
            console.error('Failed to fetch email status:', err);
            setError('Failed to load email status. Please try again later.');
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchInquiries();
    }, []);

    useEffect(() => {
        if (selectedInquiryId) {
            fetchEmailStatus(selectedInquiryId);
        }
    }, [selectedInquiryId]);

    const handleInquiryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedInquiryId(e.target.value);
    };

    const handleRefresh = () => {
        if (selectedInquiryId) {
            fetchEmailStatus(selectedInquiryId);
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'delivered':
                return <Icon icon={FiCheck} size={16} />;
            case 'sending':
                return <Icon icon={FiClock} size={16} />;
            case 'failed':
                return <Icon icon={FiX} size={16} />;
            default:
                return null;
        }
    };

    if (loading && !inquiries.length) {
        return <LoadingState>Loading email status data...</LoadingState>;
    }

    if (error) {
        return <ErrorState>{error}</ErrorState>;
    }

    if (inquiries.length === 0) {
        return (
            <NoDataState>
                <IconWrapper $marginBottom={theme.space[4]} $opacity={0.5}>
                    <Icon icon={FiMail} size={48} />
                </IconWrapper>
                <h3>No inquiries found</h3>
                <p>There are no inquiries available to check email status.</p>
            </NoDataState>
        );
    }

    return (
        <EmailStatusContainer>
            <PageTitle>Email Status Monitor</PageTitle>

            <InquirySelector>
                <InquirySelect value={selectedInquiryId} onChange={handleInquiryChange}>
                    {inquiries.map(inquiry => (
                        <option key={inquiry.id} value={inquiry.id}>
                            {inquiry.name} ({inquiry.email}) - {new Date(inquiry.created_at).toLocaleDateString()}
                        </option>
                    ))}
                </InquirySelect>

                <RefreshButton onClick={handleRefresh}>
                    <Icon icon={FiRefreshCw} size={18} />
                    Refresh
                </RefreshButton>
            </InquirySelector>

            {loading ? (
                <LoadingState>Loading email status...</LoadingState>
            ) : emailStatus ? (
                <EmailStatusCard>
                    <EmailStatusHeader>
                        <EmailStatusTitle>
                            <Icon icon={FiMail} size={24} />
                            Email Delivery Status
                        </EmailStatusTitle>

                        <EmailStatusBadge $status={emailStatus.status}>
                            {getStatusIcon(emailStatus.status)}
                            {emailStatus.status.charAt(0).toUpperCase() + emailStatus.status.slice(1)}
                        </EmailStatusBadge>
                    </EmailStatusHeader>

                    <EmailStatusGrid>
                        <EmailStatusItem>
                            <EmailStatusLabel>Inquiry ID</EmailStatusLabel>
                            <EmailStatusValue>{emailStatus.inquiry_id}</EmailStatusValue>
                        </EmailStatusItem>

                        <EmailStatusItem>
                            <EmailStatusLabel>Customer Email</EmailStatusLabel>
                            <EmailStatusValue>{emailStatus.customer_email}</EmailStatusValue>
                        </EmailStatusItem>

                        <EmailStatusItem>
                            <EmailStatusLabel>Consultant Email</EmailStatusLabel>
                            <EmailStatusValue>{emailStatus.consultant_email}</EmailStatusValue>
                        </EmailStatusItem>

                        {emailStatus.sent_at && (
                            <EmailStatusItem>
                                <EmailStatusLabel>Sent At</EmailStatusLabel>
                                <EmailStatusValue>{new Date(emailStatus.sent_at).toLocaleString()}</EmailStatusValue>
                            </EmailStatusItem>
                        )}

                        {emailStatus.delivered_at && (
                            <EmailStatusItem>
                                <EmailStatusLabel>Delivered At</EmailStatusLabel>
                                <EmailStatusValue>{new Date(emailStatus.delivered_at).toLocaleString()}</EmailStatusValue>
                            </EmailStatusItem>
                        )}

                        {emailStatus.error_message && (
                            <EmailStatusItem>
                                <EmailStatusLabel>Error Message</EmailStatusLabel>
                                <EmailStatusValue style={{ color: theme.colors.danger }}>
                                    {emailStatus.error_message}
                                </EmailStatusValue>
                            </EmailStatusItem>
                        )}
                    </EmailStatusGrid>

                    <EmailStatusTimeline>
                        <TimelineTitle>Email Delivery Timeline</TimelineTitle>

                        <TimelineItem>
                            <TimelineIcon $status="complete">
                                <Icon icon={FiCheck} size={16} />
                            </TimelineIcon>
                            <TimelineContent>
                                <TimelineText>Inquiry received</TimelineText>
                                <TimelineTime>
                                    {emailStatus.sent_at
                                        ? new Date(new Date(emailStatus.sent_at).getTime() - 60000).toLocaleString()
                                        : 'N/A'}
                                </TimelineTime>
                            </TimelineContent>
                        </TimelineItem>

                        <TimelineItem>
                            <TimelineIcon $status={emailStatus.sent_at ? "complete" : "pending"}>
                                {emailStatus.sent_at
                                    ? <Icon icon={FiCheck} size={16} />
                                    : <Icon icon={FiClock} size={16} />}
                            </TimelineIcon>
                            <TimelineContent>
                                <TimelineText>Email sent to customer and consultant</TimelineText>
                                <TimelineTime>
                                    {emailStatus.sent_at
                                        ? new Date(emailStatus.sent_at).toLocaleString()
                                        : 'Pending'}
                                </TimelineTime>
                            </TimelineContent>
                        </TimelineItem>

                        <TimelineItem>
                            <TimelineIcon $status={emailStatus.delivered_at ? "complete" : emailStatus.status === 'failed' ? "error" : "pending"}>
                                {emailStatus.delivered_at
                                    ? <Icon icon={FiCheck} size={16} />
                                    : emailStatus.status === 'failed'
                                        ? <Icon icon={FiX} size={16} />
                                        : <Icon icon={FiClock} size={16} />}
                            </TimelineIcon>
                            <TimelineContent>
                                <TimelineText>
                                    {emailStatus.delivered_at
                                        ? 'Email delivered successfully'
                                        : emailStatus.status === 'failed'
                                            ? 'Email delivery failed'
                                            : 'Awaiting delivery confirmation'}
                                </TimelineText>
                                <TimelineTime>
                                    {emailStatus.delivered_at
                                        ? new Date(emailStatus.delivered_at).toLocaleString()
                                        : emailStatus.status === 'failed'
                                            ? 'Error occurred'
                                            : 'Pending'}
                                </TimelineTime>
                            </TimelineContent>
                        </TimelineItem>
                    </EmailStatusTimeline>
                </EmailStatusCard>
            ) : (
                <NoDataState>
                    <IconWrapper $marginBottom={theme.space[4]} $opacity={0.5}>
                        <Icon icon={FiMail} size={48} />
                    </IconWrapper>
                    <h3>No email status data available</h3>
                    <p>Select an inquiry to view its email delivery status.</p>
                </NoDataState>
            )}
        </EmailStatusContainer>
    );
};

export default EmailStatusMonitor;