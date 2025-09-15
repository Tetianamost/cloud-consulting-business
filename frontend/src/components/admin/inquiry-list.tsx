import { useState, useEffect } from "react"
import {
  ArrowUpDown,
  Download,
  MoreHorizontal,
  Search,
  SlidersHorizontal,
  Archive,
  FileText,
  Globe,
  Eye,
} from "lucide-react"

import { Badge } from "../ui/Badge"
import { Button } from "../ui/Button"
import { Card } from "../ui/Card"
import { Checkbox } from "../ui/Checkbox"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/Dropdown-menu"
import { Input } from "../ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/Select"
import { Table as UITable, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/Table"
import apiService, { Inquiry } from '../../services/api'

// Loading and error states
const LoadingState = () => (
  <div className="flex justify-center items-center py-8">
    <div className="text-muted-foreground">Loading inquiries...</div>
  </div>
)

const ErrorState = ({ error, onRetry }: { error: string; onRetry: () => void }) => (
  <div className="flex flex-col items-center py-8 text-center">
    <p className="text-destructive mb-4">Failed to load inquiries</p>
    <p className="text-sm text-muted-foreground mb-4">{error}</p>
    <Button onClick={onRetry}>Try Again</Button>
  </div>
)

const handleExport = (format: string) => {
  const timestamp = new Date().toISOString().split("T")[0]
  const fileName = `inquiries_export_${timestamp}.${format}`

  console.log(`Exporting inquiries as ${fileName}`)
  // In real implementation, would generate and download the file

  // Simulate download
  const link = document.createElement("a")
  link.href = "#" // Would be actual file URL
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const handleBulkExport = (selectedInquiries: string[]) => {
  if (selectedInquiries.length === 0) {
    alert("Please select inquiries to export")
    return
  }

  const timestamp = new Date().toISOString().split("T")[0]
  const fileName = `selected_inquiries_${selectedInquiries.length}_items_${timestamp}.zip`

  console.log(`Bulk exporting ${selectedInquiries.length} inquiries as ${fileName}`)
  // In real implementation, would create ZIP with multiple formats
}

const handleGenerateReport = async (inquiry: any) => {
  console.log(`Generating Bedrock report for inquiry ${inquiry.id}`)
  // In real implementation, would call Bedrock API
  alert(`Generating AI analysis report for ${inquiry.subject}. You'll be notified when it's ready.`)
}

export function InquiryList() {
  const [inquiries, setInquiries] = useState<Inquiry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const [priorityFilter, setPriorityFilter] = useState("all")
  const [selectedInquiries, setSelectedInquiries] = useState<string[]>([])
  const [previewInquiry, setPreviewInquiry] = useState<Inquiry | null>(null)
  const [isPreviewOpen, setIsPreviewOpen] = useState(false)

  useEffect(() => {
    fetchInquiries()
  }, [])

  const fetchInquiries = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await apiService.listInquiries({ limit: 100 })
      setInquiries(response.data || [])
    } catch (err: any) {
      console.error('Failed to fetch inquiries:', err)
      setError(err.message || 'Failed to load inquiries')
      setInquiries([]) // Ensure inquiries is always an array
    } finally {
      setLoading(false)
    }
  }

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
      inquiry.id.toLowerCase().includes(searchQuery.toLowerCase())

    const matchesStatus = statusFilter === "all" || inquiry.status === statusFilter
    const matchesPriority = priorityFilter === "all" || inquiry.priority === priorityFilter

    return matchesSearch && matchesStatus && matchesPriority
  })

  // Handle checkbox selection
  const toggleSelection = (id: string) => {
    setSelectedInquiries((prev) => (prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]))
  }

  // Handle select all
  const toggleSelectAll = () => {
    if (selectedInquiries.length === filteredInquiries.length) {
      setSelectedInquiries([])
    } else {
      setSelectedInquiries(filteredInquiries.map((i) => i.id))
    }
  }

  // Format date to readable string
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date)
  }

  // Get status badge color
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return "bg-yellow-100 text-yellow-800 hover:bg-yellow-100/80"
      case "processing":
        return "bg-blue-100 text-blue-800 hover:bg-blue-100/80"
      case "reviewed":
        return "bg-purple-100 text-purple-800 hover:bg-purple-100/80"
      case "responded":
        return "bg-green-100 text-green-800 hover:bg-green-100/80"
      case "closed":
        return "bg-gray-100 text-gray-800 hover:bg-gray-100/80"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-100/80"
    }
  }

  // Get priority badge color
  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case "high":
      case "urgent":
        return "bg-red-100 text-red-800 hover:bg-red-100/80"
      case "medium":
        return "bg-amber-100 text-amber-800 hover:bg-amber-100/80"
      case "low":
        return "bg-green-100 text-green-800 hover:bg-green-100/80"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-100/80"
    }
  }

  const handleDownloadReport = async (inquiryId: string, format: 'pdf' | 'html') => {
    try {
      const blob = await apiService.downloadReport(inquiryId, format)
      
      // Create a download link
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.style.display = 'none'
      a.href = url
      a.download = `report-${inquiryId}.${format}`
      document.body.appendChild(a)
      a.click()
      
      // Clean up
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      console.error(`Failed to download ${format} report:`, err)
      alert(`Failed to download ${format} report. Please try again later.`)
    }
  }

  if (loading) {
    return <LoadingState />
  }

  if (error) {
    return <ErrorState error={error} onRetry={fetchInquiries} />
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex w-full items-center gap-2 sm:max-w-sm">
          <Search className="h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search inquiries..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="h-9 sm:max-w-sm"
          />
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
                <UITable className="mr-2 h-4 w-4" />
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
                Bulk Export Selected
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <Card>
        <div className="rounded-md border">
          <UITable>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[40px]">
                  <Checkbox
                    checked={selectedInquiries.length === filteredInquiries.length && filteredInquiries.length > 0}
                    onCheckedChange={toggleSelectAll}
                    aria-label="Select all"
                  />
                </TableHead>
                <TableHead className="w-[100px]">
                  <div className="flex items-center">
                    ID
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead>
                  <div className="flex items-center">
                    Name
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead className="hidden md:table-cell">Email</TableHead>
                <TableHead className="hidden lg:table-cell">Company</TableHead>
                <TableHead className="hidden lg:table-cell">Services</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="hidden sm:table-cell">Priority</TableHead>
                <TableHead className="hidden sm:table-cell">
                  <div className="flex items-center">
                    Date
                    <ArrowUpDown className="ml-1 h-4 w-4" />
                  </div>
                </TableHead>
                <TableHead className="w-[50px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredInquiries.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={9} className="h-24 text-center">
                    No inquiries found.
                  </TableCell>
                </TableRow>
              ) : (
                filteredInquiries.map((inquiry) => (
                  <TableRow key={inquiry.id}>
                    <TableCell>
                      <Checkbox
                        checked={selectedInquiries.includes(inquiry.id)}
                        onCheckedChange={() => toggleSelection(inquiry.id)}
                        aria-label={`Select ${inquiry.id}`}
                      />
                    </TableCell>
                    <TableCell className="font-medium">{inquiry.id}</TableCell>
                    <TableCell className="font-medium">{inquiry.name}</TableCell>
                    <TableCell className="hidden md:table-cell">{inquiry.email}</TableCell>
                    <TableCell className="hidden lg:table-cell">{inquiry.company || '—'}</TableCell>
                    <TableCell className="hidden lg:table-cell">{(inquiry.services || []).join(', ')}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className={getStatusColor(inquiry.status)}>
                        {inquiry.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">
                      <Badge variant="outline" className={getPriorityColor(inquiry.priority)}>
                        {inquiry.priority}
                      </Badge>
                    </TableCell>
                    <TableCell className="hidden sm:table-cell">{formatDate(inquiry.created_at)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => {
                            setPreviewInquiry(inquiry)
                            setIsPreviewOpen(true)
                          }}
                          title="View Report"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontal className="h-4 w-4" />
                              <span className="sr-only">Open menu</span>
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => handleDownloadReport(inquiry.id, 'pdf')}>
                              <FileText className="mr-2 h-4 w-4" />
                              Download PDF Report
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleDownloadReport(inquiry.id, 'html')}>
                              <Globe className="mr-2 h-4 w-4" />
                              Download HTML Report
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={() => handleGenerateReport(inquiry)}>
                              <FileText className="mr-2 h-4 w-4" />
                              Generate New Report
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                        {inquiry.reports && inquiry.reports.length > 0 && (
                          <Badge variant="secondary" className="ml-2 text-xs">
                            {inquiry.reports.length} Report{inquiry.reports.length > 1 ? 's' : ''}
                          </Badge>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </UITable>
        </div>
        <div className="flex items-center justify-between px-4 py-4">
          <div className="text-sm text-muted-foreground">
            Showing <strong>{filteredInquiries.length}</strong> of <strong>{inquiries.length}</strong> inquiries
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" disabled>
              Previous
            </Button>
            <Button variant="outline" size="sm">
              Next
            </Button>
          </div>
        </div>
      </Card>
      {/* Report Preview Modal */}
      {previewInquiry && isPreviewOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="bg-white rounded-xl shadow-2xl max-w-lg w-full p-0 relative">
            <div className="flex items-center justify-between px-6 pt-6 pb-2 border-b border-gray-200">
              <div>
                <h2 className="text-xl font-bold text-gray-900 mb-1">Inquiry Details</h2>
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant="outline" className="capitalize text-xs px-2 py-1">
                    {previewInquiry.status}
                  </Badge>
                  <Badge variant="secondary" className="capitalize text-xs px-2 py-1">
                    {previewInquiry.priority}
                  </Badge>
                  <span className="text-xs text-gray-500">
                    {formatDate(previewInquiry.created_at)}
                  </span>
                </div>
              </div>
              <button
                className="text-gray-400 hover:text-gray-600 text-2xl"
                onClick={() => setIsPreviewOpen(false)}
                aria-label="Close"
              >
                ×
              </button>
            </div>
            <div className="px-6 py-4 space-y-3">
              <div className="flex items-center gap-3">
                <span className="font-semibold text-gray-700 w-24">ID:</span>
                <span className="font-mono text-xs text-gray-800">{previewInquiry.id}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="font-semibold text-gray-700 w-24">Name:</span>
                <span className="text-gray-900">{previewInquiry.name}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="font-semibold text-gray-700 w-24">Email:</span>
                <span className="text-blue-700 underline">{previewInquiry.email}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="font-semibold text-gray-700 w-24">Company:</span>
                <span className="text-gray-900">{previewInquiry.company || "—"}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="font-semibold text-gray-700 w-24">Services:</span>
                <span className="text-gray-900">{(previewInquiry.services || []).join(", ")}</span>
              </div>
              {/* Link to generated report if available */}
              {previewInquiry.reports && previewInquiry.reports.length > 0 && (
                <div className="flex items-center gap-3">
                  <span className="font-semibold text-gray-700 w-24">AI Report:</span>
                  <a
                    href={previewInquiry.reports && previewInquiry.reports.length > 0 ? `/admin/reports/${previewInquiry.reports[0].id}` : "#"}
                    className="text-blue-600 hover:underline font-medium flex items-center gap-1"
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={e => {
                      if (previewInquiry.reports && previewInquiry.reports.length > 0) {
                        e.preventDefault();
                        window.open(`/admin/reports/${previewInquiry.reports[0].id}`, "_blank", "noopener,noreferrer");
                      }
                    }}
                  >
                    View Report
                  </a>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
