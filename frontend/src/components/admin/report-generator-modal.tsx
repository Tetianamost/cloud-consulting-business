import { useState } from "react"
import { Calendar, FileText, Settings, X } from "lucide-react"
import { Button } from "../ui/Button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/Card"
import { Checkbox } from "../ui/Checkbox"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/Dialog"
import { Input } from "../ui/Input"
import { Label } from "../ui/Label"
import { RadioGroup, RadioGroupItem } from "../ui/Radio-group"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/Select"
import { Textarea } from "../ui/Textarea"

interface ReportGeneratorModalProps {
  isOpen: boolean
  onClose: () => void
  onGenerate: (config: any) => void
}

export function ReportGeneratorModal({ isOpen, onClose, onGenerate }: ReportGeneratorModalProps) {
  const [reportConfig, setReportConfig] = useState({
    title: "",
    description: "",
    type: "performance",
    dateRange: "last_30_days",
    customStartDate: "",
    customEndDate: "",
    includeCharts: true,
    includeTables: true,
    includeRawData: false,
    formats: ["pdf"],
    sections: ["summary", "metrics", "trends"],
    recipients: "",
    schedule: "once",
  })

  const handleConfigChange = (key: string, value: any) => {
    setReportConfig((prev) => ({ ...prev, [key]: value }))
  }

  const handleSectionToggle = (section: string) => {
    setReportConfig((prev) => ({
      ...prev,
      sections: prev.sections.includes(section)
        ? prev.sections.filter((s) => s !== section)
        : [...prev.sections, section],
    }))
  }

  const handleFormatToggle = (format: string) => {
    setReportConfig((prev) => ({
      ...prev,
      formats: prev.formats.includes(format) ? prev.formats.filter((f) => f !== format) : [...prev.formats, format],
    }))
  }

  const handleGenerate = () => {
    if (!reportConfig.title.trim()) {
      alert("Please enter a report title")
      return
    }

    onGenerate(reportConfig)
  }

  const reportTypes = [
    { value: "performance", label: "Performance Report", description: "Overall system performance metrics" },
    { value: "inquiries", label: "Inquiry Analysis", description: "Detailed inquiry statistics and trends" },
    { value: "emails", label: "Email Analytics", description: "Email delivery and engagement metrics" },
    { value: "custom", label: "Custom Report", description: "Build a custom report with selected metrics" },
  ]

  const availableSections = [
    { id: "summary", label: "Executive Summary", description: "High-level overview and KPIs" },
    { id: "metrics", label: "Detailed Metrics", description: "Comprehensive performance data" },
    { id: "trends", label: "Trend Analysis", description: "Historical trends and comparisons" },
    { id: "breakdown", label: "Category Breakdown", description: "Data segmented by categories" },
    { id: "recommendations", label: "Recommendations", description: "AI-generated insights and suggestions" },
  ]

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div>
            <DialogTitle className="text-xl">Generate New Report</DialogTitle>
            <p className="text-sm text-muted-foreground mt-1">
              Configure and generate a custom report with your preferred settings
            </p>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </DialogHeader>

        <div className="space-y-6">
          {/* Basic Information */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <FileText className="h-5 w-5" />
                Basic Information
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="title">Report Title *</Label>
                <Input
                  id="title"
                  placeholder="Enter report title"
                  value={reportConfig.title}
                  onChange={(e) => handleConfigChange("title", e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  placeholder="Optional description for the report"
                  value={reportConfig.description}
                  onChange={(e) => handleConfigChange("description", e.target.value)}
                  rows={3}
                />
              </div>
              <div>
                <Label>Report Type</Label>
                <RadioGroup
                  value={reportConfig.type}
                  onValueChange={(value) => handleConfigChange("type", value)}
                  className="mt-2"
                >
                  {reportTypes.map((type) => (
                    <div key={type.value} className="flex items-start space-x-2 p-3 border rounded-lg">
                      <RadioGroupItem value={type.value} id={type.value} className="mt-1" />
                      <div className="flex-1">
                        <Label htmlFor={type.value} className="font-medium cursor-pointer">
                          {type.label}
                        </Label>
                        <p className="text-sm text-muted-foreground">{type.description}</p>
                      </div>
                    </div>
                  ))}
                </RadioGroup>
              </div>
            </CardContent>
          </Card>

          {/* Date Range */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Calendar className="h-5 w-5" />
                Date Range
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label>Time Period</Label>
                <Select
                  value={reportConfig.dateRange}
                  onValueChange={(value) => handleConfigChange("dateRange", value)}
                >
                  <SelectTrigger className="mt-2">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="last_7_days">Last 7 Days</SelectItem>
                    <SelectItem value="last_30_days">Last 30 Days</SelectItem>
                    <SelectItem value="last_90_days">Last 90 Days</SelectItem>
                    <SelectItem value="current_month">Current Month</SelectItem>
                    <SelectItem value="last_month">Last Month</SelectItem>
                    <SelectItem value="current_quarter">Current Quarter</SelectItem>
                    <SelectItem value="custom">Custom Range</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              {reportConfig.dateRange === "custom" && (
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="startDate">Start Date</Label>
                    <Input
                      id="startDate"
                      type="date"
                      value={reportConfig.customStartDate}
                      onChange={(e) => handleConfigChange("customStartDate", e.target.value)}
                    />
                  </div>
                  <div>
                    <Label htmlFor="endDate">End Date</Label>
                    <Input
                      id="endDate"
                      type="date"
                      value={reportConfig.customEndDate}
                      onChange={(e) => handleConfigChange("customEndDate", e.target.value)}
                    />
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Content Configuration */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Content Configuration
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-base font-medium">Include Sections</Label>
                <div className="mt-3 space-y-3">
                  {availableSections.map((section) => (
                    <div key={section.id} className="flex items-start space-x-3">
                      <Checkbox
                        id={section.id}
                        checked={reportConfig.sections.includes(section.id)}
                        onCheckedChange={() => handleSectionToggle(section.id)}
                      />
                      <div className="flex-1">
                        <Label htmlFor={section.id} className="font-medium cursor-pointer">
                          {section.label}
                        </Label>
                        <p className="text-sm text-muted-foreground">{section.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="includeCharts"
                    checked={reportConfig.includeCharts}
                    onCheckedChange={(checked) => handleConfigChange("includeCharts", checked)}
                  />
                  <Label htmlFor="includeCharts">Include Charts</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="includeTables"
                    checked={reportConfig.includeTables}
                    onCheckedChange={(checked) => handleConfigChange("includeTables", checked)}
                  />
                  <Label htmlFor="includeTables">Include Tables</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="includeRawData"
                    checked={reportConfig.includeRawData}
                    onCheckedChange={(checked) => handleConfigChange("includeRawData", checked)}
                  />
                  <Label htmlFor="includeRawData">Include Raw Data</Label>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Output Formats */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Output Formats</CardTitle>
              <CardDescription>Select the formats you want to generate</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-3 gap-4">
                {["pdf", "html", "csv"].map((format) => (
                  <div key={format} className="flex items-center space-x-2">
                    <Checkbox
                      id={format}
                      checked={reportConfig.formats.includes(format)}
                      onCheckedChange={() => handleFormatToggle(format)}
                    />
                    <Label htmlFor={format} className="cursor-pointer">
                      {format.toUpperCase()}
                    </Label>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Actions */}
          <div className="flex items-center justify-between pt-4 border-t">
            <div className="text-sm text-muted-foreground">Report will be generated and available for download</div>
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button onClick={handleGenerate}>Generate Report</Button>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
