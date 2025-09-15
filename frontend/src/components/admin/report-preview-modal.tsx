"use client"

import { useState, useEffect } from "react"
import { Download, FileText, Globe, Maximize2, Minimize2, Table, X } from "lucide-react"
import { Badge } from "../ui/Badge"
import { Button } from "../ui/Button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/Card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/Dialog"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/Dropdown-menu"
import { ScrollArea } from "../ui/Scroll-area"
import { Separator } from "../ui/Separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/Tabs"
import apiService from "../../services/api"

// Markdown rendering
import { marked } from "marked"

interface ReportPreviewModalProps {
  report?: any
  reportId?: string
  isOpen: boolean
  onClose: () => void
  onDownload: (report: any, format: 'pdf' | 'html') => Promise<void>
}

export function ReportPreviewModal({ report: initialReport, reportId, isOpen, onClose, onDownload }: ReportPreviewModalProps) {
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [previewFormat, setPreviewFormat] = useState("summary")
  const [loading, setLoading] = useState(false)
  const [report, setReport] = useState<any>(initialReport || null)
  const [error, setError] = useState<string | null>(null)
  const [downloadLoading, setDownloadLoading] = useState<string | null>(null)

  // Fetch report from backend if reportId is provided
  useEffect(() => {
    if (isOpen && reportId) {
      setLoading(true)
      setError(null)
      apiService.getReport(reportId)
        .then((data: any) => {
          setReport(data)
        })
        .catch(() => setError("Failed to load report"))
        .finally(() => setLoading(false))
    } else if (isOpen && initialReport) {
      setReport(initialReport)
    }
  }, [isOpen, reportId])

  if (!isOpen) return null

  const formatDate = (dateString: string) => {
    if (!dateString) return ""
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  // Parse and structure the report content
  const inquiry = report?.inquiry || {};
  const reportContent = report?.content || "";
  
  // Parse the content into sections based on common patterns
  const parseReportContent = (content: string) => {
    if (!content) return { summary: "", analysis: "", recommendations: "", nextSteps: "" };
    
    // Split content into sections based on common headers
    const sections = {
      summary: "",
      analysis: "",
      recommendations: "",
      nextSteps: ""
    };
    
    // Try to extract sections based on common patterns
    const summaryMatch = content.match(/(?:EXECUTIVE SUMMARY|Summary|Overview)([\s\S]*?)(?=\n(?:TECHNICAL ANALYSIS|Analysis|RECOMMENDATIONS|Recommendations|NEXT STEPS|Next Steps)|$)/i);
    const analysisMatch = content.match(/(?:TECHNICAL ANALYSIS|Analysis|Technical Details)([\s\S]*?)(?=\n(?:RECOMMENDATIONS|Recommendations|NEXT STEPS|Next Steps)|$)/i);
    const recommendationsMatch = content.match(/(?:RECOMMENDATIONS|Recommendations|Action Items)([\s\S]*?)(?=\n(?:NEXT STEPS|Next Steps)|$)/i);
    const nextStepsMatch = content.match(/(?:NEXT STEPS|Next Steps|Timeline)([\s\S]*?)$/i);
    
    sections.summary = summaryMatch ? summaryMatch[1].trim() : content.substring(0, 500) + "...";
    sections.analysis = analysisMatch ? analysisMatch[1].trim() : "";
    sections.recommendations = recommendationsMatch ? recommendationsMatch[1].trim() : "";
    sections.nextSteps = nextStepsMatch ? nextStepsMatch[1].trim() : "";
    
    return sections;
  };
  
  const parsedContent = parseReportContent(reportContent);
  
  const reportData = {
    summary: {
      customerName: inquiry.name || "Unknown Customer",
      serviceType: report?.type || (inquiry.services ? inquiry.services.join(", ") : ""),
      inquiryDate: report?.created_at || inquiry.created_at || "",
      priority: inquiry.priority || "Medium",
      generatedBy: report?.generated_by || "AI Assistant",
      status: report?.status || "Draft",
    },
    content: {
      summary: parsedContent.summary,
      analysis: parsedContent.analysis,
      recommendations: parsedContent.recommendations,
      nextSteps: parsedContent.nextSteps,
      fullContent: reportContent,
    }
  }

  // Markdown renderer using marked
  function markdownToHtml(md: string): string {
    if (!md) return ""
    try {
      const result = marked.parse(md)
      return typeof result === 'string' ? result : md.replace(/\n/g, "<br/>")
    } catch {
      return md.replace(/\n/g, "<br/>")
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent
        className={`w-full max-w-3xl mx-auto ${isFullscreen ? "h-[95vh] max-h-[95vh]" : "max-h-[80vh]"} p-0 sm:px-8 sm:py-6`}
        style={{ borderRadius: 12 }}
      >
        <>
          <DialogHeader className="pb-0 px-6 pt-6 border-b border-gray-200">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
              <div>
                <DialogTitle className="text-xl font-semibold text-gray-900 truncate mb-1">
                  {report?.title || "Report Preview"}
                </DialogTitle>
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant="outline" className="capitalize text-xs px-2 py-1">
                    {report?.type || 'Report'}
                  </Badge>
                  <Badge variant={report?.status === 'completed' ? 'default' : 'secondary'} className="capitalize text-xs px-2 py-1">
                    {report?.status || 'Draft'}
                  </Badge>
                  <span className="text-xs text-gray-500">
                    {formatDate(report?.created_at || report?.generatedAt)}
                  </span>
                  <span className="text-xs text-gray-500">
                    Generated by: {reportData.summary.generatedBy}
                  </span>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsFullscreen(!isFullscreen)}
                  title={isFullscreen ? "Exit fullscreen" : "Enter fullscreen"}
                >
                  {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
                </Button>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline" size="sm" className="bg-blue-50 hover:bg-blue-100 border-blue-200">
                      <Download className="mr-2 h-4 w-4" />
                      Download
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-48">
                    <DropdownMenuItem
                      onClick={async () => {
                        setDownloadLoading('pdf');
                        try {
                          await onDownload(report, "pdf");
                        } catch (error) {
                          console.error('PDF download failed:', error);
                        } finally {
                          setDownloadLoading(null);
                        }
                      }}
                      className="cursor-pointer"
                      disabled={downloadLoading === 'pdf'}
                    >
                      {downloadLoading === 'pdf' ? (
                        <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-red-600 border-t-transparent" />
                      ) : (
                        <FileText className="mr-2 h-4 w-4 text-red-600" />
                      )}
                      <div>
                        <div className="font-medium">
                          {downloadLoading === 'pdf' ? 'Generating PDF...' : 'Download PDF'}
                        </div>
                        <div className="text-xs text-gray-500">Formatted for printing</div>
                      </div>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={async () => {
                        setDownloadLoading('html');
                        try {
                          await onDownload(report, "html");
                        } catch (error) {
                          console.error('HTML download failed:', error);
                        } finally {
                          setDownloadLoading(null);
                        }
                      }}
                      className="cursor-pointer"
                      disabled={downloadLoading === 'html'}
                    >
                      {downloadLoading === 'html' ? (
                        <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-blue-600 border-t-transparent" />
                      ) : (
                        <Globe className="mr-2 h-4 w-4 text-blue-600" />
                      )}
                      <div>
                        <div className="font-medium">
                          {downloadLoading === 'html' ? 'Generating HTML...' : 'Download HTML'}
                        </div>
                        <div className="text-xs text-gray-500">Web-friendly format</div>
                      </div>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
                <Button variant="ghost" size="sm" onClick={onClose} title="Close">
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </DialogHeader>

          {loading ? (
            <div className="flex flex-col items-center justify-center h-64 space-y-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <div className="text-center">
                <p className="text-gray-600 font-medium">Loading report...</p>
                <p className="text-sm text-gray-500">Please wait while we fetch the report data</p>
              </div>
            </div>
          ) : error ? (
            <div className="flex flex-col items-center justify-center h-64 space-y-4">
              <div className="p-3 bg-red-100 rounded-full">
                <X className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-center">
                <p className="text-red-600 font-medium">Failed to load report</p>
                <p className="text-sm text-gray-500">{error}</p>
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={() => window.location.reload()} 
                  className="mt-2"
                >
                  Try Again
                </Button>
              </div>
            </div>
          ) : (
            <Tabs value={previewFormat} onValueChange={setPreviewFormat} className="flex-1 px-2 pb-4">
              <TabsList className="grid w-full grid-cols-2 gap-2 mb-4">
                <TabsTrigger value="summary" className="w-full">Report Overview</TabsTrigger>
                <TabsTrigger value="fullContent" className="w-full">Full Report</TabsTrigger>
              </TabsList>

              <TabsContent value="summary" className="space-y-4 px-0 py-2">
                <ScrollArea className={`${isFullscreen ? "h-[calc(95vh-220px)]" : "h-[420px]"} mt-2`}>
                  <div className="space-y-6">
                    {/* Key Information Cards */}
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                      <Card className="bg-blue-50 border-blue-200">
                        <CardContent className="p-4">
                          <div className="flex items-center space-x-2">
                            <div className="p-2 bg-blue-100 rounded-lg">
                              <FileText className="h-4 w-4 text-blue-600" />
                            </div>
                            <div>
                              <p className="text-sm font-medium text-blue-900">Customer</p>
                              <p className="text-lg font-bold text-blue-800">{reportData.summary.customerName}</p>
                              {reportData.summary.customerName !== inquiry.email && (
                                <p className="text-xs text-blue-700">{inquiry.email}</p>
                              )}
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                      
                      <Card className="bg-green-50 border-green-200">
                        <CardContent className="p-4">
                          <div className="flex items-center space-x-2">
                            <div className="p-2 bg-green-100 rounded-lg">
                              <Globe className="h-4 w-4 text-green-600" />
                            </div>
                            <div>
                              <p className="text-sm font-medium text-green-900">Service</p>
                              <p className="text-lg font-bold text-green-800">{reportData.summary.serviceType}</p>
                              {inquiry.company && (
                                <p className="text-xs text-green-700">{inquiry.company}</p>
                              )}
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                      
                      <Card className="bg-purple-50 border-purple-200">
                        <CardContent className="p-4">
                          <div className="flex items-center space-x-2">
                            <div className="p-2 bg-purple-100 rounded-lg">
                              <Table className="h-4 w-4 text-purple-600" />
                            </div>
                            <div>
                              <p className="text-sm font-medium text-purple-900">Status</p>
                              <p className="text-lg font-bold text-purple-800 capitalize">{reportData.summary.status}</p>
                              <p className="text-xs text-purple-700">{formatDate(reportData.summary.inquiryDate)}</p>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    </div>
                    
                    {/* Report Details */}
                    <Card>
                      <CardHeader className="pb-3">
                        <CardTitle className="text-lg flex items-center">
                          <FileText className="h-5 w-5 mr-2 text-blue-600" />
                          Report Details
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <dl className="grid grid-cols-2 gap-x-4 gap-y-2 text-xs">
                          <dt className="font-medium text-gray-600 col-span-1">Generated:</dt>
                          <dd className="text-gray-900 col-span-1">{formatDate(reportData.summary.inquiryDate)}</dd>
                          <dt className="font-medium text-gray-600 col-span-1">Generated By:</dt>
                          <dd className="text-gray-900 col-span-1">{reportData.summary.generatedBy}</dd>
                          <dt className="font-medium text-gray-600 col-span-1">Priority:</dt>
                          <dd className="col-span-1">
                            <Badge variant={reportData.summary.priority === 'High' ? 'destructive' : 'secondary'}>
                              {reportData.summary.priority}
                            </Badge>
                          </dd>
                          <dt className="font-medium text-gray-600 col-span-1">Report ID:</dt>
                          <dd className="font-mono text-xs text-gray-700 col-span-1">{report?.id}</dd>
                        </dl>
                      </CardContent>
                    </Card>
                    
                    {/* Executive Summary */}
                    {reportData.content.summary && (
                      <Card>
                        <CardHeader className="pb-3">
                          <CardTitle className="text-lg flex items-center">
                            <FileText className="h-5 w-5 mr-2 text-green-600" />
                            Executive Summary
                          </CardTitle>
                          <CardDescription>Key insights and overview</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="prose prose-xs max-w-none text-gray-700 leading-snug">
                            <div dangerouslySetInnerHTML={{ __html: markdownToHtml(reportData.content.summary) }} />
                          </div>
                        </CardContent>
                      </Card>
                    )}
                  </div>
                </ScrollArea>
              </TabsContent>

              <TabsContent value="fullContent" className="space-y-4 px-2 py-2">
                <ScrollArea className={`${isFullscreen ? "h-[calc(95vh-220px)]" : "h-[420px]"} mt-2`}>
                  <Card>
                    <CardHeader className="pb-4">
                      <CardTitle className="text-lg flex items-center">
                        <FileText className="h-5 w-5 mr-2 text-blue-600" />
                        Complete Report Analysis
                      </CardTitle>
                      <CardDescription>
                        Full AI-generated analysis with detailed insights and recommendations
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      {reportData.content.fullContent ? (
                        <div className="prose prose-xs max-w-none leading-snug text-gray-700">
                          <div
                            dangerouslySetInnerHTML={{ __html: markdownToHtml(reportData.content.fullContent) }}
                          />
                        </div>
                      ) : (
                        <div className="text-center py-8">
                          <FileText className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                          <p className="text-gray-600 font-medium">No content available</p>
                          <p className="text-sm text-gray-500">This report may still be processing or may have encountered an error.</p>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                </ScrollArea>
              </TabsContent>
            </Tabs>
          )}
        </>
      </DialogContent>
    </Dialog>
  );
}
