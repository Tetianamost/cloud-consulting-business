import { useState } from "react"
import { Bot, FileText, Loader2, Settings, Sparkles, X } from "lucide-react"
import { Button } from "../ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog"
import { Label } from "../ui/label"
import { RadioGroup, RadioGroupItem } from "../ui/radio-group"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select"
import { Textarea } from "../ui/textarea"
import { Badge } from "../ui/badge"
import { Progress } from "../ui/progress"

interface BedrockReportGeneratorProps {
  isOpen: boolean
  onClose: () => void
  onGenerate: (config: any) => void
  inquiry?: any
}

export function BedrockReportGenerator({ isOpen, onClose, onGenerate, inquiry }: BedrockReportGeneratorProps) {
  const [reportConfig, setReportConfig] = useState({
    inquiryId: inquiry?.id || "",
    customerEmail: inquiry?.email || "",
    serviceType: "",
    analysisType: "comprehensive",
    bedrockModel: "claude-3-sonnet",
    focusAreas: ["technical_requirements", "cost_analysis", "timeline_estimation"],
    customPrompt: "",
    includeRecommendations: true,
    includeNextSteps: true,
    includeRiskAssessment: true,
    confidenceThreshold: 0.8,
  })

  const [isGenerating, setIsGenerating] = useState(false)
  const [generationProgress, setGenerationProgress] = useState(0)

  const handleConfigChange = (key: string, value: any) => {
    setReportConfig((prev) => ({ ...prev, [key]: value }))
  }

  const handleFocusAreaToggle = (area: string) => {
    setReportConfig((prev) => ({
      ...prev,
      focusAreas: prev.focusAreas.includes(area)
        ? prev.focusAreas.filter((a) => a !== area)
        : [...prev.focusAreas, area],
    }))
  }

  const handleGenerate = async () => {
    setIsGenerating(true)
    setGenerationProgress(0)

    // Simulate Bedrock API call with progress updates
    const progressSteps = [
      { step: 20, message: "Analyzing customer inquiry..." },
      { step: 40, message: "Extracting technical requirements..." },
      { step: 60, message: "Generating cost estimates..." },
      { step: 80, message: "Creating recommendations..." },
      { step: 100, message: "Finalizing report..." },
    ]

    for (const { step, message } of progressSteps) {
      await new Promise((resolve) => setTimeout(resolve, 1000))
      setGenerationProgress(step)
    }

    onGenerate(reportConfig)
    setIsGenerating(false)
    setGenerationProgress(0)
  }

  const bedrockModels = [
    {
      value: "claude-3-sonnet",
      label: "Claude 3 Sonnet",
      description: "Best for complex analysis and detailed recommendations",
      cost: "$$",
    },
    {
      value: "claude-3-haiku",
      label: "Claude 3 Haiku",
      description: "Fast and cost-effective for standard analysis",
      cost: "$",
    },
    {
      value: "claude-3-opus",
      label: "Claude 3 Opus",
      description: "Most capable model for critical business decisions",
      cost: "$$$",
    },
  ]

  const analysisTypes = [
    { value: "quick", label: "Quick Analysis", description: "Basic requirements and recommendations (2-3 min)" },
    {
      value: "standard",
      label: "Standard Analysis",
      description: "Comprehensive review with cost estimates (5-7 min)",
    },
    {
      value: "comprehensive",
      label: "Comprehensive Analysis",
      description: "Deep dive with risk assessment (10-15 min)",
    },
  ]

  const focusAreas = [
    {
      id: "technical_requirements",
      label: "Technical Requirements",
      description: "Infrastructure and architecture needs",
    },
    { id: "cost_analysis", label: "Cost Analysis", description: "Pricing estimates and budget planning" },
    { id: "timeline_estimation", label: "Timeline Estimation", description: "Project phases and delivery schedule" },
    { id: "risk_assessment", label: "Risk Assessment", description: "Potential challenges and mitigation strategies" },
    { id: "compliance_review", label: "Compliance Review", description: "Security and regulatory requirements" },
    { id: "migration_strategy", label: "Migration Strategy", description: "Step-by-step migration approach" },
  ]

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div className="flex items-center gap-2">
            <Bot className="h-6 w-6 text-blue-600" />
            <div>
              <DialogTitle className="text-xl">Generate AI Analysis Report</DialogTitle>
              <p className="text-sm text-muted-foreground mt-1">
                Use Amazon Bedrock to analyze customer inquiry and generate actionable insights
              </p>
            </div>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </DialogHeader>

        {isGenerating ? (
          <div className="space-y-6 py-8">
            <div className="text-center">
              <Loader2 className="h-12 w-12 animate-spin mx-auto text-blue-600 mb-4" />
              <h3 className="text-lg font-medium mb-2">Generating AI Analysis...</h3>
              <p className="text-sm text-muted-foreground mb-4">
                Bedrock is analyzing the customer inquiry and generating insights
              </p>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span>Progress</span>
                <span>{generationProgress}%</span>
              </div>
              <Progress value={generationProgress} className="h-2" />
            </div>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Inquiry Information */}
            {inquiry && (
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    Inquiry Details
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <Label className="text-sm font-medium">Inquiry ID</Label>
                      <p className="text-sm text-muted-foreground">{inquiry.id}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium">Customer</Label>
                      <p className="text-sm text-muted-foreground">{inquiry.customer}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium">Subject</Label>
                      <p className="text-sm text-muted-foreground">{inquiry.subject}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium">Priority</Label>
                      <Badge
                        variant="outline"
                        className={inquiry.priority === "High" ? "border-red-200 text-red-800" : ""}
                      >
                        {inquiry.priority}
                      </Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Bedrock Configuration */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Bot className="h-5 w-5" />
                  AI Model Configuration
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label>Bedrock Model</Label>
                  <RadioGroup
                    value={reportConfig.bedrockModel}
                    onValueChange={(value) => handleConfigChange("bedrockModel", value)}
                    className="mt-2"
                  >
                    {bedrockModels.map((model) => (
                      <div key={model.value} className="flex items-start space-x-2 p-3 border rounded-lg">
                        <RadioGroupItem value={model.value} id={model.value} className="mt-1" />
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <Label htmlFor={model.value} className="font-medium cursor-pointer">
                              {model.label}
                            </Label>
                            <Badge variant="outline" className="text-xs">
                              {model.cost}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground">{model.description}</p>
                        </div>
                      </div>
                    ))}
                  </RadioGroup>
                </div>

                <div>
                  <Label>Analysis Type</Label>
                  <Select
                    value={reportConfig.analysisType}
                    onValueChange={(value) => handleConfigChange("analysisType", value)}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {analysisTypes.map((type) => (
                        <SelectItem key={type.value} value={type.value}>
                          <div>
                            <div className="font-medium">{type.label}</div>
                            <div className="text-xs text-muted-foreground">{type.description}</div>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label>Service Type</Label>
                  <Select
                    value={reportConfig.serviceType}
                    onValueChange={(value) => handleConfigChange("serviceType", value)}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Select service type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="aws_migration">AWS Migration</SelectItem>
                      <SelectItem value="multi_cloud">Multi-Cloud Strategy</SelectItem>
                      <SelectItem value="analytics_platform">Analytics Platform</SelectItem>
                      <SelectItem value="security_compliance">Security & Compliance</SelectItem>
                      <SelectItem value="devops_automation">DevOps Automation</SelectItem>
                      <SelectItem value="data_modernization">Data Modernization</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>

            {/* Analysis Focus Areas */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Settings className="h-5 w-5" />
                  Analysis Focus Areas
                </CardTitle>
                <CardDescription>Select the areas you want Bedrock to focus on in the analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {focusAreas.map((area) => (
                    <div
                      key={area.id}
                      className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                        reportConfig.focusAreas.includes(area.id)
                          ? "border-blue-200 bg-blue-50"
                          : "border-gray-200 hover:border-gray-300"
                      }`}
                      onClick={() => handleFocusAreaToggle(area.id)}
                    >
                      <div className="flex items-center gap-2">
                        <input
                          type="checkbox"
                          checked={reportConfig.focusAreas.includes(area.id)}
                          onChange={() => handleFocusAreaToggle(area.id)}
                          className="rounded"
                        />
                        <div>
                          <div className="font-medium text-sm">{area.label}</div>
                          <div className="text-xs text-muted-foreground">{area.description}</div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Custom Instructions */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Sparkles className="h-5 w-5" />
                  Custom Instructions (Optional)
                </CardTitle>
                <CardDescription>
                  Provide additional context or specific requirements for the AI analysis
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Textarea
                  placeholder="e.g., Focus on cost optimization, consider regulatory requirements for healthcare industry, emphasize security best practices..."
                  value={reportConfig.customPrompt}
                  onChange={(e) => handleConfigChange("customPrompt", e.target.value)}
                  rows={4}
                />
              </CardContent>
            </Card>

            {/* Advanced Settings */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Advanced Settings</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label>Confidence Threshold</Label>
                  <Select
                    value={reportConfig.confidenceThreshold.toString()}
                    onValueChange={(value) => handleConfigChange("confidenceThreshold", Number.parseFloat(value))}
                  >
                    <SelectTrigger className="mt-2">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="0.6">60% - Include all insights</SelectItem>
                      <SelectItem value="0.7">70% - Standard confidence</SelectItem>
                      <SelectItem value="0.8">80% - High confidence (Recommended)</SelectItem>
                      <SelectItem value="0.9">90% - Very high confidence only</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="grid grid-cols-3 gap-4">
                  <div className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      id="includeRecommendations"
                      checked={reportConfig.includeRecommendations}
                      onChange={(e) => handleConfigChange("includeRecommendations", e.target.checked)}
                      className="rounded"
                    />
                    <Label htmlFor="includeRecommendations" className="text-sm">
                      Include Recommendations
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      id="includeNextSteps"
                      checked={reportConfig.includeNextSteps}
                      onChange={(e) => handleConfigChange("includeNextSteps", e.target.checked)}
                      className="rounded"
                    />
                    <Label htmlFor="includeNextSteps" className="text-sm">
                      Include Next Steps
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      id="includeRiskAssessment"
                      checked={reportConfig.includeRiskAssessment}
                      onChange={(e) => handleConfigChange("includeRiskAssessment", e.target.checked)}
                      className="rounded"
                    />
                    <Label htmlFor="includeRiskAssessment" className="text-sm">
                      Include Risk Assessment
                    </Label>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Actions */}
            <div className="flex items-center justify-between pt-4 border-t">
              <div className="text-sm text-muted-foreground">
                Estimated generation time:{" "}
                {reportConfig.analysisType === "quick"
                  ? "2-3"
                  : reportConfig.analysisType === "standard"
                    ? "5-7"
                    : "10-15"}{" "}
                minutes
              </div>
              <div className="flex items-center gap-2">
                <Button variant="outline" onClick={onClose}>
                  Cancel
                </Button>
                <Button onClick={handleGenerate} disabled={!reportConfig.serviceType}>
                  <Bot className="mr-2 h-4 w-4" />
                  Generate AI Report
                </Button>
              </div>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
