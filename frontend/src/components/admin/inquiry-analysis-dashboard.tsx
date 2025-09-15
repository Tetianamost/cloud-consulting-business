import { useState } from "react"
import { Bot, TrendingUp, Clock, AlertTriangle, CheckCircle, Eye, Download } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/Card"
import { Badge } from "../ui/Badge"
import { Button } from "../ui/Button"
import { Progress } from "../ui/Progress"
import { BedrockReportGenerator } from "./bedrock-report-generator"

// Sample data for AI-generated reports
const aiReports = [
  {
    id: "INQ-RPT-2025-07-001",
    inquiryId: "INQ-1001",
    title: "Cloud Migration Quote Analysis - TechCorp Inc",
    customerName: "TechCorp Inc",
    customerEmail: "john.smith@techcorp.com",
    serviceType: "AWS Migration",
    priority: "High",
    status: "Ready",
    confidence: 0.92,
    bedrockModel: "Claude-3-Sonnet",
    generatedAt: "2025-07-26T10:30:00",
    processingTime: "8.5 minutes",
    keyInsights: [
      "Complex legacy system requiring phased migration approach",
      "Estimated cost savings of 35% after migration",
      "High compliance requirements for financial data",
    ],
    recommendedActions: [
      { action: "Schedule technical discovery call", priority: "Immediate", timeline: "Within 24 hours" },
      { action: "Prepare detailed cost breakdown", priority: "High", timeline: "Within 3 days" },
      { action: "Review compliance framework", priority: "High", timeline: "Within 1 week" },
    ],
    riskLevel: "Medium",
    estimatedValue: "$350K - $500K",
    timeline: "6-9 months",
  },
  {
    id: "INQ-RPT-2025-07-002",
    inquiryId: "INQ-1002",
    title: "Multi-Cloud Strategy Analysis - StartupXYZ",
    customerName: "StartupXYZ",
    customerEmail: "cto@startupxyz.com",
    serviceType: "Multi-Cloud Architecture",
    priority: "Medium",
    status: "Ready",
    confidence: 0.87,
    bedrockModel: "Claude-3-Haiku",
    generatedAt: "2025-07-25T14:45:00",
    processingTime: "4.2 minutes",
    keyInsights: [
      "Growing startup with scalability concerns",
      "Budget-conscious but quality-focused approach",
      "Need for vendor lock-in avoidance strategy",
    ],
    recommendedActions: [
      { action: "Send multi-cloud architecture whitepaper", priority: "High", timeline: "Within 2 days" },
      { action: "Schedule discovery session", priority: "Medium", timeline: "Within 1 week" },
    ],
    riskLevel: "Low",
    estimatedValue: "$75K - $150K",
    timeline: "3-4 months",
  },
  {
    id: "INQ-RPT-2025-07-003",
    inquiryId: "INQ-1003",
    title: "Data Analytics Platform Quote - RetailCorp",
    customerName: "RetailCorp",
    customerEmail: "data-team@retailcorp.com",
    serviceType: "Analytics Platform",
    priority: "High",
    status: "Ready",
    confidence: 0.94,
    bedrockModel: "Claude-3-Sonnet",
    generatedAt: "2025-07-24T09:15:00",
    processingTime: "12.1 minutes",
    keyInsights: [
      "Large retail chain with massive data volumes",
      "Real-time analytics requirements for inventory",
      "Integration with existing POS systems critical",
    ],
    recommendedActions: [
      { action: "Prepare demo environment", priority: "Immediate", timeline: "Within 48 hours" },
      { action: "Schedule technical deep-dive", priority: "High", timeline: "Within 3 days" },
      { action: "Review data governance requirements", priority: "High", timeline: "Within 1 week" },
    ],
    riskLevel: "High",
    estimatedValue: "$800K - $1.2M",
    timeline: "9-12 months",
  },
]

export function InquiryAnalysisDashboard() {
  const [selectedReport, setSelectedReport] = useState<any>(null)
  const [showGenerator, setShowGenerator] = useState(false)
  const [showPreview, setShowPreview] = useState(false)

  const handleViewReport = (report: any) => {
    setSelectedReport(report)
    setShowPreview(true)
  }

  const handleGenerateReport = (inquiry: any) => {
    setSelectedReport(inquiry)
    setShowGenerator(true)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "Critical":
        return "bg-red-100 text-red-800 border-red-200"
      case "High":
        return "bg-orange-100 text-orange-800 border-orange-200"
      case "Medium":
        return "bg-yellow-100 text-yellow-800 border-yellow-200"
      case "Low":
        return "bg-green-100 text-green-800 border-green-200"
      default:
        return "bg-gray-100 text-gray-800 border-gray-200"
    }
  }

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case "High":
        return "text-red-600"
      case "Medium":
        return "text-yellow-600"
      case "Low":
        return "text-green-600"
      default:
        return "text-gray-600"
    }
  }

  return (
    <div className="space-y-6">
      {/* Overview Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">AI Reports Generated</CardTitle>
            <Bot className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">24</div>
            <p className="text-xs text-muted-foreground">
              <span className="text-emerald-500">+8</span> this week
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Confidence Score</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">89.3%</div>
            <p className="text-xs text-muted-foreground">
              <span className="text-emerald-500">+2.1%</span> from last month
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Processing Time</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">7.2min</div>
            <p className="text-xs text-muted-foreground">
              <span className="text-emerald-500">-1.3min</span> improvement
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">High-Value Opportunities</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">7</div>
            <p className="text-xs text-muted-foreground">Requiring immediate attention</p>
          </CardContent>
        </Card>
      </div>

      {/* AI Reports List */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Bot className="h-5 w-5 text-blue-600" />
                AI-Generated Inquiry Analysis Reports
              </CardTitle>
              <CardDescription>Bedrock-powered analysis of customer inquiries with actionable insights</CardDescription>
            </div>
            <Button onClick={() => setShowGenerator(true)}>
              <Bot className="mr-2 h-4 w-4" />
              Generate New Report
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {aiReports.map((report) => (
              <div key={report.id} className="border rounded-lg p-4 space-y-3">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h3 className="font-semibold">{report.title}</h3>
                      <Badge variant="outline" className={getPriorityColor(report.priority)}>
                        {report.priority}
                      </Badge>
                      <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">
                        {report.bedrockModel}
                      </Badge>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-muted-foreground mb-3">
                      <div>
                        <span className="font-medium">Customer:</span> {report.customerName}
                      </div>
                      <div>
                        <span className="font-medium">Service:</span> {report.serviceType}
                      </div>
                      <div>
                        <span className="font-medium">Value:</span> {report.estimatedValue}
                      </div>
                      <div>
                        <span className="font-medium">Timeline:</span> {report.timeline}
                      </div>
                    </div>
                    <div className="flex items-center gap-4 mb-3">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Confidence:</span>
                        <Progress value={report.confidence * 100} className="w-20 h-2" />
                        <span className="text-sm text-muted-foreground">{(report.confidence * 100).toFixed(0)}%</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">Risk:</span>
                        <span className={`text-sm font-medium ${getRiskColor(report.riskLevel)}`}>
                          {report.riskLevel}
                        </span>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        Generated: {formatDate(report.generatedAt)} • {report.processingTime}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button variant="outline" size="sm" onClick={() => handleViewReport(report)}>
                      <Eye className="mr-2 h-4 w-4" />
                      View
                    </Button>
                    <Button variant="outline" size="sm">
                      <Download className="mr-2 h-4 w-4" />
                      Download
                    </Button>
                  </div>
                </div>

                {/* Key Insights */}
                <div>
                  <h4 className="text-sm font-medium mb-2">Key Insights:</h4>
                  <div className="space-y-1">
                    {report.keyInsights.map((insight, index) => (
                      <div key={index} className="flex items-start gap-2 text-sm">
                        <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                        <span>{insight}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Recommended Actions */}
                <div>
                  <h4 className="text-sm font-medium mb-2">Recommended Actions:</h4>
                  <div className="flex flex-wrap gap-2">
                    {report.recommendedActions.map((action, index) => (
                      <Badge
                        key={index}
                        variant="outline"
                        className={`${
                          action.priority === "Immediate"
                            ? "bg-red-50 text-red-700 border-red-200"
                            : action.priority === "High"
                              ? "bg-orange-50 text-orange-700 border-orange-200"
                              : "bg-blue-50 text-blue-700 border-blue-200"
                        }`}
                      >
                        {action.action} • {action.timeline}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Bedrock Report Generator Modal */}
      <BedrockReportGenerator
        isOpen={showGenerator}
        onClose={() => setShowGenerator(false)}
        onGenerate={(config) => {
          console.log("Generating report with config:", config)
          setShowGenerator(false)
        }}
        inquiry={selectedReport}
      />
    </div>
  )
}
