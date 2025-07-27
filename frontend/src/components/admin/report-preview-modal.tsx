"use client"

import { useState } from "react"
import { Download, FileText, Globe, Maximize2, Minimize2, Table, X } from "lucide-react"

import { Badge } from "../ui/badge"
import { Button } from "../ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu"
import { ScrollArea } from "../ui/scroll-area"
import { Separator } from "../ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs"

interface ReportPreviewModalProps {
  report: any
  isOpen: boolean
  onClose: () => void
  onDownload: (report: any, format: 'pdf' | 'html') => Promise<void>
}

export function ReportPreviewModal({ report, isOpen, onClose, onDownload }: ReportPreviewModalProps) {
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [previewFormat, setPreviewFormat] = useState("summary")

  if (!report) return null

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  // Sample report data for preview
  const reportData = {
    summary: {
      customerName: report.title.split(" - ")[1] || "Unknown Customer",
      serviceType: report.serviceType,
      inquiryDate: report.generatedAt,
      priority: report.priority,
      confidence: report.confidence,
      bedrockModel: report.bedrockModel,
    },
    analysis: {
      keyRequirements: [
        "Migrate 50+ legacy applications to cloud",
        "Ensure 99.9% uptime SLA",
        "Implement disaster recovery",
        "Maintain regulatory compliance",
      ],
      technicalComplexity: "High",
      estimatedTimeline: "6-9 months",
      budgetRange: "$250K - $400K",
      riskFactors: ["Legacy system dependencies", "Data migration complexity", "Compliance requirements"],
    },
    recommendations: [
      {
        action: "Schedule technical discovery call",
        priority: "Immediate",
        reasoning: "Customer has complex legacy environment requiring detailed assessment",
      },
      {
        action: "Prepare detailed cost breakdown",
        priority: "High",
        reasoning: "Budget-conscious customer needs transparent pricing",
      },
      {
        action: "Review compliance framework",
        priority: "High",
        reasoning: "Regulatory requirements mentioned multiple times in inquiry",
      },
    ],
    nextSteps: [
      "Send technical questionnaire within 24 hours",
      "Schedule discovery call within 3 business days",
      "Prepare preliminary architecture proposal",
      "Identify compliance requirements and certifications needed",
    ],
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className={`max-w-4xl ${isFullscreen ? "h-[95vh] max-h-[95vh]" : "max-h-[80vh]"}`}>
        <DialogHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div>
            <DialogTitle className="text-xl">{report.title}</DialogTitle>
            <div className="flex items-center gap-2 mt-2">
              <Badge variant="outline">{report.type}</Badge>
              <span className="text-sm text-muted-foreground">Generated: {formatDate(report.generatedAt)}</span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" onClick={() => setIsFullscreen(!isFullscreen)}>
              {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm">
                  <Download className="mr-2 h-4 w-4" />
                  Download
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => onDownload(report, "pdf")}>
                  <FileText className="mr-2 h-4 w-4" />
                  Download PDF
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => onDownload(report, "html")}>
                  <Globe className="mr-2 h-4 w-4" />
                  Download HTML
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        </DialogHeader>

        <Tabs value={previewFormat} onValueChange={setPreviewFormat} className="flex-1">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="summary">Executive Summary</TabsTrigger>
            <TabsTrigger value="detailed">Technical Analysis</TabsTrigger>
            <TabsTrigger value="recommendations">Recommendations</TabsTrigger>
            <TabsTrigger value="nextSteps">Next Steps</TabsTrigger>
          </TabsList>

          <ScrollArea className={`${isFullscreen ? "h-[calc(95vh-200px)]" : "h-[500px]"} mt-4`}>
            <TabsContent value="summary" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Executive Summary</CardTitle>
                  <CardDescription>Customer information and AI confidence</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <h4 className="font-medium">Customer Name</h4>
                      <p className="text-sm text-muted-foreground">{reportData.summary.customerName}</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Service Type</h4>
                      <p className="text-sm text-muted-foreground">{reportData.summary.serviceType}</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Inquiry Date</h4>
                      <p className="text-sm text-muted-foreground">{formatDate(reportData.summary.inquiryDate)}</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Priority</h4>
                      <p className="text-sm text-muted-foreground">{reportData.summary.priority}</p>
                    </div>
                    <div>
                      <h4 className="font-medium">AI Confidence</h4>
                      <p className="text-sm text-muted-foreground">{reportData.summary.confidence}</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Bedrock Model</h4>
                      <p className="text-sm text-muted-foreground">{reportData.summary.bedrockModel}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="detailed" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Technical Analysis</CardTitle>
                  <CardDescription>Key requirements and complexity</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <h4 className="font-medium">Key Requirements</h4>
                    <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                      {reportData.analysis.keyRequirements.map((req, index) => (
                        <li key={index}>{req}</li>
                      ))}
                    </ul>
                    <Separator />
                    <h4 className="font-medium">Technical Complexity</h4>
                    <p className="text-sm text-muted-foreground">{reportData.analysis.technicalComplexity}</p>
                    <Separator />
                    <h4 className="font-medium">Estimated Timeline</h4>
                    <p className="text-sm text-muted-foreground">{reportData.analysis.estimatedTimeline}</p>
                    <Separator />
                    <h4 className="font-medium">Budget Range</h4>
                    <p className="text-sm text-muted-foreground">{reportData.analysis.budgetRange}</p>
                    <Separator />
                    <h4 className="font-medium">Risk Factors</h4>
                    <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                      {reportData.analysis.riskFactors.map((risk, index) => (
                        <li key={index}>{risk}</li>
                      ))}
                    </ul>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="recommendations" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Recommendations</CardTitle>
                  <CardDescription>Prioritized action items</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {reportData.recommendations.map((rec, index) => (
                      <div key={index} className="p-4 border rounded-lg">
                        <h4 className="font-medium">{rec.action}</h4>
                        <p className="text-sm text-muted-foreground">
                          Priority: {rec.priority}
                          <br />
                          Reasoning: {rec.reasoning}
                        </p>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="nextSteps" className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Next Steps</CardTitle>
                  <CardDescription>Timeline and deliverables</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                      {reportData.nextSteps.map((step, index) => (
                        <li key={index}>{step}</li>
                      ))}
                    </ul>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>
          </ScrollArea>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
