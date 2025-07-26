import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { FiSearch, FiDownload, FiEye } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import apiService, { Inquiry } from '../../services/api';
import Icon from '../ui/Icon';

const InquiriesContainer = styled.div`
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

const FiltersContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: ${theme.space[4]};
  margin-bottom: ${theme.space[4]};
  
  @media (max-width: ${theme.breakpoints.md}) {
    flex-direction: column;
  }
`;

const SearchContainer = styled.div`
  position: relative;
  flex: 1;
  min-width: 250px;
`;

const SearchInput = styled.input`
  width: 100%;
  padding: ${theme.space[3]} ${theme.space[3]} ${theme.space[3]} ${theme.space[10]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
  }
`;

const SearchIcon = styled.div`
  position: absolute;
  left: ${theme.space[3]};
  top: 50%;
  transform: translateY(-50%);
  color: ${theme.colors.gray500};
`;

const FilterSelect = styled.select`
  padding: ${theme.space[3]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  background-color: ${theme.colors.white};
  min-width: 150px;
  
  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
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

const ActionButton = styled.button`
  background: none;
  border: none;
  color: ${theme.colors.primary};
  cursor: pointer;
  padding: ${theme.space[2]};
  border-radius: ${theme.borderRadius.md};
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.gray100};
  }
  
  &:not(:last-child) {
    margin-right: ${theme.space[2]};
  }
`;

const PaginationContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: ${theme.space[4]};
`;

const PaginationInfo = styled.div`
  color: ${theme.colors.gray600};
`;

const PaginationButtons = styled.div`
  display: flex;
  gap: ${theme.space[2]};
`;

const PaginationButton = styled.button<{ $active?: boolean }>`
  padding: ${theme.space[2]} ${theme.space[4]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  background-color: ${props => props.$active ? theme.colors.primary : theme.colors.white};
  color: ${props => props.$active ? theme.colors.white : theme.colors.gray800};
  cursor: ${props => props.$active ? 'default' : 'pointer'};
  transition: ${theme.transitions.normal};
  
  &:hover:not(:disabled) {
    background-color: ${props => props.$active ? theme.colors.primary : theme.colors.gray100};
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
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

const InquiriesList: React.FC = () => {
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    status: '',
    service: '',
    search: '',
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    pages: 0,
  });
  
  const fetchInquiries = async () => {
    try {
      setLoading(true);
      
      const offset = (pagination.page - 1) * pagination.limit;
      
      // Build filter object
      const apiFilters: any = {
        limit: pagination.limit,
        offset: offset,
      };
      
      if (filters.status) {
        apiFilters.status = filters.status;
      }
      
      if (filters.service) {
        apiFilters.service = filters.service;
      }
      
      // Note: search is not implemented in the backend yet
      // This would typically be handled by the backend
      
      const response = await apiService.listInquiries(apiFilters);
      
      setInquiries(response.data);
      setPagination(prev => ({
        ...prev,
        total: response.total,
        pages: response.pages,
      }));
      
      setLoading(false);
    } catch (err) {
      console.error('Failed to fetch inquiries:', err);
      setError('Failed to load inquiries. Please try again later.');
      setLoading(false);
    }
  };
  
  useEffect(() => {
    fetchInquiries();
  }, [pagination.page, filters.status, filters.service]);
  
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters(prev => ({ ...prev, search: e.target.value }));
  };
  
  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setFilters(prev => ({ ...prev, status: e.target.value }));
    setPagination(prev => ({ ...prev, page: 1 })); // Reset to first page
  };
  
  const handleServiceChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setFilters(prev => ({ ...prev, service: e.target.value }));
    setPagination(prev => ({ ...prev, page: 1 })); // Reset to first page
  };
  
  const handlePageChange = (page: number) => {
    setPagination(prev => ({ ...prev, page }));
  };
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchInquiries();
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
  
  if (loading && inquiries.length === 0) {
    return <LoadingState>Loading inquiries...</LoadingState>;
  }
  
  if (error) {
    return <ErrorState>{error}</ErrorState>;
  }
  
  return (
    <InquiriesContainer>
      <PageTitle>Inquiries</PageTitle>
      
      <FiltersContainer>
        <SearchContainer>
          <form onSubmit={handleSearch}>
            <SearchInput
              type="text"
              placeholder="Search by name or email..."
              value={filters.search}
              onChange={handleSearchChange}
            />
            <SearchIcon>
              <Icon icon={FiSearch} size={18} />
            </SearchIcon>
          </form>
        </SearchContainer>
        
        <FilterSelect value={filters.status} onChange={handleStatusChange}>
          <option value="">All Statuses</option>
          <option value="pending">Pending</option>
          <option value="processing">Processing</option>
          <option value="reviewed">Reviewed</option>
          <option value="responded">Responded</option>
          <option value="closed">Closed</option>
        </FilterSelect>
        
        <FilterSelect value={filters.service} onChange={handleServiceChange}>
          <option value="">All Services</option>
          <option value="assessment">Assessment</option>
          <option value="migration">Migration</option>
          <option value="optimization">Optimization</option>
          <option value="architecture_review">Architecture Review</option>
        </FilterSelect>
      </FiltersContainer>
      
      {inquiries.length > 0 ? (
        <>
          <Table>
            <TableHead>
              <TableRow>
                <TableHeader>Name</TableHeader>
                <TableHeader>Email</TableHeader>
                <TableHeader>Service</TableHeader>
                <TableHeader>Status</TableHeader>
                <TableHeader>Date</TableHeader>
                <TableHeader>Actions</TableHeader>
              </TableRow>
            </TableHead>
            <tbody>
              {inquiries.map((inquiry) => (
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
                  <TableCell>
                    <ActionButton 
                      title="View Details"
                      onClick={() => alert(`View details for inquiry ${inquiry.id}`)}
                    >
                      <Icon icon={FiEye} size={18} />
                    </ActionButton>
                    <ActionButton
                      title="Download PDF"
                      onClick={() => handleDownloadReport(inquiry.id, 'pdf')}
                    >
                      <Icon icon={FiDownload} size={18} />
                    </ActionButton>
                  </TableCell>
                </TableRow>
              ))}
            </tbody>
          </Table>
          
          <PaginationContainer>
            <PaginationInfo>
              Showing {Math.min((pagination.page - 1) * pagination.limit + 1, pagination.total)} to {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} inquiries
            </PaginationInfo>
            
            <PaginationButtons>
              <PaginationButton
                onClick={() => handlePageChange(pagination.page - 1)}
                disabled={pagination.page === 1}
              >
                Previous
              </PaginationButton>
              
              {Array.from({ length: Math.min(5, pagination.pages) }, (_, i) => {
                // Show pages around current page
                let pageNum;
                if (pagination.pages <= 5) {
                  pageNum = i + 1;
                } else if (pagination.page <= 3) {
                  pageNum = i + 1;
                } else if (pagination.page >= pagination.pages - 2) {
                  pageNum = pagination.pages - 4 + i;
                } else {
                  pageNum = pagination.page - 2 + i;
                }
                
                return (
                  <PaginationButton
                    key={pageNum}
                    $active={pageNum === pagination.page}
                    onClick={() => handlePageChange(pageNum)}
                  >
                    {pageNum}
                  </PaginationButton>
                );
              })}
              
              <PaginationButton
                onClick={() => handlePageChange(pagination.page + 1)}
                disabled={pagination.page === pagination.pages}
              >
                Next
              </PaginationButton>
            </PaginationButtons>
          </PaginationContainer>
        </>
      ) : (
        <div>No inquiries found matching your filters.</div>
      )}
    </InquiriesContainer>
  );
};

export default InquiriesList;