package services

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// pdfService implements the PDFService interface using gofpdf
type pdfService struct {
	logger      *logrus.Logger
	initialized bool
}

// NewPDFService creates a new PDF service instance
func NewPDFService(logger *logrus.Logger) interfaces.PDFService {
	service := &pdfService{
		logger:      logger,
		initialized: true,
	}
	
	logger.Info("PDF service initialized successfully with gofpdf")
	
	return service
}

// GeneratePDF generates a PDF from HTML content using gofpdf
func (p *pdfService) GeneratePDF(ctx context.Context, htmlContent string, options *interfaces.PDFOptions) ([]byte, error) {
	if !p.initialized {
		return nil, fmt.Errorf("PDF service not initialized")
	}
	
	// Set default options if not provided
	if options == nil {
		options = p.getDefaultOptions()
	}
	
	return p.generatePDFWithGofpdf(ctx, htmlContent, options)
}

// generatePDFWithGofpdf generates PDF using gofpdf as fallback (plain text)
func (p *pdfService) generatePDFWithGofpdf(ctx context.Context, htmlContent string, options *interfaces.PDFOptions) ([]byte, error) {
	// Convert HTML to plain text for PDF generation
	plainText := p.htmlToPlainText(htmlContent)
	
	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	
	// Set margins
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	
	// Add a page
	pdf.AddPage()
	
	// Set font
	pdf.SetFont("Arial", "", 12)
	
	// Add header if provided
	if options.HeaderHTML != "" {
		headerText := p.htmlToPlainText(options.HeaderHTML)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 10, headerText)
		pdf.Ln(15)
		pdf.SetFont("Arial", "", 12)
	}
	
	// Process the content
	p.addContentToPDF(pdf, plainText)
	
	// Add footer if provided
	if options.FooterHTML != "" {
		footerText := p.htmlToPlainText(options.FooterHTML)
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(0, 10, footerText)
	}
	
	// Check for errors first
	if pdf.Error() != nil {
		p.logger.WithError(pdf.Error()).Error("Failed to generate PDF")
		return nil, fmt.Errorf("failed to generate PDF: %w", pdf.Error())
	}
	
	// Get PDF bytes using a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		p.logger.WithError(err).Error("Failed to output PDF to buffer")
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}
	
	pdfBytes := buf.Bytes()
	
	p.logger.WithFields(logrus.Fields{
		"pdf_size": len(pdfBytes),
		"method":   "gofpdf_fallback",
	}).Info("PDF generated successfully")
	
	return pdfBytes, nil
}

// GeneratePDFFromURL generates a PDF from a URL (not supported with gofpdf)
func (p *pdfService) GeneratePDFFromURL(ctx context.Context, url string, options *interfaces.PDFOptions) ([]byte, error) {
	return nil, fmt.Errorf("PDF generation from URL not supported with gofpdf - use GeneratePDF with HTML content instead")
}

// IsHealthy checks if the PDF service is healthy and ready to generate PDFs
func (p *pdfService) IsHealthy() bool {
	if !p.initialized {
		return false
	}
	
	// Test with a simple HTML content
	testHTML := "<html><body><h1>Test</h1></body></html>"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := p.GeneratePDF(ctx, testHTML, nil)
	return err == nil
}

// GetVersion returns the version of the gofpdf library
func (p *pdfService) GetVersion() string {
	return "gofpdf v1.16.2"
}

// getDefaultOptions returns default PDF generation options
func (p *pdfService) getDefaultOptions() *interfaces.PDFOptions {
	return &interfaces.PDFOptions{
		PageSize:         "A4",
		Orientation:      "Portrait",
		MarginTop:        "1in",
		MarginRight:      "0.75in",
		MarginBottom:     "1in",
		MarginLeft:       "0.75in",
		EnableJavaScript: false,
		LoadTimeout:      30,
		Quality:          94,
		CustomOptions:    make(map[string]string),
	}
}

// htmlToPlainText converts HTML content to plain text for PDF generation
func (p *pdfService) htmlToPlainText(htmlContent string) string {
	// Remove HTML tags
	text := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(htmlContent, "")
	
	// Convert HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	
	// Clean up extra whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	
	return text
}

// addContentToPDF adds formatted content to the PDF
func (p *pdfService) addContentToPDF(pdf *gofpdf.Fpdf, content string) {
	// Split content into lines
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			pdf.Ln(5) // Add some space for empty lines
			continue
		}
		
		// Check if this looks like a header (all caps or starts with numbers)
		if p.isHeader(line) {
			pdf.Ln(5)
			pdf.SetFont("Arial", "B", 14)
			pdf.Cell(0, 10, line)
			pdf.Ln(10)
			pdf.SetFont("Arial", "", 12)
		} else {
			// Regular text - handle word wrapping
			p.addWrappedText(pdf, line)
		}
	}
}

// isHeader determines if a line should be treated as a header
func (p *pdfService) isHeader(line string) bool {
	// Check for numbered headers (1., 2., etc.)
	if matched, _ := regexp.MatchString(`^\d+\.`, line); matched {
		return true
	}
	
	// Check for all caps headers (but not too long)
	if strings.ToUpper(line) == line && len(line) < 100 && !strings.Contains(line, "\n") {
		return true
	}
	
	// Check for headers with specific keywords
	headerKeywords := []string{
		"EXECUTIVE SUMMARY", "CURRENT STATE", "RECOMMENDATIONS", "NEXT STEPS",
		"ASSESSMENT", "MIGRATION", "OPTIMIZATION", "ARCHITECTURE",
		"PRIORITY LEVEL", "URGENCY ASSESSMENT", "CONTACT INFORMATION",
	}
	
	lineUpper := strings.ToUpper(line)
	for _, keyword := range headerKeywords {
		if lineUpper == keyword || strings.HasPrefix(lineUpper, keyword+":") {
			return true
		}
	}
	
	return false
}

// addWrappedText adds text with word wrapping to the PDF
func (p *pdfService) addWrappedText(pdf *gofpdf.Fpdf, text string) {
	// Simple word wrapping - split long lines
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}
	
	const maxLineLength = 80 // Approximate characters per line
	currentLine := ""
	
	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxLineLength {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			if currentLine != "" {
				pdf.Cell(0, 6, currentLine)
				pdf.Ln(6)
			}
			currentLine = word
		}
	}
	
	if currentLine != "" {
		pdf.Cell(0, 6, currentLine)
		pdf.Ln(6)
	}
}

// enhanceHTMLForPDF enhances HTML content for better PDF rendering
func (p *pdfService) enhanceHTMLForPDF(htmlContent string, options *interfaces.PDFOptions) string {
	// If the HTML already has a complete structure, use it as-is
	if strings.Contains(strings.ToLower(htmlContent), "<!doctype") || 
	   strings.Contains(strings.ToLower(htmlContent), "<html") {
		return p.addPrintStyles(htmlContent)
	}
	
	// Otherwise, wrap in a complete HTML structure
	enhancedHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cloud Consulting Report</title>
    <style>
        /* Reset and base styles */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            font-size: 11pt;
            line-height: 1.6;
            color: #333;
            background: white;
            margin: 0;
            padding: 20px;
        }
        
        /* Typography */
        h1 {
            font-size: 24pt;
            font-weight: 700;
            color: #007cba;
            margin: 20px 0 15px 0;
            border-bottom: 2px solid #007cba;
            padding-bottom: 10px;
            page-break-after: avoid;
        }
        
        h2 {
            font-size: 18pt;
            font-weight: 600;
            color: #007cba;
            margin: 25px 0 12px 0;
            page-break-after: avoid;
        }
        
        h3 {
            font-size: 14pt;
            font-weight: 600;
            color: #333;
            margin: 20px 0 10px 0;
            page-break-after: avoid;
        }
        
        h4 {
            font-size: 12pt;
            font-weight: 600;
            color: #333;
            margin: 15px 0 8px 0;
            page-break-after: avoid;
        }
        
        p {
            margin: 12px 0;
            text-align: justify;
            orphans: 2;
            widows: 2;
        }
        
        /* Lists */
        ul, ol {
            margin: 15px 0;
            padding-left: 25px;
        }
        
        li {
            margin: 8px 0;
            line-height: 1.5;
        }
        
        /* Emphasis */
        strong {
            font-weight: 600;
            color: #007cba;
        }
        
        em {
            font-style: italic;
            color: #666;
        }
        
        /* Special elements */
        .highlight {
            background-color: #f0f8ff;
            padding: 15px;
            border-left: 4px solid #007cba;
            margin: 20px 0;
            border-radius: 4px;
        }
        
        .alert {
            padding: 15px;
            border-radius: 4px;
            margin: 20px 0;
            border-left: 4px solid;
        }
        
        .alert-danger {
            background: #f8d7da;
            border-left-color: #dc3545;
            color: #721c24;
        }
        
        .alert-warning {
            background: #fff3cd;
            border-left-color: #ffc107;
            color: #856404;
        }
        
        .alert-success {
            background: #d4edda;
            border-left-color: #28a745;
            color: #155724;
        }
        
        .alert-info {
            background: #d1ecf1;
            border-left-color: #17a2b8;
            color: #0c5460;
        }
        
        /* Page breaks */
        .page-break {
            page-break-before: always;
        }
        
        .no-break {
            page-break-inside: avoid;
        }
        
        /* Print-specific styles */
        @media print {
            body {
                font-size: 10pt;
                margin: 0;
                padding: 0;
            }
            
            h1 {
                font-size: 20pt;
            }
            
            h2 {
                font-size: 16pt;
            }
            
            h3 {
                font-size: 13pt;
            }
            
            .no-print {
                display: none !important;
            }
            
            /* Ensure colors print correctly */
            * {
                -webkit-print-color-adjust: exact !important;
                color-adjust: exact !important;
            }
        }
    </style>
</head>
<body>
    %s
</body>
</html>`, htmlContent)
	
	return enhancedHTML
}

// addPrintStyles adds print-specific styles to existing HTML
func (p *pdfService) addPrintStyles(htmlContent string) string {
	// Add print-specific styles if not already present
	printStyles := `
<style>
    @media print {
        body { font-size: 10pt; margin: 0; padding: 20px; }
        h1 { font-size: 20pt; page-break-after: avoid; }
        h2 { font-size: 16pt; page-break-after: avoid; }
        h3 { font-size: 13pt; page-break-after: avoid; }
        .no-print { display: none !important; }
        * { -webkit-print-color-adjust: exact !important; color-adjust: exact !important; }
    }
</style>`
	
	// Insert before closing head tag if it exists
	if strings.Contains(htmlContent, "</head>") {
		return strings.Replace(htmlContent, "</head>", printStyles+"\n</head>", 1)
	}
	
	// Otherwise, add at the beginning
	return printStyles + "\n" + htmlContent
}

// convertMarginToMM converts margin strings to millimeters
func (p *pdfService) convertMarginToMM(margin string) uint {
	if margin == "" {
		return 25 // Default 25mm
	}
	
	// Remove spaces
	margin = strings.TrimSpace(margin)
	
	// Handle different units
	if strings.HasSuffix(margin, "in") {
		// Convert inches to mm (1 inch = 25.4 mm)
		value := strings.TrimSuffix(margin, "in")
		if inches := p.parseFloat(value); inches > 0 {
			return uint(inches * 25.4)
		}
	} else if strings.HasSuffix(margin, "cm") {
		// Convert cm to mm (1 cm = 10 mm)
		value := strings.TrimSuffix(margin, "cm")
		if cm := p.parseFloat(value); cm > 0 {
			return uint(cm * 10)
		}
	} else if strings.HasSuffix(margin, "mm") {
		// Already in mm
		value := strings.TrimSuffix(margin, "mm")
		if mm := p.parseFloat(value); mm > 0 {
			return uint(mm)
		}
	} else {
		// Assume pixels and convert (rough approximation: 1px ≈ 0.26mm)
		if px := p.parseFloat(margin); px > 0 {
			return uint(px * 0.26)
		}
	}
	
	return 25 // Default fallback
}

// parseFloat safely parses a string to float64
func (p *pdfService) parseFloat(s string) float64 {
	// Simple float parsing - handle common decimal formats
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	
	// Handle common formats like "1", "1.0", "0.75"
	var result float64
	if _, err := fmt.Sscanf(s, "%f", &result); err == nil {
		return result
	}
	
	return 0
}

// getReportPDFOptions returns optimized PDF options for reports
func getReportPDFOptions() *interfaces.PDFOptions {
	return &interfaces.PDFOptions{
		PageSize:         "A4",
		Orientation:      "Portrait",
		MarginTop:        "1in",
		MarginRight:      "0.75in",
		MarginBottom:     "1in",
		MarginLeft:       "0.75in",
		EnableJavaScript: false,
		LoadTimeout:      30,
		Quality:          94,
		HeaderHTML: `
			<div style="text-align: center; font-size: 10px; color: #666; padding: 10px;">
				Cloud Consulting Report - Generated on ` + time.Now().Format("January 2, 2006") + `
			</div>
		`,
		FooterHTML: `
			<div style="text-align: center; font-size: 10px; color: #666; padding: 10px;">
				Page [page] of [topage] - Confidential
			</div>
		`,
		CustomOptions: map[string]string{
			"enable-local-file-access": "",
			"print-media-type":         "",
			"disable-smart-shrinking":  "",
		},
	}
}