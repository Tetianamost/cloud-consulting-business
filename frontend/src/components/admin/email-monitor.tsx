import { useState, useEffect } from "react"
import {
  AlertCircle,
  ArrowUpDown,
  CheckCircle2,
  Clock,
  Download,
  Filter,
  MailCheck,
  MailX,
  Search,
  XCircle,
  Archive,
  FileText,
  Globe,
} from "lucide-react"

import { Badge } from "../ui/badge"
import { Button } from "../ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card"
import { Input } from "../ui/input"
import { Progress } from "../ui/progress"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select"
import { Table as UITable, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import apiService, { Inquiry, EmailStatus } from '../../services/api'

// Sample data for demonstration
const emailStats = {
  total: 1390,
  delivered: 1320,
  opened: 890,
  clicked: 456,
  failed: 70,
  bounced: 45,
  spam: 25,
  deliveryRate: 94.9,
  openRate: 67.4,
  clickRate: 34.5,
}

const emailEvents = [
  {
    id: "EMAIL-1001",
    recipient: "john.doe@example.com",
    subject: "Your Monthly Newsletter",
    status: "Delivered",
    sentAt: "2025-07-20T10:30:00",
    openedAt: "2025-07-20T11:15:00",
    clickedAt: "2025-07-20T11:20:00",
  },
  {
    id: "EMAIL-1002",
    recipient: "jane.smith@example.com",
    subject: "Order Confirmation #12345",
    status: "Delivered",
    sentAt: "2025-07-20T09:45:00",
    openedAt: "2025-07-20T10:10:00",
    clickedAt: null,
  },
  {
    id: "EMAIL-1003",
    recipient: "mike.wilson@example.com",
    subject: "Your Account Statement",
    status: "Failed",
    sentAt: "2025-07-20T08:30:00",
    openedAt: null,
    clickedAt: null,
    failureReason: "Invalid email address",
  },
  {
    id: "EMAIL-1004",
    recipient: "sarah.johnson@example.com",
    subject: "Special Offer Inside!",
    status: "Bounced",
    sentAt: "2025-07-19T15:20:00",
    openedAt: null,
    clickedAt: null,
    failureReason: "Mailbox full",
  },
  {
    id: "EMAIL-1005",
    recipient: "robert.brown@example.com",
    subject: "Your Subscription Renewal",
    status: "Delivered",
    sentAt: "2025-07-19T14:10:00",
    openedAt: "2025-07-19T16:30:00",
    clickedAt: "2025-07-19T16:35:00",
  },
  {
    id: "EMAIL-1006",
    recipient: "emily.davis@example.com",
    subject: "Important Security Update",
    status: "Delivered",
    sentAt: "2025-07-19T11:45:00",
    openedAt: null,
    clickedAt: null,
  },
  {
    id: "EMAIL-1007",
    recipient: "david.miller@example.com",
    subject: "Your Support Ticket #5678",
    status: "Spam",
    sentAt: "2025-07-18T16:20:00",
    openedAt: null,
    clickedAt: null,
    failureReason: "Marked as spam",
  },
  {
    id: "EMAIL-1008",
    recipient: "lisa.taylor@example.com",
    subject: "Invitation to Webinar",
    status: "Delivered",
    sentAt: "2025-07-18T10:15:00",
    openedAt: "2025-07-18T13:40:00",
    clickedAt: "2025-07-18T13:45:00",
  },
]

export function EmailMonitor() {
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const [timeRange, setTimeRange] = useState("24h")

  // Filter email events based on search and filters
  const filteredEmails = emailEvents.filter((email) => {
    const matchesSearch =
      searchQuery === "" ||
      email.recipient.toLowerCase().includes(searchQuery.toLowerCase()) ||
      email.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
      email.id.toLowerCase().includes(searchQuery.toLowerCase())

    const matchesStatus = statusFilter === "all" || email.status === statusFilter

    return matchesSearch && matchesStatus
  })

  // Format date to readable string
  const formatDate = (dateString: string | null) => {
    if (!dateString) return "—"
    const date = new Date(dateString)
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date)
  }

  // Get status icon
  const getStatusIcon = (status: string) => {
    switch (status) {
      case "Delivered":
        return <CheckCircle2 className="h-4 w-4 text-green-500" />
      case "Failed":
        return <XCircle className="h-4 w-4 text-red-500" />
      case "Bounced":
        return <AlertCircle className="h-4 w-4 text-amber-500" />
      case "Spam":
        return <MailX className="h-4 w-4 text-red-500" />
      default:
        return <Clock className="h-4 w-4 text-gray-500" />
    }
  }

  // Get status badge color
  const getStatusColor = (status: string) => {
    switch (status) {
      case "Delivered":
        return "bg-green-100 text-green-800 hover:bg-green-100/80"
      case "Failed":
        return "bg-red-100 text-red-800 hover:bg-red-100/80"
      case "Bounced":
        return "bg-amber-100 text-amber-800 hover:bg-amber-100/80"
      case "Spam":
        return "bg-red-100 text-red-800 hover:bg-red-100/80"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-100/80"
    }
  }

  const handleEmailExport = (format: string) => {
    const timestamp = new Date().toISOString().split("T")[0]
    const timeRangeLabel = timeRange.replace("_", "-")

    let fileName = ""
    switch (format) {
      case "csv":
        fileName = `email_events_${timeRangeLabel}_${timestamp}.csv`
        break
      case "pdf":
        fileName = `email_analytics_report_${timeRangeLabel}_${timestamp}.pdf`
        break
      case "html":
        fileName = `email_dashboard_${timeRangeLabel}_${timestamp}.html`
        break
      case "zip":
        fileName = `email_data_complete_${timeRangeLabel}_${timestamp}.zip`
        break
      default:
        fileName = `email_export_${timestamp}.${format}`
    }

    console.log(`Exporting email data as ${fileName}`)
    // In real implementation, would generate and download the file

    // Simulate download
    const link = document.createElement("a")
    link.href = "#" // Would be actual file URL
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  return (
    <div className="space-y-4">
      <Tabs defaultValue="overview" className="space-y-4">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="events">Email Events</TabsTrigger>
            <TabsTrigger value="analytics">Analytics</TabsTrigger>
          </TabsList>
          <div className="flex items-center gap-2">
            <Select value={timeRange} onValueChange={setTimeRange}>
              <SelectTrigger className="h-8 w-[120px]">
                <SelectValue placeholder="Time Range" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="24h">Last 24 Hours</SelectItem>
                <SelectItem value="7d">Last 7 Days</SelectItem>
                <SelectItem value="30d">Last 30 Days</SelectItem>
                <SelectItem value="90d">Last 90 Days</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Delivery Rate</CardTitle>
                <MailCheck className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{emailStats.deliveryRate}%</div>
                <Progress value={emailStats.deliveryRate} className="mt-2 h-1" />
                <p className="mt-2 text-xs text-muted-foreground">
                  {emailStats.delivered} of {emailStats.total} emails delivered
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Open Rate</CardTitle>
                <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{emailStats.openRate}%</div>
                <Progress value={emailStats.openRate} className="mt-2 h-1" />
                <p className="mt-2 text-xs text-muted-foreground">
                  {emailStats.opened} of {emailStats.delivered} emails opened
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Click Rate</CardTitle>
                <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{emailStats.clickRate}%</div>
                <Progress value={emailStats.clickRate} className="mt-2 h-1" />
                <p className="mt-2 text-xs text-muted-foreground">
                  {emailStats.clicked} of {emailStats.opened} emails clicked
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Failed Emails</CardTitle>
                <AlertCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{emailStats.failed}</div>
                <div className="mt-2 flex items-center gap-2">
                  <Badge variant="outline" className="bg-red-100 text-red-800">
                    Bounced: {emailStats.bounced}
                  </Badge>
                  <Badge variant="outline" className="bg-amber-100 text-amber-800">
                    Spam: {emailStats.spam}
                  </Badge>
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Email Delivery Status</CardTitle>
              <CardDescription>Overview of email delivery performance for the selected time period.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-8">
                <div>
                  <div className="mb-2 flex items-center justify-between">
                    <div className="text-sm font-medium">Delivered</div>
                    <div className="text-sm text-muted-foreground">{emailStats.deliveryRate}%</div>
                  </div>
                  <Progress value={emailStats.deliveryRate} className="h-2" />
                </div>
                <div>
                  <div className="mb-2 flex items-center justify-between">
                    <div className="text-sm font-medium">Opened</div>
                    <div className="text-sm text-muted-foreground">{emailStats.openRate}%</div>
                  </div>
                  <Progress value={emailStats.openRate} className="h-2" />
                </div>
                <div>
                  <div className="mb-2 flex items-center justify-between">
                    <div className="text-sm font-medium">Clicked</div>
                    <div className="text-sm text-muted-foreground">{emailStats.clickRate}%</div>
                  </div>
                  <Progress value={emailStats.clickRate} className="h-2" />
                </div>
                <div>
                  <div className="mb-2 flex items-center justify-between">
                    <div className="text-sm font-medium">Failed</div>
                    <div className="text-sm text-muted-foreground">
                      {((emailStats.failed / emailStats.total) * 100).toFixed(1)}%
                    </div>
                  </div>
                  <Progress
                    value={(emailStats.failed / emailStats.total) * 100}
                    className="h-2 bg-muted [&>div]:bg-red-500"
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="events" className="space-y-4">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex w-full items-center gap-2 sm:max-w-sm">
              <Search className="h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search emails..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="h-9 sm:max-w-sm"
              />
            </div>
            <div className="flex flex-col gap-2 sm:flex-row">
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="h-9 w-[160px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Filter by status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Statuses</SelectItem>
                  <SelectItem value="Delivered">Delivered</SelectItem>
                  <SelectItem value="Failed">Failed</SelectItem>
                  <SelectItem value="Bounced">Bounced</SelectItem>
                  <SelectItem value="Spam">Spam</SelectItem>
                </SelectContent>
              </Select>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="sm" className="h-9 gap-1 bg-transparent">
                    <Download className="h-4 w-4" />
                    <span className="hidden sm:inline-block">Export</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => handleEmailExport("csv")}>
                    <UITable className="mr-2 h-4 w-4" />
                    Export Email Events (CSV)
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleEmailExport("pdf")}>
                    <FileText className="mr-2 h-4 w-4" />
                    Export Analytics Report (PDF)
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleEmailExport("html")}>
                    <Globe className="mr-2 h-4 w-4" />
                    Export Dashboard (HTML)
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => handleEmailExport("zip")}>
                    <Archive className="mr-2 h-4 w-4" />
                    Export All Formats (ZIP)
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
                    <TableHead className="w-[100px]">
                      <div className="flex items-center">
                        ID
                        <ArrowUpDown className="ml-1 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead>
                      <div className="flex items-center">
                        Recipient
                        <ArrowUpDown className="ml-1 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead className="hidden md:table-cell">Subject</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="hidden sm:table-cell">
                      <div className="flex items-center">
                        Sent
                        <ArrowUpDown className="ml-1 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead className="hidden lg:table-cell">Opened</TableHead>
                    <TableHead className="hidden lg:table-cell">Clicked</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredEmails.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="h-24 text-center">
                        No email events found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredEmails.map((email) => (
                      <TableRow key={email.id}>
                        <TableCell className="font-medium">{email.id}</TableCell>
                        <TableCell>{email.recipient}</TableCell>
                        <TableCell className="hidden md:table-cell max-w-[200px] truncate">{email.subject}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {getStatusIcon(email.status)}
                            <Badge variant="outline" className={getStatusColor(email.status)}>
                              {email.status}
                            </Badge>
                          </div>
                        </TableCell>
                        <TableCell className="hidden sm:table-cell">{formatDate(email.sentAt)}</TableCell>
                        <TableCell className="hidden lg:table-cell">{formatDate(email.openedAt)}</TableCell>
                        <TableCell className="hidden lg:table-cell">{formatDate(email.clickedAt)}</TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </UITable>
            </div>
            <div className="flex items-center justify-between px-4 py-4">
              <div className="text-sm text-muted-foreground">
                Showing <strong>{filteredEmails.length}</strong> of <strong>{emailEvents.length}</strong> email events
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
        </TabsContent>

        <TabsContent value="analytics" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Email Performance Analytics</CardTitle>
              <CardDescription>Detailed analytics for email campaigns and deliverability.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-8">
                <div>
                  <h3 className="mb-4 text-lg font-medium">Delivery by Domain</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">gmail.com</div>
                        <div className="text-sm text-muted-foreground">98.2%</div>
                      </div>
                      <Progress value={98.2} className="h-2" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">outlook.com</div>
                        <div className="text-sm text-muted-foreground">96.5%</div>
                      </div>
                      <Progress value={96.5} className="h-2" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">yahoo.com</div>
                        <div className="text-sm text-muted-foreground">94.8%</div>
                      </div>
                      <Progress value={94.8} className="h-2" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">hotmail.com</div>
                        <div className="text-sm text-muted-foreground">92.3%</div>
                      </div>
                      <Progress value={92.3} className="h-2" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Other domains</div>
                        <div className="text-sm text-muted-foreground">90.1%</div>
                      </div>
                      <Progress value={90.1} className="h-2" />
                    </div>
                  </div>
                </div>

                <div>
                  <h3 className="mb-4 text-lg font-medium">Failure Reasons</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Invalid email address</div>
                        <div className="text-sm text-muted-foreground">42%</div>
                      </div>
                      <Progress value={42} className="h-2 bg-muted [&>div]:bg-red-500" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Mailbox full</div>
                        <div className="text-sm text-muted-foreground">28%</div>
                      </div>
                      <Progress value={28} className="h-2 bg-muted [&>div]:bg-red-500" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Marked as spam</div>
                        <div className="text-sm text-muted-foreground">18%</div>
                      </div>
                      <Progress value={18} className="h-2 bg-muted [&>div]:bg-red-500" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Server error</div>
                        <div className="text-sm text-muted-foreground">8%</div>
                      </div>
                      <Progress value={8} className="h-2 bg-muted [&>div]:bg-red-500" />
                    </div>
                    <div>
                      <div className="mb-1 flex items-center justify-between">
                        <div className="text-sm">Other reasons</div>
                        <div className="text-sm text-muted-foreground">4%</div>
                      </div>
                      <Progress value={4} className="h-2 bg-muted [&>div]:bg-red-500" />
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
