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
  Plus,
  X,
  Filter,
  ChevronDown,
  Calendar,
  User,
  Building,
  Tag,
  Trash2,
  Mail,
  CheckCircle2,
  Clock,
  ArrowUp,
  ArrowDown
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
  const [serviceFilter, setServiceFilter] = useState("all");
  const [dateFilter, setDateFilter] = useState("all");
  const [selectedInquiries, setSelectedInquiries] = useState<string[]>([]);
  const [selectedReport, setSelectedReport] = useState<AnalysisReport | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [sortField, setSortField] = useState<string>("created_at");
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("desc");
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(false);
  const [searchFocused, setSearchFocused] = useState(false);
  const [bulkActionMode, setBulkActionMode] = useState(false);
  const [isExporting, setIsExporting] = useState(false);
  const [lastSearchTime, setLastSearchTime] = useState<number>(0);

  useEffect(() => {
    fetchInquiries();
  }, []);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Ctrl/Cmd + K to focus search
      if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
        event.preventDefault();
        const searchInput = document.querySelector('input[placeholder*="Search"]') as HTMLInputElement;
        if (searchInput) {
          searchInput.focus();
        }
      }
      
      // Escape to clear search and filters
      if (event.key === 'Escape') {
        if (searchQuery || getActiveFilterCount() > 0) {
          clearAllFilters();
        }
      }
      
      // Ctrl/Cmd + A to select all visible inquiries
      if ((event.ctrlKey || event.metaKey) && event.key === 'a' && !searchFocused) {
        event.preventDefault();
        const currentFilteredInquiries = inquiries.filter((inquiry) => {
          const matchesSearch =
            searchQuery === "" ||
            inquiry.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            inquiry.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (inquiry.company || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
            (inquiry.services || []).some(service => 
              service.toLowerCase().includes(searchQuery.toLowerCase())
            ) ||
            inquiry.id.toLowerCase().includes(searchQuery.toLowerCase());
          return matchesSearch;
        });
        
        if (selectedInquiries.length === currentFilteredInquiries.length) {
          setSelectedInquiries([]);
          setBulkActionMode(false);
        } else {
          setSelectedInquiries(currentFilteredInquiries.map((i) => i.id));
          setBulkActionMode(true);
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [searchQuery, selectedInquiries, searchFocused]);

  // Debounced search performance tracking
  useEffect(() => {
    if (searchQuery) {
      setLastSearchTime(Date.now());
    }
  }, [searchQuery]);

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

  // Get unique values for filter options
  const getUniqueServices = () => {
    const services = new Set<string>();
    inquiries.forEach(inquiry => {
      (inquiry.services || []).forEach(service => services.add(service));
    });
    return Array.from(services).sort();
  };

  // Enhanced filtering logic
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
    
    const matchesService = serviceFilter === "all" || 
      (inquiry.services || []).some(service => service === serviceFilter);
    
    const matchesDate = (() => {
      if (dateFilter === "all") return true;
      const inquiryDate = new Date(inquiry.created_at);
      const now = new Date();
      
      switch (dateFilter) {
        case "today":
          return inquiryDate.toDateString() === now.toDateString();
        case "week":
          const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
          return inquiryDate >= weekAgo;
        case "month":
          const monthAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
          return inquiryDate >= monthAgo;
        case "quarter":
          const quarterAgo = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000);
          return inquiryDate >= quarterAgo;
        default:
          return true;
      }
    })();

    return matchesSearch && matchesStatus && matchesPriority && matchesService && matchesDate;
  });

  // Enhanced sorting logic
  const sortedInquiries = [...filteredInquiries].sort((a, b) => {
    let aValue: any, bValue: any;
    
    switch (sortField) {
      case "name":
        aValue = a.name.toLowerCase();
        bValue = b.name.toLowerCase();
        break;
      case "email":
        aValue = a.email.toLowerCase();
        bValue = b.email.toLowerCase();
        break;
      case "company":
        aValue = (a.company || '').toLowerCase();
        bValue = (b.company || '').toLowerCase();
        break;
      case "status":
        aValue = a.status.toLowerCase();
        bValue = b.status.toLowerCase();
        break;
      case "priority":
        const priorityOrder = { low: 1, medium: 2, high: 3, urgent: 4 };
        aValue = priorityOrder[a.priority.toLowerCase() as keyof typeof priorityOrder] || 0;
        bValue = priorityOrder[b.priority.toLowerCase() as keyof typeof priorityOrder] || 0;
        break;
      case "created_at":
      default:
        aValue = new Date(a.created_at).getTime();
        bValue = new Date(b.created_at).getTime();
        break;
    }
    
    if (aValue < bValue) return sortDirection === "asc" ? -1 : 1;
    if (aValue > bValue) return sortDirection === "asc" ? 1 : -1;
    return 0;
  });

  // Handle checkbox selection
  const toggleSelection = (id: string) => {
    setSelectedInquiries((prev) => (prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]));
  };

  // Handle select all
  const toggleSelectAll = () => {
    if (selectedInquiries.length === sortedInquiries.length) {
      setSelectedInquiries([]);
      setBulkActionMode(false);
    } else {
      setSelectedInquiries(sortedInquiries.map((i) => i.id));
      setBulkActionMode(true);
    }
  };

  // Handle sorting
  const handleSort = (field: string) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  // Clear all filters
  const clearAllFilters = () => {
    setSearchQuery("");
    setStatusFilter("all");
    setPriorityFilter("all");
    setServiceFilter("all");
    setDateFilter("all");
    setSortField("created_at");
    setSortDirection("desc");
  };

  // Get active filter count
  const getActiveFilterCount = () => {
    let count = 0;
    if (searchQuery) count++;
    if (statusFilter !== "all") count++;
    if (priorityFilter !== "all") count++;
    if (serviceFilter !== "all") count++;
    if (dateFilter !== "all") count++;
    return count;
  };

  // Bulk actions
  const handleBulkStatusUpdate = (newStatus: string) => {
    if (selectedInquiries.length === 0) return;
    
    console.log(`Updating ${selectedInquiries.length} inquiries to status: ${newStatus}`);
    // In real implementation, would call API to update status
    alert(`Updated ${selectedInquiries.length} inquiries to ${newStatus} status`);
    setSelectedInquiries([]);
    setBulkActionMode(false);
  };

  const handleBulkDelete = () => {
    if (selectedInquiries.length === 0) return;
    
    const confirmed = window.confirm(`Are you sure you want to delete ${selectedInquiries.length} inquiries? This action cannot be undone.`);
    if (confirmed) {
      console.log(`Deleting ${selectedInquiries.length} inquiries`);
      // In real implementation, would call API to delete inquiries
      alert(`Deleted ${selectedInquiries.length} inquiries`);
      setSelectedInquiries([]);
      setBulkActionMode(false);
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

  const handleExport = async (format: string) => {
    setIsExporting(true);
    const timestamp = new Date().toISOString().split("T")[0];
    const fileName = `inquiries_export_${timestamp}.${format}`;

    try {
      console.log(`Exporting inquiries as ${fileName}`);
      // Simulate export delay for better UX
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // In real implementation, would generate and download the file
      alert(`Successfully exported ${filteredInquiries.length} inquiries as ${format.toUpperCase()}`);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Export failed. Please try again.');
    } finally {
      setIsExporting(false);
    }
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
    <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-3 sm:gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0">
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Inquiries</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage and analyze customer inquiries with AI-powered insights
          </p>
        </div>
        <div className="flex items-center gap-2 flex-shrink-0">
          <Button variant="outline" size="sm" className="bg-transparent">
            <Plus className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">New Inquiry</span>
            <span className="sm:hidden">New</span>
          </Button>
        </div>
      </div>

      {/* Enhanced Search and Filters */}
      <div className="flex flex-col gap-3 sm:gap-4">
        {/* Search bar - full width on mobile */}
        <div className="flex w-full items-center gap-2">
          <div className="relative flex-1 group">
            <Search className={`absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 transition-all duration-200 ${
              searchFocused ? 'text-blue-500 scale-110' : 'text-muted-foreground group-hover:text-gray-600'
            }`} />
            <Input
              placeholder="Search inquiries..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setSearchFocused(false)}
              className={`pl-10 h-9 transition-all duration-200 ${
                searchFocused 
                  ? 'ring-2 ring-blue-500 ring-opacity-20 border-blue-300 shadow-sm' 
                  : 'hover:border-gray-300 hover:shadow-sm'
              } ${searchQuery ? 'bg-blue-50/30 border-blue-200' : 'bg-white'}`}
            />
            {searchQuery && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSearchQuery("")}
                className="absolute right-1 top-1/2 transform -translate-y-1/2 h-7 w-7 p-0 hover:bg-gray-100 transition-all duration-150"
                aria-label="Clear search"
              >
                <X className="h-3 w-3" />
              </Button>
            )}
          </div>
          
          {/* Clear filters button - always visible when filters active */}
          {getActiveFilterCount() > 0 && (
            <Button
              variant="outline"
              size="sm"
              onClick={clearAllFilters}
              className="h-9 px-2 sm:px-3 text-xs bg-blue-50 border-blue-200 text-blue-700 hover:bg-blue-100 transition-all duration-150 shadow-sm hover:shadow flex-shrink-0"
            >
              <X className="h-3 w-3 sm:mr-1" />
              <span className="hidden sm:inline">Clear ({getActiveFilterCount()})</span>
            </Button>
          )}
        </div>

        {/* Search results indicator */}
        {searchQuery && (
          <div className="text-xs text-muted-foreground">
            {filteredInquiries.length} result{filteredInquiries.length !== 1 ? 's' : ''} found
          </div>
        )}

        {/* Filters row */}
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          {/* Quick Filters */}
          <div className="flex items-center gap-2 flex-wrap">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className={`h-9 w-[110px] sm:w-[130px] bg-white border-gray-200 hover:border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:ring-opacity-20 transition-all duration-150 shadow-sm hover:shadow ${
                statusFilter !== 'all' ? 'bg-blue-50 border-blue-200 text-blue-700' : ''
              }`}>
                <div className="flex items-center gap-1 sm:gap-2">
                  <Tag className={`h-3 w-3 transition-colors ${
                    statusFilter !== 'all' ? 'text-blue-600' : 'text-gray-500'
                  }`} />
                  <SelectValue placeholder="Status" />
                </div>
              </SelectTrigger>
              <SelectContent className="border-gray-200 shadow-lg">
                <SelectItem value="all" className="hover:bg-gray-50">All Statuses</SelectItem>
                <SelectItem value="pending" className="hover:bg-yellow-50">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-yellow-400"></div>
                    Pending
                  </div>
                </SelectItem>
                <SelectItem value="processing" className="hover:bg-blue-50">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-blue-400"></div>
                    Processing
                  </div>
                </SelectItem>
                <SelectItem value="reviewed" className="hover:bg-purple-50">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-purple-400"></div>
                    Reviewed
                  </div>
                </SelectItem>
                <SelectItem value="responded" className="hover:bg-green-50">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-green-400"></div>
                    Responded
                  </div>
                </SelectItem>
                <SelectItem value="closed" className="hover:bg-gray-50">
                  <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-gray-400"></div>
                    Closed
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>

            <Select value={priorityFilter} onValueChange={setPriorityFilter}>
              <SelectTrigger className={`h-9 w-[110px] sm:w-[130px] bg-white border-gray-200 hover:border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:ring-opacity-20 transition-all duration-150 shadow-sm hover:shadow ${
                priorityFilter !== 'all' ? 'bg-amber-50 border-amber-200 text-amber-700' : ''
              }`}>
                <div className="flex items-center gap-1 sm:gap-2">
                  <AlertTriangle className={`h-3 w-3 transition-colors ${
                    priorityFilter !== 'all' ? 'text-amber-600' : 'text-gray-500'
                  }`} />
                  <SelectValue placeholder="Priority" />
                </div>
              </SelectTrigger>
              <SelectContent className="border-gray-200 shadow-lg">
                <SelectItem value="all" className="hover:bg-gray-50">All Priorities</SelectItem>
                <SelectItem value="high" className="hover:bg-red-50">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-3 w-3 text-red-500" />
                    High
                  </div>
                </SelectItem>
                <SelectItem value="medium" className="hover:bg-amber-50">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-3 w-3 text-amber-500" />
                    Medium
                  </div>
                </SelectItem>
                <SelectItem value="low" className="hover:bg-green-50">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-3 w-3 text-green-500" />
                    Low
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Advanced Filters Toggle */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button 
                variant="outline" 
                size="sm" 
                className={`h-9 gap-1 transition-all duration-200 shadow-sm hover:shadow ${
                  getActiveFilterCount() > 2 
                    ? 'bg-blue-50 border-blue-200 text-blue-700 shadow-blue-100' 
                    : 'bg-white hover:bg-blue-50 hover:border-blue-300'
                }`}
              >
                <Filter className={`h-4 w-4 transition-transform ${
                  getActiveFilterCount() > 2 ? 'text-blue-600' : ''
                }`} />
                <span className="hidden sm:inline-block">More Filters</span>
                {getActiveFilterCount() > 2 && (
                  <span className="ml-1 px-1.5 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full animate-pulse">
                    {getActiveFilterCount() - 2}
                  </span>
                )}
                <ChevronDown className="h-3 w-3 transition-transform group-data-[state=open]:rotate-180" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-[300px] p-4 border-gray-200 shadow-xl">
              <div className="space-y-4">
                <div className="pb-2 border-b border-gray-100">
                  <h4 className="text-sm font-semibold text-gray-900 flex items-center gap-2">
                    <SlidersHorizontal className="h-4 w-4" />
                    Advanced Filters
                  </h4>
                </div>

                <div>
                  <label className="text-sm font-medium text-gray-700 mb-2 block flex items-center gap-2">
                    <Building className="h-3 w-3" />
                    Service Type
                  </label>
                  <Select value={serviceFilter} onValueChange={setServiceFilter}>
                    <SelectTrigger className={`h-9 transition-all duration-150 ${
                      serviceFilter !== 'all' ? 'bg-green-50 border-green-200 text-green-700' : ''
                    }`}>
                      <SelectValue placeholder="All Services" />
                    </SelectTrigger>
                    <SelectContent className="border-gray-200 shadow-lg">
                      <SelectItem value="all" className="hover:bg-gray-50">All Services</SelectItem>
                      {getUniqueServices().map(service => (
                        <SelectItem key={service} value={service} className="hover:bg-green-50">
                          <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full bg-green-400"></div>
                            {service}
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <label className="text-sm font-medium text-gray-700 mb-2 block flex items-center gap-2">
                    <Calendar className="h-3 w-3" />
                    Date Range
                  </label>
                  <Select value={dateFilter} onValueChange={setDateFilter}>
                    <SelectTrigger className={`h-9 transition-all duration-150 ${
                      dateFilter !== 'all' ? 'bg-purple-50 border-purple-200 text-purple-700' : ''
                    }`}>
                      <div className="flex items-center gap-2">
                        <Calendar className={`h-3 w-3 ${
                          dateFilter !== 'all' ? 'text-purple-600' : 'text-gray-500'
                        }`} />
                        <SelectValue placeholder="All Time" />
                      </div>
                    </SelectTrigger>
                    <SelectContent className="border-gray-200 shadow-lg">
                      <SelectItem value="all" className="hover:bg-gray-50">All Time</SelectItem>
                      <SelectItem value="today" className="hover:bg-purple-50">
                        <div className="flex items-center gap-2">
                          <Clock className="h-3 w-3 text-purple-500" />
                          Today
                        </div>
                      </SelectItem>
                      <SelectItem value="week" className="hover:bg-purple-50">
                        <div className="flex items-center gap-2">
                          <Calendar className="h-3 w-3 text-purple-500" />
                          Last 7 days
                        </div>
                      </SelectItem>
                      <SelectItem value="month" className="hover:bg-purple-50">
                        <div className="flex items-center gap-2">
                          <Calendar className="h-3 w-3 text-purple-500" />
                          Last 30 days
                        </div>
                      </SelectItem>
                      <SelectItem value="quarter" className="hover:bg-purple-50">
                        <div className="flex items-center gap-2">
                          <Calendar className="h-3 w-3 text-purple-500" />
                          Last 90 days
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="pt-3 border-t border-gray-100">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs text-gray-500">Active filters: {getActiveFilterCount()}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={clearAllFilters}
                      className="h-6 px-2 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                    >
                      <X className="h-3 w-3 mr-1" />
                      Reset All
                    </Button>
                  </div>
                </div>
              </div>
            </DropdownMenuContent>
          </DropdownMenu>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button 
                variant="outline" 
                size="sm" 
                className="h-9 gap-1 bg-white hover:bg-gray-50 shadow-sm hover:shadow transition-all duration-150"
              >
                <Download className="h-4 w-4" />
                <span className="hidden sm:inline-block">Export</span>
                <ChevronDown className="h-3 w-3 ml-1 transition-transform group-data-[state=open]:rotate-180" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="border-gray-200 shadow-lg">
              <div className="px-3 py-2 border-b border-gray-100">
                <p className="text-xs font-medium text-gray-700">Export Options</p>
                <p className="text-xs text-gray-500">{filteredInquiries.length} inquiries</p>
              </div>
              <DropdownMenuItem 
                onClick={() => handleExport("csv")}
                className="hover:bg-green-50 transition-colors"
              >
                <FileText className="mr-2 h-4 w-4 text-green-600" />
                <div className="flex flex-col">
                  <span>Export as CSV</span>
                  <span className="text-xs text-gray-500">Spreadsheet format</span>
                </div>
              </DropdownMenuItem>
              <DropdownMenuItem 
                onClick={() => handleExport("pdf")}
                className="hover:bg-red-50 transition-colors"
              >
                <FileText className="mr-2 h-4 w-4 text-red-600" />
                <div className="flex flex-col">
                  <span>Export as PDF</span>
                  <span className="text-xs text-gray-500">Printable format</span>
                </div>
              </DropdownMenuItem>
              <DropdownMenuItem 
                onClick={() => handleExport("html")}
                className="hover:bg-blue-50 transition-colors"
              >
                <Globe className="mr-2 h-4 w-4 text-blue-600" />
                <div className="flex flex-col">
                  <span>Export as HTML</span>
                  <span className="text-xs text-gray-500">Web format</span>
                </div>
              </DropdownMenuItem>
              {selectedInquiries.length > 0 && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    onClick={() => handleBulkExport(selectedInquiries)}
                    className="hover:bg-purple-50 transition-colors"
                  >
                    <Archive className="mr-2 h-4 w-4 text-purple-600" />
                    <div className="flex flex-col">
                      <span>Bulk Export Selected</span>
                      <span className="text-xs text-gray-500">{selectedInquiries.length} selected items</span>
                    </div>
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Enhanced Bulk Actions Bar */}
      {selectedInquiries.length > 0 && (
        <div className="bg-gradient-to-r from-blue-50 to-blue-100/50 border border-blue-200 rounded-lg p-4 flex items-center justify-between shadow-sm animate-in slide-in-from-top-2 duration-200">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <div className="relative">
                <CheckCircle2 className="h-5 w-5 text-blue-600" />
                <div className="absolute -top-1 -right-1 w-3 h-3 bg-blue-600 rounded-full flex items-center justify-center">
                  <span className="text-xs text-white font-bold">{selectedInquiries.length}</span>
                </div>
              </div>
              <span className="font-medium text-blue-900">
                {selectedInquiries.length} inquir{selectedInquiries.length === 1 ? 'y' : 'ies'} selected
              </span>
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setSelectedInquiries([]);
                setBulkActionMode(false);
              }}
              className="h-7 px-2 text-blue-700 hover:bg-blue-100 transition-all duration-150"
            >
              <X className="h-3 w-3 mr-1" />
              Clear selection
            </Button>
          </div>
          
          <div className="flex items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="h-8 bg-white border-blue-200 text-blue-700 hover:bg-blue-50 shadow-sm hover:shadow transition-all duration-150"
                >
                  <Tag className="h-4 w-4 mr-2" />
                  Update Status
                  <ChevronDown className="h-3 w-3 ml-1 transition-transform group-data-[state=open]:rotate-180" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="border-gray-200 shadow-lg">
                <DropdownMenuItem 
                  onClick={() => handleBulkStatusUpdate("pending")}
                  className="hover:bg-yellow-50 transition-colors"
                >
                  <Clock className="mr-2 h-4 w-4 text-yellow-500" />
                  <span>Mark as Pending</span>
                  <span className="ml-auto text-xs text-gray-500">({selectedInquiries.length})</span>
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => handleBulkStatusUpdate("processing")}
                  className="hover:bg-blue-50 transition-colors"
                >
                  <ArrowUp className="mr-2 h-4 w-4 text-blue-500" />
                  <span>Mark as Processing</span>
                  <span className="ml-auto text-xs text-gray-500">({selectedInquiries.length})</span>
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => handleBulkStatusUpdate("reviewed")}
                  className="hover:bg-purple-50 transition-colors"
                >
                  <Eye className="mr-2 h-4 w-4 text-purple-500" />
                  <span>Mark as Reviewed</span>
                  <span className="ml-auto text-xs text-gray-500">({selectedInquiries.length})</span>
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => handleBulkStatusUpdate("responded")}
                  className="hover:bg-green-50 transition-colors"
                >
                  <Mail className="mr-2 h-4 w-4 text-green-500" />
                  <span>Mark as Responded</span>
                  <span className="ml-auto text-xs text-gray-500">({selectedInquiries.length})</span>
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => handleBulkStatusUpdate("closed")}
                  className="hover:bg-gray-50 transition-colors"
                >
                  <Archive className="mr-2 h-4 w-4 text-gray-500" />
                  <span>Mark as Closed</span>
                  <span className="ml-auto text-xs text-gray-500">({selectedInquiries.length})</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <Button
              variant="outline"
              size="sm"
              onClick={() => handleBulkExport(selectedInquiries)}
              className="h-8 bg-white border-blue-200 text-blue-700 hover:bg-blue-50 shadow-sm hover:shadow transition-all duration-150"
            >
              <Download className="h-4 w-4 mr-2" />
              Export Selected ({selectedInquiries.length})
            </Button>

            <Button
              variant="outline"
              size="sm"
              onClick={handleBulkDelete}
              className="h-8 bg-white border-red-200 text-red-700 hover:bg-red-50 shadow-sm hover:shadow transition-all duration-150"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete ({selectedInquiries.length})
            </Button>
          </div>
        </div>
      )}

      {/* Table */}
      <Card className="border border-gray-200 shadow-sm">
        <div className="rounded-lg border border-gray-200 overflow-hidden">
          {/* Mobile-friendly table wrapper with horizontal scroll */}
          <div className="overflow-x-auto">
            <Table className="min-w-[800px]">
            <TableHeader>
              <TableRow className="bg-gray-50/50">
                <TableHead className="w-[40px]">
                  <Checkbox
                    checked={selectedInquiries.length === sortedInquiries.length && sortedInquiries.length > 0}
                    onCheckedChange={toggleSelectAll}
                    aria-label="Select all"
                  />
                </TableHead>
                <TableHead className="w-[100px]">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("id")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    ID
                    {sortField === "id" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("name")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    <User className="mr-1 h-4 w-4" />
                    Name
                    {sortField === "name" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="hidden md:table-cell">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("email")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    <Mail className="mr-1 h-4 w-4" />
                    Email
                    {sortField === "email" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="hidden lg:table-cell">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("company")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    <Building className="mr-1 h-4 w-4" />
                    Company
                    {sortField === "company" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="hidden lg:table-cell font-medium text-gray-700">Services</TableHead>
                <TableHead>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("status")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    Status
                    {sortField === "status" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="hidden sm:table-cell">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("priority")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    Priority
                    {sortField === "priority" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="hidden sm:table-cell">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleSort("created_at")}
                    className="h-8 px-2 font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  >
                    <Calendar className="mr-1 h-4 w-4" />
                    Date
                    {sortField === "created_at" ? (
                      sortDirection === "asc" ? <ArrowUp className="ml-1 h-4 w-4" /> : <ArrowDown className="ml-1 h-4 w-4" />
                    ) : (
                      <ArrowUpDown className="ml-1 h-4 w-4 opacity-50" />
                    )}
                  </Button>
                </TableHead>
                <TableHead className="w-[120px] font-medium text-gray-700">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedInquiries.length === 0 ? (
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
                sortedInquiries.map((inquiry) => (
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
        </div>
        
        {/* Enhanced Pagination and Stats */}
        {sortedInquiries.length > 0 && (
          <div className="flex items-center justify-between px-6 py-4 border-t border-gray-200 bg-gray-50/30">
            <div className="text-sm text-muted-foreground">
              Showing <strong>{sortedInquiries.length}</strong> of <strong>{inquiries.length}</strong> inquiries
              {getActiveFilterCount() > 0 && (
                <span className="ml-2 text-amber-600">
                  • <strong>{getActiveFilterCount()}</strong> filter{getActiveFilterCount() === 1 ? '' : 's'} applied
                </span>
              )}
              {selectedInquiries.length > 0 && (
                <span className="ml-2 text-blue-600">
                  • <strong>{selectedInquiries.length}</strong> selected
                </span>
              )}
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" disabled className="bg-white">
                Previous
              </Button>
              <span className="px-3 py-1 text-sm text-gray-600">Page 1 of 1</span>
              <Button variant="outline" size="sm" disabled className="bg-white">
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