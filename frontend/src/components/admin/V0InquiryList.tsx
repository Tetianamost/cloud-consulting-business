import React, { useState, useEffect } from "react";
import { 
  FiSearch, 
  FiDownload, 
  FiMoreHorizontal, 
  FiArrowUp,
  FiFilter,
  FiEye,
  FiCheckCircle,
  FiClock,
  FiAlertCircle
} from "react-icons/fi";
import styled from 'styled-components';
import { theme } from '../../styles/theme';
import apiService, { Inquiry } from '../../services/api';
import Icon from '../ui/Icon';

// Styled Components
const Container = styled.div`
  padding: ${theme.space[6]};
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: ${theme.space[6]};
`;

const Title = styled.h1`
  font-size: ${theme.fontSizes['2xl']};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
`;

const Controls = styled.div`
  display: flex;
  gap: ${theme.space[3]};
  align-items: center;
`;

const SearchContainer = styled.div`
  position: relative;
  width: 300px;
`;

const SearchInput = styled.input`
  width: 100%;
  padding: ${theme.space[2]} ${theme.space[3]} ${theme.space[2]} ${theme.space[10]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.sm};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
  }
`;

const SearchIcon = styled.div`
  position: absolute;
  left: ${theme.space[3]};
  top: 50%;
  transform: translateY(-50%);
  color: ${theme.colors.gray400};
`;

const Button = styled.button<{ $variant?: 'primary' | 'outline' }>`
  display: flex;
  align-items: center;
  gap: ${theme.space[2]};
  padding: ${theme.space[2]} ${theme.space[4]};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  ${props => props.$variant === 'primary' ? `
    background-color: ${theme.colors.primary};
    color: ${theme.colors.white};
    border: 1px solid ${theme.colors.primary};
    
    &:hover {
      background-color: ${theme.colors.primary};
      opacity: 0.9;
    }
  ` : `
    background-color: ${theme.colors.white};
    color: ${theme.colors.gray700};
    border: 1px solid ${theme.colors.gray300};
    
    &:hover {
      background-color: ${theme.colors.gray100};
    }
  `}
`;

const Select = styled.select`
  padding: ${theme.space[2]} ${theme.space[3]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.sm};
  background-color: ${theme.colors.white};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
  }
`;

const Card = styled.div`
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  border: 1px solid ${theme.colors.gray200};
  overflow: hidden;
  box-shadow: ${theme.shadows.sm};
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

const TableHeader = styled.thead`
  background-color: ${theme.colors.gray100};
`;

const TableRow = styled.tr`
  border-bottom: 1px solid ${theme.colors.gray200};
  
  &:hover {
    background-color: ${theme.colors.gray100};
  }
`;

const TableHead = styled.th`
  padding: ${theme.space[3]} ${theme.space[4]};
  text-align: left;
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray700};
  font-size: ${theme.fontSizes.sm};
`;

const TableCell = styled.td`
  padding: ${theme.space[3]} ${theme.space[4]};
  font-size: ${theme.fontSizes.sm};
  color: ${theme.colors.gray900};
`;

const Badge = styled.span<{ $variant: 'success' | 'warning' | 'danger' | 'info' }>`
  display: inline-flex;
  align-items: center;
  padding: ${theme.space[1]} ${theme.space[2]};
  border-radius: ${theme.borderRadius.full};
  font-size: ${theme.fontSizes.xs};
  font-weight: ${theme.fontWeights.medium};
  
  ${props => {
    switch (props.$variant) {
      case 'success':
        return `
          background-color: ${theme.colors.success}20;
          color: ${theme.colors.success};
        `;
      case 'warning':
        return `
          background-color: ${theme.colors.warning}20;
          color: ${theme.colors.warning};
        `;
      case 'danger':
        return `
          background-color: ${theme.colors.danger}20;
          color: ${theme.colors.danger};
        `;
      case 'info':
      default:
        return `
          background-color: ${theme.colors.info}20;
          color: ${theme.colors.info};
        `;
    }
  }}
`;

const ActionButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  border-radius: ${theme.borderRadius.md};
  color: ${theme.colors.gray500};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.gray100};
    color: ${theme.colors.gray700};
  }
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

const V0InquiryList: React.FC = () => {
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");

  useEffect(() => {
    fetchInquiries();
  }, []);

  const fetchInquiries = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiService.listInquiries({ limit: 100 });
      setInquiries(response.data);
    } catch (err: any) {
      console.error('Failed to fetch inquiries:', err);
      setError(err.message || 'Failed to load inquiries');
    } finally {
      setLoading(false);
    }
  };

  // Filter inquiries based on search and filters
  const filteredInquiries = inquiries.filter((inquiry) => {
    const matchesSearch =
      searchQuery === "" ||
      inquiry.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      inquiry.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      inquiry.company.toLowerCase().includes(searchQuery.toLowerCase()) ||
      inquiry.services.some(service => 
        service.toLowerCase().includes(searchQuery.toLowerCase())
      );

    const matchesStatus = statusFilter === "all" || inquiry.status === statusFilter;

    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return <Badge $variant="warning">Pending</Badge>;
      case 'processing':
        return <Badge $variant="info">Processing</Badge>;
      case 'reviewed':
        return <Badge $variant="success">Reviewed</Badge>;
      case 'responded':
        return <Badge $variant="success">Responded</Badge>;
      case 'closed':
        return <Badge $variant="success">Closed</Badge>;
      default:
        return <Badge $variant="info">{status}</Badge>;
    }
  };

  const getPriorityBadge = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'high':
      case 'urgent':
        return <Badge $variant="danger">{priority}</Badge>;
      case 'medium':
        return <Badge $variant="warning">{priority}</Badge>;
      case 'low':
        return <Badge $variant="success">{priority}</Badge>;
      default:
        return <Badge $variant="info">{priority}</Badge>;
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  };

  const handleDownloadReport = async (inquiryId: string, format: 'pdf' | 'html') => {
    try {
      const blob = await apiService.downloadReport(inquiryId, format);
      
      // Create a download link
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.style.display = 'none';
      a.href = url;
      a.download = `report-${inquiryId}.${format}`;
      document.body.appendChild(a);
      a.click();
      
      // Clean up
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      console.error(`Failed to download ${format} report:`, err);
      alert(`Failed to download ${format} report. Please try again later.`);
    }
  };

  if (loading) {
    return <LoadingState>Loading inquiries...</LoadingState>;
  }

  if (error) {
    return (
      <ErrorState>
        <p>Failed to load inquiries</p>
        <p>{error}</p>
        <Button $variant="primary" onClick={fetchInquiries} style={{ marginTop: theme.space[4] }}>
          Try Again
        </Button>
      </ErrorState>
    );
  }

  return (
    <Container>
      <Header>
        <Title>Inquiries</Title>
        <Controls>
          <SearchContainer>
            <SearchIcon>
              <Icon icon={FiSearch} size={16} />
            </SearchIcon>
            <SearchInput
              placeholder="Search inquiries..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </SearchContainer>
          
          <Select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
            <option value="all">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="processing">Processing</option>
            <option value="reviewed">Reviewed</option>
            <option value="responded">Responded</option>
            <option value="closed">Closed</option>
          </Select>
          
          <Button $variant="outline">
            <Icon icon={FiDownload} size={16} />
            Export
          </Button>
        </Controls>
      </Header>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>
                <div style={{ display: 'flex', alignItems: 'center', gap: theme.space[1] }}>
                  Name
                  <Icon icon={FiArrowUp} size={12} />
                </div>
              </TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Company</TableHead>
              <TableHead>Services</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Priority</TableHead>
              <TableHead>
                <div style={{ display: 'flex', alignItems: 'center', gap: theme.space[1] }}>
                  Date
                  <Icon icon={FiArrowUp} size={12} />
                </div>
              </TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <tbody>
            {filteredInquiries.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} style={{ textAlign: 'center', padding: theme.space[8] }}>
                  No inquiries found.
                </TableCell>
              </TableRow>
            ) : (
              filteredInquiries.map((inquiry) => (
                <TableRow key={inquiry.id}>
                  <TableCell style={{ fontWeight: theme.fontWeights.medium }}>
                    {inquiry.name}
                  </TableCell>
                  <TableCell>{inquiry.email}</TableCell>
                  <TableCell>{inquiry.company || '—'}</TableCell>
                  <TableCell>{inquiry.services.join(', ')}</TableCell>
                  <TableCell>{getStatusBadge(inquiry.status)}</TableCell>
                  <TableCell>{getPriorityBadge(inquiry.priority)}</TableCell>
                  <TableCell>{formatDate(inquiry.created_at)}</TableCell>
                  <TableCell>
                    <div style={{ display: 'flex', gap: theme.space[1] }}>
                      <ActionButton
                        title="View Details"
                        onClick={() => alert(`View details for inquiry ${inquiry.id}`)}
                      >
                        <Icon icon={FiEye} size={16} />
                      </ActionButton>
                      <ActionButton
                        title="Download PDF"
                        onClick={() => handleDownloadReport(inquiry.id, 'pdf')}
                      >
                        <Icon icon={FiDownload} size={16} />
                      </ActionButton>
                      <ActionButton
                        title="More Actions"
                        onClick={() => alert(`More actions for inquiry ${inquiry.id}`)}
                      >
                        <Icon icon={FiMoreHorizontal} size={16} />
                      </ActionButton>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </tbody>
        </Table>
        
        {filteredInquiries.length > 0 && (
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: theme.space[4],
            borderTop: `1px solid ${theme.colors.gray200}`
          }}>
            <div style={{ fontSize: theme.fontSizes.sm, color: theme.colors.gray600 }}>
              Showing <strong>{filteredInquiries.length}</strong> of <strong>{inquiries.length}</strong> inquiries
            </div>
            <div style={{ display: 'flex', gap: theme.space[2] }}>
              <Button $variant="outline" disabled>
                Previous
              </Button>
              <Button $variant="outline">
                Next
              </Button>
            </div>
          </div>
        )}
      </Card>
    </Container>
  );
};

export default V0InquiryList;