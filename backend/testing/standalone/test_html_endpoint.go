package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// Wait a moment for server to be ready
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8061"

	// Create a test inquiry
	inquiry := map[string]interface{}{
		"name":     "Jane Smith",
		"email":    "jane.smith@example.com",
		"company":  "Innovation Labs",
		"phone":    "+1-555-0199",
		"services": []string{"assessment"},
		"message":  "We need urgent cloud assessment for our infrastructure. Can we schedule a meeting today?",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(inquiry)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Create inquiry
	fmt.Println("Creating test inquiry...")
	resp, err := http.Post(baseURL+"/api/v1/inquiries", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating inquiry: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to create inquiry. Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}

	// Parse response to get inquiry ID
	var createResponse struct {
		Success bool `json:"success"`
		Data    struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	err = json.Unmarshal(body, &createResponse)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	if !createResponse.Success {
		fmt.Printf("Inquiry creation failed: %s\n", string(body))
		return
	}

	inquiryID := createResponse.Data.ID
	fmt.Printf("✓ Created inquiry with ID: %s\n", inquiryID)

	// Wait for report generation
	fmt.Println("Waiting for report generation...")
	time.Sleep(5 * time.Second)

	// Test HTML report endpoint
	fmt.Println("Fetching HTML report...")
	htmlResp, err := http.Get(baseURL + "/api/v1/inquiries/" + inquiryID + "/report/html")
	if err != nil {
		fmt.Printf("Error fetching HTML report: %v\n", err)
		return
	}
	defer htmlResp.Body.Close()

	if htmlResp.StatusCode != http.StatusOK {
		htmlBody, _ := io.ReadAll(htmlResp.Body)
		fmt.Printf("Failed to fetch HTML report. Status: %d, Body: %s\n", htmlResp.StatusCode, string(htmlBody))
		return
	}

	// Save HTML report to file
	htmlContent, err := io.ReadAll(htmlResp.Body)
	if err != nil {
		fmt.Printf("Error reading HTML content: %v\n", err)
		return
	}

	filename := fmt.Sprintf("test_endpoint_report_%s.html", inquiryID)
	err = os.WriteFile(filename, htmlContent, 0644)
	if err != nil {
		fmt.Printf("Error saving HTML file: %v\n", err)
		return
	}

	fmt.Printf("✓ HTML report saved to: %s\n", filename)
	fmt.Printf("✓ HTML report size: %d bytes\n", len(htmlContent))

	// Check if content looks like HTML
	htmlStr := string(htmlContent)
	if len(htmlStr) > 100 && 
		(bytes.Contains(htmlContent, []byte("<!DOCTYPE html")) || 
		 bytes.Contains(htmlContent, []byte("<html"))) {
		fmt.Println("✓ Response appears to be valid HTML")
	} else {
		fmt.Println("⚠ Response may not be valid HTML")
		fmt.Printf("First 200 chars: %s\n", htmlStr[:min(200, len(htmlStr))])
	}

	fmt.Println("\nHTML report endpoint test completed successfully!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}