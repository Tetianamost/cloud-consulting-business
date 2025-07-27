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

type CreateInquiryRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Company  string   `json:"company"`
	Phone    string   `json:"phone"`
	Services []string `json:"services"`
	Message  string   `json:"message"`
}

type InquiryResponse struct {
	Success bool   `json:"success"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
}

func main() {
	baseURL := "http://localhost:8061"
	
	fmt.Println("=== PDF Endpoints Test ===")
	
	// Test 1: Create an inquiry first
	fmt.Println("\n1. Creating test inquiry...")
	
	inquiry := CreateInquiryRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Company:  "Test Company",
		Phone:    "+1-555-0123",
		Services: []string{"assessment"},
		Message:  "This is a test inquiry for PDF generation testing.",
	}
	
	jsonData, err := json.Marshal(inquiry)
	if err != nil {
		fmt.Printf("❌ Failed to marshal inquiry: %v\n", err)
		return
	}
	
	resp, err := http.Post(baseURL+"/api/v1/inquiries", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create inquiry: %v\n", err)
		fmt.Println("Make sure the server is running: go run cmd/server/main.go")
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Failed to create inquiry - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	var inquiryResp InquiryResponse
	err = json.NewDecoder(resp.Body).Decode(&inquiryResp)
	if err != nil {
		fmt.Printf("❌ Failed to decode inquiry response: %v\n", err)
		return
	}
	
	inquiryID := inquiryResp.Data.ID
	fmt.Printf("✅ Inquiry created successfully - ID: %s\n", inquiryID)
	
	// Wait a moment for report generation
	fmt.Println("⏳ Waiting for report generation...")
	time.Sleep(3 * time.Second)
	
	// Test 2: Test PDF endpoint
	fmt.Println("\n2. Testing PDF endpoint...")
	
	pdfURL := fmt.Sprintf("%s/api/v1/inquiries/%s/report/pdf", baseURL, inquiryID)
	resp, err = http.Get(pdfURL)
	if err != nil {
		fmt.Printf("❌ Failed to get PDF: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ PDF endpoint failed - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		fmt.Printf("❌ Wrong content type - Expected: application/pdf, Got: %s\n", contentType)
		return
	}
	
	// Read PDF content
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read PDF data: %v\n", err)
		return
	}
	
	// Save PDF
	err = os.WriteFile("test_endpoint_pdf.pdf", pdfData, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save PDF: %v\n", err)
		return
	}
	
	fmt.Printf("✅ PDF endpoint works - Size: %d bytes, Content-Type: %s\n", len(pdfData), contentType)
	
	// Test 3: Test download endpoint with PDF format
	fmt.Println("\n3. Testing download endpoint (PDF format)...")
	
	downloadURL := fmt.Sprintf("%s/api/v1/inquiries/%s/report/download?format=pdf", baseURL, inquiryID)
	resp, err = http.Get(downloadURL)
	if err != nil {
		fmt.Printf("❌ Failed to download PDF: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Download endpoint failed - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	// Check headers
	contentType = resp.Header.Get("Content-Type")
	contentDisposition := resp.Header.Get("Content-Disposition")
	
	if contentType != "application/pdf" {
		fmt.Printf("❌ Wrong content type - Expected: application/pdf, Got: %s\n", contentType)
		return
	}
	
	if contentDisposition == "" {
		fmt.Printf("❌ Missing Content-Disposition header\n")
		return
	}
	
	// Read download data
	downloadData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read download data: %v\n", err)
		return
	}
	
	// Save download
	err = os.WriteFile("test_download_pdf.pdf", downloadData, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save download: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Download endpoint works - Size: %d bytes\n", len(downloadData))
	fmt.Printf("   Content-Type: %s\n", contentType)
	fmt.Printf("   Content-Disposition: %s\n", contentDisposition)
	
	// Test 4: Test download endpoint with HTML format
	fmt.Println("\n4. Testing download endpoint (HTML format)...")
	
	htmlDownloadURL := fmt.Sprintf("%s/api/v1/inquiries/%s/report/download?format=html", baseURL, inquiryID)
	resp, err = http.Get(htmlDownloadURL)
	if err != nil {
		fmt.Printf("❌ Failed to download HTML: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ HTML download endpoint failed - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	// Check headers
	contentType = resp.Header.Get("Content-Type")
	contentDisposition = resp.Header.Get("Content-Disposition")
	
	if contentType != "text/html; charset=utf-8" {
		fmt.Printf("❌ Wrong content type - Expected: text/html; charset=utf-8, Got: %s\n", contentType)
		return
	}
	
	// Read HTML data
	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read HTML data: %v\n", err)
		return
	}
	
	// Save HTML
	err = os.WriteFile("test_download_report.html", htmlData, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save HTML: %v\n", err)
		return
	}
	
	fmt.Printf("✅ HTML download endpoint works - Size: %d bytes\n", len(htmlData))
	fmt.Printf("   Content-Type: %s\n", contentType)
	fmt.Printf("   Content-Disposition: %s\n", contentDisposition)
	
	// Test 5: Test admin download endpoint (requires authentication)
	fmt.Println("\n5. Testing admin download endpoint...")
	
	// First, login to get auth token
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	
	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("❌ Failed to marshal login data: %v\n", err)
		return
	}
	
	resp, err = http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Printf("❌ Failed to login: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Login failed - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	var loginResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		fmt.Printf("❌ Failed to decode login response: %v\n", err)
		return
	}
	
	token, ok := loginResp["token"].(string)
	if !ok {
		fmt.Printf("❌ No token in login response\n")
		return
	}
	
	// Test admin PDF download
	adminPDFURL := fmt.Sprintf("%s/api/v1/admin/reports/%s/download/pdf", baseURL, inquiryID)
	req, err := http.NewRequest("GET", adminPDFURL, nil)
	if err != nil {
		fmt.Printf("❌ Failed to create admin request: %v\n", err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to get admin PDF: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Admin PDF endpoint failed - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return
	}
	
	// Read admin PDF data
	adminPDFData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read admin PDF data: %v\n", err)
		return
	}
	
	// Save admin PDF
	err = os.WriteFile("test_admin_pdf.pdf", adminPDFData, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save admin PDF: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Admin PDF endpoint works - Size: %d bytes\n", len(adminPDFData))
	
	fmt.Println("\n=== PDF Endpoints Test Summary ===")
	fmt.Println("✅ Inquiry Creation: PASSED")
	fmt.Println("✅ PDF Endpoint: PASSED")
	fmt.Println("✅ Download Endpoint (PDF): PASSED")
	fmt.Println("✅ Download Endpoint (HTML): PASSED")
	fmt.Println("✅ Admin Download Endpoint: PASSED")
	fmt.Println("\nGenerated test files:")
	fmt.Println("- test_endpoint_pdf.pdf")
	fmt.Println("- test_download_pdf.pdf")
	fmt.Println("- test_download_report.html")
	fmt.Println("- test_admin_pdf.pdf")
	fmt.Println("\n🎉 All PDF endpoint tests PASSED!")
}