package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

func main() {
	fmt.Println("Testing email MIME structure...")

	// Create a test email message
	email := &interfaces.EmailMessage{
		From:     "noreply@cloudpartner.pro",
		To:       []string{"test@example.com"},
		Subject:  "Test Email with HTML Formatting",
		ReplyTo:  "info@cloudpartner.pro",
		TextBody: "This is the plain text version of the email.",
		HTMLBody: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test Email</title>
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .header { background: #667eea; color: white; padding: 20px; }
        .content { padding: 20px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Test Email</h1>
    </div>
    <div class="content">
        <p>This is a <strong>test email</strong> with HTML formatting.</p>
        <p>It should render properly in email clients.</p>
    </div>
</body>
</html>`,
	}

	// Use reflection to access the private buildRawMessage method
	// We'll create a mock that captures the raw message
	mockSES := &MIMECaptureSESService{}
	
	// Send the email through our mock to capture the MIME structure
	err := mockSES.SendEmail(context.Background(), email)
	if err != nil {
		fmt.Printf("Failed to send test email: %v\n", err)
		return
	}

	fmt.Println("✅ MIME structure test completed successfully!")
	fmt.Println("Check the generated .eml file to see the actual MIME structure.")
}

// MIMECaptureSESService captures the raw MIME message for inspection
type MIMECaptureSESService struct{}

func (m *MIMECaptureSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Build the raw message using the same logic as the real SES service
	rawMessage, err := m.buildRawMessage(email)
	if err != nil {
		return fmt.Errorf("failed to build raw email message: %w", err)
	}

	// Write the raw MIME message to a file for inspection
	filename := fmt.Sprintf("test_email_mime_%d.eml", time.Now().Unix())
	err = os.WriteFile(filename, rawMessage, 0644)
	if err != nil {
		return fmt.Errorf("failed to write MIME file: %w", err)
	}

	fmt.Printf("Raw MIME message written to: %s\n", filename)
	fmt.Printf("MIME message size: %d bytes\n", len(rawMessage))
	
	// Show a preview of the MIME structure
	lines := strings.Split(string(rawMessage), "\n")
	fmt.Println("\nMIME Structure Preview (first 30 lines):")
	fmt.Println("=" + strings.Repeat("=", 50))
	for i, line := range lines {
		if i >= 30 {
			fmt.Println("... (truncated)")
			break
		}
		fmt.Println(line)
	}
	fmt.Println("=" + strings.Repeat("=", 50))

	return nil
}

func (m *MIMECaptureSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *MIMECaptureSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     1,
		SentLast24Hours: 0,
	}, nil
}

// buildRawMessage builds a raw MIME email message (copied from ses.go)
func (m *MIMECaptureSESService) buildRawMessage(email *interfaces.EmailMessage) ([]byte, error) {
	var buf strings.Builder

	// Write headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	if email.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: multipart/alternative; boundary=boundary123\r\n")
	buf.WriteString("\r\n")

	// Add text part
	if email.TextBody != "" {
		buf.WriteString("--boundary123\r\n")
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: 7bit\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(email.TextBody)
		buf.WriteString("\r\n")
	}

	// Add HTML part
	if email.HTMLBody != "" {
		buf.WriteString("--boundary123\r\n")
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: 7bit\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(email.HTMLBody)
		buf.WriteString("\r\n")
	}

	// Close boundary
	buf.WriteString("--boundary123--\r\n")

	return []byte(buf.String()), nil
}