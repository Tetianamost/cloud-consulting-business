import { useState, useEffect } from "react"
import {
  ArrowDown,
  ArrowUp,
  BarChart3,
  Mail,
  MessageSquare,
  Users,
  Archive,
  ChevronDown,
  Eye,
  FileText,
  Globe,
  Table,
  Download,
} from "lucide-react"
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts"

import { Button } from "../ui/Button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/Card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/Tabs"
import { Checkbox } from "../ui/Checkbox"
import { Badge } from "../ui/Badge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/Dropdown-menu"
import apiService, { SystemMetrics, Inquiry } from '../../services/api'

// Sample data for demonstration
const overviewData = [
  { name: "Jan", inquiries: 65, emails: 120, users: 45 },
  { name: "Feb", inquiries: 59, emails: 80, users: 49 },
  { name: "Mar", inquiries: 80, emails: 150, users: 60 },
  { name: "Apr", inquiries: 81, emails: 170, users: 65 },
  { name: "May", inquiries: 56, emails: 140, users: 78 },
  { name: "Jun", inquiries: 55, emails: 130, users: 88 },
  { name: "Jul", inquiries: 40, emails: 100, users: 90 },
]

const emailData = [
  { name: "Mon", delivered: 240, failed: 5, opened: 180 },
  { name: "Tue", delivered: 300, failed: 8, opened: 230 },
  { name: "Wed", delivered: 320, failed: 10, opened: 250 },
  { name: "Thu", delivered: 280, failed: 7, opened: 220 },
  { name: "Fri", delivered: 250, failed: 6, opened: 190 },
  { name: "Sat", delivered: 150, failed: 3, opened: 100 },
  { name: "Sun", delivered: 100, failed: 2, opened: 70 },
]

const inquiryData = [
  { name: "General", value: 40 },
  { name: "Support", value: 30 },
  { name: "Sales", value: 20 },
  { name: "Feedback", value: 10 },
]

export function MetricsDashboard() {
  const [timeRange, setTimeRange] = useState("7d")
  const [selectedReports, setSelectedReports] = useState<string[]>([])
  const [showPreview, setShowPreview] = useState(false)
  const [previewReport, setPreviewReport] = useState<any>(null)
  const [showReportGenerator, setShowReportGenerator] = useState(false)
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null)
  const [inquiries, setInquiries] = useState<Inquiry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      setLoading(true)
      setError(null)
      
      // Fetch metrics and inquiries in parallel
      const [metricsResponse, inquiriesResponse] = await Promise.all([
        apiService.getSystemMetrics(),
        apiService.listInquiries({ limit: 100 })
      ])
      
      setMetrics(metricsResponse.data)
      setInquiries(inquiriesResponse.data)
    } catch (err: any) {
      console.error('Failed to fetch dashboard data:', err)
      setError(err.message || 'Failed to load dashboard data')
    } finally {
      setLoading(false)
    }
  }

  const reports = [
    {
      id: "INQ-RPT-2025-07-001",
      title: "Cloud Migration Quote Analysis - TechCorp Inc",
      type: "Quote Analysis",
      size: "1.8MB",
      status: "Ready",
      generatedAt: "2025-07-26T10:30:00",
      expiresAt: "2025-08-26T10:30:00",
      inquiryId: "INQ-1001",
      customerEmail: "john.smith@techcorp.com",
      serviceType: "AWS Migration",
      priority: "High",
      bedrockModel: "Claude-3-Sonnet",
      confidence: 0.92,
      recommendedActions: ["Schedule technical call", "Prepare cost estimate", "Review compliance requirements"],
    },
    {
      id: "INQ-RPT-2025-07-002",
      title: "Multi-Cloud Strategy Inquiry - StartupXYZ",
      type: "Strategy Analysis",
      size: "2.1MB",
      status: "Ready",
      generatedAt: "2025-07-25T14:45:00",
      expiresAt: "2025-08-25T14:45:00",
      inquiryId: "INQ-1002",
      customerEmail: "cto@startupxyz.com",
      serviceType: "Multi-Cloud Architecture",
      priority: "Medium",
      bedrockModel: "Claude-3-Haiku",
      confidence: 0.87,
      recommendedActions: ["Send architecture whitepaper", "Schedule discovery session"],
    },
    {
      id: "INQ-RPT-2025-07-003",
      title: "Data Analytics Platform Quote - RetailCorp",
      type: "Quote Analysis",
      size: "2.4MB",
      status: "Ready",
      generatedAt: "2025-07-24T09:15:00",
      expiresAt: "2025-08-24T09:15:00",
      inquiryId: "INQ-1003",
      customerEmail: "data-team@retailcorp.com",
      serviceType: "Analytics Platform",
      priority: "High",
      bedrockModel: "Claude-3-Sonnet",
      confidence: 0.94,
      recommendedActions: ["Prepare demo environment", "Schedule technical deep-dive", "Review data governance"],
    },
    {
      id: "INQ-RPT-2025-07-004",
      title: "Security Compliance Assessment - FinanceFirst",
      type: "Compliance Analysis",
      size: "3.2MB",
      status: "Generating",
      generatedAt: "2025-07-26T08:15:00",
      expiresAt: null,
      inquiryId: "INQ-1004",
      customerEmail: "security@financefirst.com",
      serviceType: "Security & Compliance",
      priority: "Critical",
      bedrockModel: "Claude-3-Sonnet",
      confidence: null,
      recommendedActions: [],
    },
  ]

  const [reportType, setReportType] = useState("Monthly Performance")
  const [startDate, setStartDate] = useState("")
  const [endDate, setEndDate] = useState("")

  // Add these handler functions
  const toggleReportSelection = (reportId: string) => {
    setSelectedReports((prev) => (prev.includes(reportId) ? prev.filter((id) => id !== reportId) : [...prev, reportId]))
  }

  const handleReportPreview = (report: any) => {
    setPreviewReport(report)
    setShowPreview(true)
  }

  const handleDownload = async (report: any, format: string) => {
    // Simulate download with proper file naming
    const fileName = generateFileName(report, format)

    // In a real implementation, this would trigger the actual download
    console.log(`Downloading ${fileName}`)

    // Simulate file download
    const link = document.createElement("a")
    link.href = "#" // In real app, this would be the file URL
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const handleBulkDownload = async () => {
    if (selectedReports.length === 0) {
      alert("Please select reports to download")
      return
    }

    const selectedReportData = reports.filter((r) => selectedReports.includes(r.id))
    const zipFileName = `bulk_reports_${new Date().toISOString().split("T")[0]}.zip`

    console.log(`Bulk downloading ${selectedReports.length} reports as ${zipFileName}`)
    // In real implementation, would create and download ZIP file
  }

  const handleGenerateReport = (reportConfig: any) => {
    console.log("Generating new report:", reportConfig)
    setShowReportGenerator(false)
  }

  const generateFileName = (report: any, format: string) => {
    const date = new Date(report.generatedAt).toISOString().split("T")[0]
    const sanitizedTitle = report.title.toLowerCase().replace(/[^a-z0-9]/g, "_")
    return `${report.id}_${sanitizedTitle}_${date}.${format}`
  }

  const getReportStatusColor = (status: string) => {
    switch (status) {
      case "Ready":
        return "bg-green-100 text-green-800 hover:bg-green-100/80"
      case "Generating":
        return "bg-amber-100 text-amber-800 hover:bg-amber-100/80"
      case "Failed":
        return "bg-red-100 text-red-800 hover:bg-red-100/80"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-100/80"
    }
  }

  const ReportPreviewModal = ({
    report,
    isOpen,
    onClose,
    onDownload,
  }: { report: any; isOpen: boolean; onClose: () => void; onDownload: (report: any, format: string) => void }) => {
    if (!isOpen) return null

    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
        <div className="bg-white p-6 rounded-md shadow-lg">
          <h2 className="text-lg font-semibold mb-4">Report Preview: {report?.title}</h2>
          <p>Report ID: {report?.id}</p>
          <p>Type: {report?.type}</p>
          <p>Size: {report?.size}</p>
          <p>Status: {report?.status}</p>
          <p>Generated At: {report?.generatedAt}</p>
          {report?.expiresAt && <p>Expires At: {report?.expiresAt}</p>}
          {report?.customerEmail && <p>Customer Email: {report.customerEmail}</p>}
          {report?.serviceType && <p>Service Type: {report.serviceType}</p>}
          {report?.bedrockModel && <p>Bedrock Model: {report.bedrockModel}</p>}
          {report?.confidence && <p>Confidence: {report.confidence}</p>}
          {report?.recommendedActions && (
            <div>
              <p>Recommended Actions:</p>
              <div className="flex gap-2">
                {report.recommendedActions.map((action: string) => (
                  <Badge key={action}>{action}</Badge>
                ))}
              </div>
            </div>
          )}

          <div className="mt-4 flex justify-end gap-2">
            <Button variant="outline" onClick={onClose}>
              Close
            </Button>
            <Button onClick={() => onDownload(report, "pdf")}>Download PDF</Button>
          </div>
        </div>
      </div>
    )
  }

  const ReportGeneratorModal = ({
    isOpen,
    onClose,
    onGenerate,
  }: { isOpen: boolean; onClose: () => void; onGenerate: (reportConfig: any) => void }) => {
    if (!isOpen) return null

    const handleSubmit = () => {
      const reportConfig = {
        type: reportType,
        startDate: startDate,
        endDate: endDate,
      }
      onGenerate(reportConfig)
    }

    const reportTypes = [
      { value: "quote_analysis", label: "Quote Analysis", description: "AI analysis of customer quote requests" },
      {
        value: "strategy_analysis",
        label: "Strategy Analysis",
        description: "Strategic recommendations for cloud adoption",
      },
      { value: "compliance_analysis", label: "Compliance Analysis", description: "Security and compliance assessment" },
      {
        value: "technical_analysis",
        label: "Technical Analysis",
        description: "Technical feasibility and architecture review",
      },
    ]

    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
        <div className="bg-white p-6 rounded-md shadow-lg">
          <h2 className="text-lg font-semibold mb-4">Generate New Report</h2>

          <div className="mb-4">
            <label htmlFor="reportType" className="block text-sm font-medium text-gray-700">
              Report Type
            </label>
            <select
              id="reportType"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              value={reportType}
              onChange={(e) => setReportType(e.target.value)}
            >
              {reportTypes.map((type) => (
                <option key={type.value} value={type.value}>
                  {type.label}
                </option>
              ))}
            </select>
          </div>

          <div className="mb-4">
            <label htmlFor="startDate" className="block text-sm font-medium text-gray-700">
              Start Date
            </label>
            <input
              type="date"
              id="startDate"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
            />
          </div>

          <div className="mb-4">
            <label htmlFor="endDate" className="block text-sm font-medium text-gray-700">
              End Date
            </label>
            <input
              type="date"
              id="endDate"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <Button variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button onClick={handleSubmit}>Generate</Button>
          </div>
        </div>
      </div>
    )
  }

  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return "N/A"
    const date = new Date(dateString)
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "Critical":
        return "bg-red-500 text-white"
      case "High":
        return "bg-orange-500 text-white"
      case "Medium":
        return "bg-yellow-500 text-gray-800"
      case "Low":
        return "bg-green-500 text-white"
      default:
        return "bg-gray-500 text-white"
    }
  }

  return (
    <div className="space-y-4">
      <Tabs defaultValue="overview" className="space-y-4">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="analytics">Analytics</TabsTrigger>
            <TabsTrigger value="reports">Reports</TabsTrigger>
          </TabsList>
          <div className="flex items-center gap-2">
            <TabsList>
              <TabsTrigger
                value="1d"
                onClick={() => setTimeRange("1d")}
                className={timeRange === "1d" ? "bg-primary text-primary-foreground" : ""}
              >
                1d
              </TabsTrigger>
              <TabsTrigger
                value="7d"
                onClick={() => setTimeRange("7d")}
                className={timeRange === "7d" ? "bg-primary text-primary-foreground" : ""}
              >
                7d
              </TabsTrigger>
              <TabsTrigger
                value="30d"
                onClick={() => setTimeRange("30d")}
                className={timeRange === "30d" ? "bg-primary text-primary-foreground" : ""}
              >
                30d
              </TabsTrigger>
              <TabsTrigger
                value="90d"
                onClick={() => setTimeRange("90d")}
                className={timeRange === "90d" ? "bg-primary text-primary-foreground" : ""}
              >
                90d
              </TabsTrigger>
            </TabsList>
          </div>
        </div>
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Inquiries</CardTitle>
                <MessageSquare className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.total_inquiries || 0}</div>
                <p className="text-xs text-muted-foreground">
                  <span className="text-emerald-500 flex items-center">
                    <ArrowUp className="mr-1 h-4 w-4" />
                    12%
                  </span>{" "}
                  from last month
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Emails Sent</CardTitle>
                <Mail className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.emails_sent || 0}</div>
                <p className="text-xs text-muted-foreground">
                  <span className="text-emerald-500 flex items-center">
                    <ArrowUp className="mr-1 h-4 w-4" />
                    8%
                  </span>{" "}
                  from last month
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Reports Generated</CardTitle>
                <BarChart3 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.reports_generated || 0}</div>
                <p className="text-xs text-muted-foreground">
                  <span className="text-emerald-500 flex items-center">
                    <ArrowUp className="mr-1 h-4 w-4" />
                    20%
                  </span>{" "}
                  from last month
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Email Delivery Rate</CardTitle>
                <Mail className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.email_delivery_rate?.toFixed(1) || 0}%</div>
                <p className="text-xs text-muted-foreground">
                  <span className="text-rose-500 flex items-center">
                    <ArrowDown className="mr-1 h-4 w-4" />
                    3%
                  </span>{" "}
                  from last month
                </p>
              </CardContent>
            </Card>
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
            <Card className="col-span-4">
              <CardHeader>
                <CardTitle>Overview</CardTitle>
                <CardDescription>System activity for the selected period.</CardDescription>
              </CardHeader>
              <CardContent className="pl-2">
                <ResponsiveContainer width="100%" height={350}>
                  <LineChart data={overviewData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Line type="monotone" dataKey="inquiries" stroke="#8884d8" activeDot={{ r: 8 }} />
                    <Line type="monotone" dataKey="emails" stroke="#82ca9d" />
                    <Line type="monotone" dataKey="users" stroke="#ffc658" />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Inquiry Types</CardTitle>
                <CardDescription>Distribution of inquiry categories.</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={350}>
                  <BarChart data={inquiryData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="value" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
        <TabsContent value="analytics" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Email Performance</CardTitle>
              <CardDescription>Email delivery and engagement metrics.</CardDescription>
            </CardHeader>
            <CardContent className="pl-2">
              <ResponsiveContainer width="100%" height={350}>
                <AreaChart data={emailData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Area type="monotone" dataKey="delivered" stackId="1" stroke="#8884d8" fill="#8884d8" />
                  <Area type="monotone" dataKey="opened" stackId="2" stroke="#82ca9d" fill="#82ca9d" />
                  <Area type="monotone" dataKey="failed" stackId="3" stroke="#ffc658" fill="#ffc658" />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="reports" className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-medium">System Reports</h3>
              <p className="text-sm text-muted-foreground">Generate and download detailed system reports</p>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={() => handleBulkDownload()}>
                <Download className="mr-2 h-4 w-4" />
                Bulk Download
              </Button>
              <Button size="sm" onClick={() => setShowReportGenerator(true)}>
                Generate New Report
              </Button>
            </div>
          </div>

          <Card>
            <CardContent className="p-6">
              <div className="space-y-4">
                {reports.map((report) => (
                  <div key={report.id} className="flex items-center justify-between border-b pb-4 last:border-b-0">
                    <div className="flex items-center gap-3">
                      <Checkbox
                        checked={selectedReports.includes(report.id)}
                        onCheckedChange={() => toggleReportSelection(report.id)}
                      />
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <p className="font-medium">{report.title}</p>
                          <Badge variant="outline" className={getReportStatusColor(report.status)}>
                            {report.status}
                          </Badge>
                          {report.priority && (
                            <Badge className={getPriorityColor(report.priority)}>{report.priority}</Badge>
                          )}
                        </div>
                        <div className="flex items-center gap-4 text-sm text-muted-foreground">
                          <span>
                            {report.type} • {report.size}
                          </span>
                          <span>Generated: {formatDate(report.generatedAt)}</span>
                          {report.expiresAt && <span>Expires: {formatDate(report.expiresAt)}</span>}
                          {report.customerEmail && <span>Customer Email: {report.customerEmail}</span>}
                          {report.serviceType && <span>Service Type: {report.serviceType}</span>}
                          {report.bedrockModel && <span>Bedrock Model: {report.bedrockModel}</span>}
                          {report.confidence && (
                            <span>Confidence: {report.confidence ? report.confidence.toFixed(2) : "N/A"}</span>
                          )}
                          {report.recommendedActions && report.recommendedActions.length > 0 && (
                            <div className="flex items-center gap-2">
                              Recommended Actions:
                              {report.recommendedActions.map((action: string) => (
                                <Badge key={action}>{action}</Badge>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleReportPreview(report)}
                        disabled={report.status !== "Ready"}
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="outline" size="sm" disabled={report.status !== "Ready"}>
                            <Download className="mr-2 h-4 w-4" />
                            Download
                            <ChevronDown className="ml-1 h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleDownload(report, "pdf")}>
                            <FileText className="mr-2 h-4 w-4" />
                            Download PDF
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleDownload(report, "html")}>
                            <Globe className="mr-2 h-4 w-4" />
                            Download HTML
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleDownload(report, "csv")}>
                            <Table className="mr-2 h-4 w-4" />
                            Download CSV
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem onClick={() => handleDownload(report, "zip")}>
                            <Archive className="mr-2 h-4 w-4" />
                            Download All Formats (ZIP)
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Report Preview Modal */}
          <ReportPreviewModal
            report={previewReport}
            isOpen={showPreview}
            onClose={() => setShowPreview(false)}
            onDownload={handleDownload}
          />

          {/* Report Generator Modal */}
          <ReportGeneratorModal
            isOpen={showReportGenerator}
            onClose={() => setShowReportGenerator(false)}
            onGenerate={handleGenerateReport}
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}
