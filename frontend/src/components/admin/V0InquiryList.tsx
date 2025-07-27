import React, { useState, useEffect } from "react";
import { 
  Search, 
  Download, 
  MoreHorizontal, 
  ArrowUpDown,
  SlidersHorizontal,
  Eye,
  FileText,
  Globe,
  Archive,
  AlertTriangle,
  Plus
} from "lucide-react";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { Card } from "../ui/card";
import { Checkbox } from "../ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { Input } from "../ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table";
import apiService, { Inquiry } from '../../services/api';
import V0ReportModal from './V0ReportModal';
import { V0DataAdapter } from './V0DataAdapter';
import { AnalysisReport } from './V0InquiryAnalysisSection';

// Loading and error states
const LoadingState = () => (
  <div className="flex justify-center items-center py-12">
    <div className="flex flex-col items-center space-y-4">
      <div className="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent"></div>
      <div className="text-muted-foreground">Loading inquiries...</div>
    </div>
  </div>
);

const ErrorState = ({ error, onRetry }: { error: string; onRetry: () => void }) => (
  <div className="flex flex-col items-center py-12 text-center">
    <div className="p-3 bg-red-100 rounded-full mb-4">
      <AlertTriangle className="w-6 h-6 text-red-600" />
    </div>
    <h3 className="text-lg font-semibold text-gray-900 mb-2">Failed to load inquiries</h3>
    <p className="text-sm text-muted-foreground mb-6 max-w-md">{error}</p>
    <Button onClick={onRetry} className="bg-blue-600 hover:bg-blue-700">
      <Download className="w-4 h-4 mr-2" />
      Try Again
    </Button>
  </div>
);

const V0InquiryList: React.FC = () => {
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [priorityFilter, setPriorityFilter] = useState("all");
  const [selectedInquiries, setSelectedInquiries] = useState<string[]>([]);
  const [selectedReport, setSelectedReport] = useState<AnalysisReport | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    fetchInquiries();
  }, []);

  const fetchInquiries = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiService.listInquiries({ limit: 100 });
      setInquiries(response.data || []);
    } catch (err: any) {
      console.error('Failed to fetch inquiries:', err);
      setError(err.message || 'Failed to load inquiries');
      setInquiries([]);
    } finally {
      setLoading(false);
    }
  };

  // Filter inquiries based on search and filters
  const filteredInquiries = (inquiries || []).filter((inquiry) => {
    const matchesSearch =
      searchQuery === "" ||
      inquiry.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      inquiry.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (inquiry.company || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
      (inquiry.services || []).some(service => 
        service.toLowerCase().includes(searchQuery.toLowerCase())
      ) ||
      inquiry.id.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesStatus = statusFilter === "all" || inquiry.status === statusFilter;
    const matchesPriority = priorityFilter === "all" || inquiry.priority === priorityFilter;

    return matchesSearch && matchesStatus && matchesPriority;
  });

  // Handle checkbox selection
  const toggleSelection = (id: string) => {
    setSelectedInquiries((prev) => (prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]));
  };

  // Handle select all
  const toggleSelectAll = () => {
    if (selectedInquiries.length === filteredInquiries.length) {
      setSelectedInquiries([]);
    } else {
      setSelectedInquiries(filteredInquiries.map((i) => i.id));
    }
  };

  // Format date to readable string
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  };

  // Get status badge color with v0 styling
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return "bg-yellow-50 text-yellow-700 border-yellow-200 hover:bg-yellow-100";
      case "processing":
        return "bg-blue-50 text-blue-700 border-blue-200 hover:bg-blue-100";
      case "reviewed":
        return "bg-purple-50 text-purple-700 border-purple-200 hover:bg-purple-100";
      case "responded":
        return "bg-green-50 text-green-700 border-green-200 hover:bg-green-100";
      case "closed":
        return "bg-gray-50 text-gray-700 border-gray-200 hover:bg-gray-100";
      default:
        return "bg-gray-50 text-gray-700 border-gray-200 hover:bg-gray-100";
    }
  };

  // Get priority badge color with v0 styling
  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case "high":
      case "urgent":
        return "bg-red-50 text-red-700 border-red-200 hover:bg-red-100";
      case "medium":
        return "bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-100";
      case "low":
        return "bg-green-50 text-green-700 border-green-200 hover:bg-green-100";
      default:
        return "bg-gray-50 text-gray-700 border-gray-200 hover:bg-gray-100";
    }
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

  const handlePreviewReport = (inquiry: Inquiry) => {
    console.log('Preview report clicked for inquiry:', inquiry);
    
    // Generate a preview report using the data adapter
    try {
      const analysisReport = V0DataAdapter.safeAdaptInquiryToAnalysisReport(inquiry, 0);
      console.log('Generated analysis report:', analysisReport);
      
      if (analysisReport) {
        setSelectedReport(analysisReport);
        setIsModalOpen(true);
        console.log('Modal should be opening...');
      } else {
        console.error('Analysis report is null');
        alert('Unable to generate preview for this inquiry.');
      }
    } catch (error) {
      console.error('Error generating preview report:', error);
      alert('Error generating preview report. Please try again.');
    }
  };

  const handleExport = (format: string) => {
    const timestamp = new Date().toISOString().split("T")[0];
    const fileName = `inquiries_export_${timestamp}.${format}`;

    console.log(`Exporting inquiries as ${fileName}`);
    // In real implementation, would generate and download the file
    alert(`Exporting ${filteredInquiries.length} inquiries as ${format.toUpperCase()}`);
  };

  const handleBulkExport = (selectedInquiries: string[]) => {
    if (selectedInquiries.length === 0) {
      alert("Please select inquiries to export");
      return;
    }

    const timestamp = new Date().toISOString().split("T")[0];
    const fileName = `selected_inquiries_${selectedInquiries.length}_items_${timestamp}.zip`;

    console.log(`Bulk exporting ${selectedInquiries.length} inquiries as ${fileName}`);
    alert(`Bulk exporting ${selectedInquiries.length} selected inquiries`);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedReport(null);
  };

  if (loading) {
    return <LoadingState />;
  }

  if (error) {
    return <ErrorState error={error} onRetry={fetchInquiries} />;
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Inquiries</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage and analyze customer inquiries with AI-powered insights
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" className="bg-transparent">
            <Plus className="h-4 w-4 mr-2" />
            New Inquiry
          </Button>
        </div>
      </div>

      {/* Search and Filters */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex w-full items-center gap-2 sm:max-w-sm">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search inquiries..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 h-9"
            />
          </div>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-9 gap-1 bg-transparent">
                <SlidersHorizontal className="h-4 w-4" />
                <span className="hidden sm:inline-block">Filters</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[200px]">
              <DropdownMenuLabel>Filter by</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <div className="p-2">
                <p className="mb-2 text-xs font-medium">Status</p>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger className="h-8">
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Statuses</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="processing">Processing</SelectItem>
                    <SelectItem value="reviewed">Reviewed</SelectItem>
                    <SelectItem value="responded">Responded</SelectItem>
                    <SelectItem value="closed">Closed</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <DropdownMenuSeparator />
              <div className="p-2">
                <p className="mb-2 text-xs font-medium">Priority</p>
                <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                  <SelectTrigger className="h-8">
                    <SelectValue placeholder="Select priority" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Priorities</SelectItem>
                    <SelectItem value="high">High</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="low">Low</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </DropdownMenuContent>
          </DropdownMenu>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-9 gap-1 bg-transparent">
                <Download className="h-4 w-4" />
                <span className="hidden sm:inline-block">Export</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => handleExport("csv")}>
                <FileText className="mr-2 h-4 w-4" />
                Export as CSV
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleExport("pdf")}>
                <FileText className="mr-2 h-4 w-4" />
                Export as PDF
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleExport("html")}>
                <Globe className="mr-2 h-4 w-4" />
                Export as HTML
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => handleBulkExport(selectedInquiries)}>
                <Archive className="mr-2 h-4 w-4" />
                Bulk Export Selected ({selectedInquiries.length})
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Table */}
      <Card className="border border-gray-200 shadow-sm">
        <div className="rounded-lg border border-gray-200 overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow className="bg-gray-50/50">
                <TableHead className="w-[40px]">
                  <Checkbox
                    checked={selectedInquiries.length === filteredInquiries.length && filteredInquiries.length > 0}
                    onCheckedChange={toggleSelectAll}
                    aria-label="Select all"
                  />
                </TableHead>
                <TableHead className="w-[100px]">
                  <div className="flex items-center font-medium text-gray-700">
                    ID
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead>
                  <div className="flex items-center font-medium text-gray-700">
                    Name
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead className="hidden md:table-cell font-medium text-gray-700">Email</TableHead>
                <TableHead className="hidden lg:table-cell font-medium text-gray-700">Company</TableHead>
                <TableHead className="hidden lg:table-cell font-medium text-gray-700">Services</TableHead>
                <TableHead className="font-medium text-gray-700">Status</TableHead>
                <TableHead className="hidden sm:table-cell font-medium text-gray-700">Priority</TableHead>
                <TableHead className="hidden sm:table-cell">
                  <div className="flex items-center font-medium text-gray-700">
                    Date
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead className="w-[120px] font-medium text-gray-700">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredInquiries.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={10} className="h-32 text-center">
                    <div className="flex flex-col items-center space-y-2">
                      <FileText className="h-8 w-8 text-gray-400" />
                      <p className="text-gray-500">No inquiries found.</p>
                      <p className="text-sm text-gray-400">Try adjusting your search or filters.</p>
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                filteredInquiries.map((inquiry) => (
                  <TableRow key={inquiry.id} className="hover:bg-gray-50/50 transition-colors">
                    <TableCell>
                      <Checkbox
                        checked={selectedInquiries.includes(inquiry.id)}
                        onCheckedChange={() => toggleSelection(inquiry.id)}
                        aria-label={`Select ${inquiry.id}`}
                      />
                    </TableCell>
                    <TableCell className="font-mono text-sm text-gray-600">{inquiry.id}</TableCell>
                    <TableCell className="font-medium text-gray-900">{inquiry.name}</TableCell>
                    <TableCell className="hidden md:table-cell text-gray-600">{inquiry.email}</TableCell>
                    <TableCell className="hidden lg:table-cell text-gray-600">{inquiry.company || '—'}</TableCell>
                    <TableCell className="hidden lg:table-cell text-gray-600">
                      <div className="max-w-[200px] truncate">
                        {(inquiry.services || []).join(', ')}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className={`border ${getStatusColor(inquiry.status)}`}>
                        {inquiry.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className={`border ${getPriorityColor(inquiry.priority)}`}>
                        {inquiry.priority}
                      </Badge>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell text-gray-600 text-sm">
                      {formatDate(inquiry.created_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handlePreviewReport(inquiry)}
                          title="Preview Report"
                          className="h-8 w-8 p-0 hover:bg-blue-50 hover:text-blue-600"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm" className="h-8 w-8 p-0 hover:bg-gray-100">
                              <MoreHorizontal className="h-4 w-4" />
                              <span className="sr-only">Open menu</span>
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => handlePreviewReport(inquiry)}>
                              <Eye className="mr-2 h-4 w-4" />
                              Preview Report
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={() => handleDownloadReport(inquiry.id, 'pdf')}>
                              <FileText className="mr-2 h-4 w-4" />
                              Download PDF Report
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleDownloadReport(inquiry.id, 'html')}>
                              <Globe className="mr-2 h-4 w-4" />
                              Download HTML Report
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
        
        {/* Pagination */}
        {filteredInquiries.length > 0 && (
          <div className="flex items-center justify-between px-6 py-4 border-t border-gray-200 bg-gray-50/30">
            <div className="text-sm text-muted-foreground">
              Showing <strong>{filteredInquiries.length}</strong> of <strong>{inquiries.length}</strong> inquiries
              {selectedInquiries.length > 0 && (
                <span className="ml-2 text-blue-600">
                  • <strong>{selectedInquiries.length}</strong> selected
                </span>
              )}
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" disabled className="bg-transparent">
                Previous
              </Button>
              <Button variant="outline" size="sm" className="bg-transparent">
                Next
              </Button>
            </div>
          </div>
        )}
      </Card>

      {/* Report Preview Modal */}
      <V0ReportModal
        report={selectedReport}
        isOpen={isModalOpen}
        onClose={closeModal}
        onDownload={handleDownloadReport}
      />
    </div>
  );
};

export default V0InquiryList;