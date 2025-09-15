package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("🔧 Direct Template Loading Test...")

	// Test direct template loading
	templatePath := filepath.Join("templates", "email", "consultant_notification.html")

	fmt.Printf("Checking template path: %s\n", templatePath)

	// Check if file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("❌ Template file does not exist: %s\n", templatePath)
		return
	}

	fmt.Println("✅ Template file exists")

	// Try to parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("❌ Failed to parse template: %v\n", err)
		return
	}

	fmt.Println("✅ Template parsed successfully")

	// Test with simple data
	data := map[string]interface{}{
		"Name":           "Test User",
		"Email":          "test@example.com",
		"Company":        "Test Company",
		"Phone":          "555-0123",
		"Services":       "assessment, migration",
		"Message":        "Test message",
		"ID":             "test-123",
		"IsHighPriority": true,
		"Priority":       "HIGH",
		"Report": map[string]interface{}{
			"ID":          "report-123",
			"Content":     "Test report content",
			"HTMLContent": template.HTML("<p>Test HTML content</p>"),
		},
	}

	// Try to execute the template
	fmt.Println("Testing template execution...")
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("❌ Failed to execute template: %v\n", err)
		return
	}

	fmt.Println("\n✅ Template executed successfully!")
}
