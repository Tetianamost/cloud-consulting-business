package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Wait a moment for server to be ready
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8061"

	// Create a test inquiry with urgent language
	inquiry := map[string]interface{}{
		"name":     "Michael Chen",
		"email":    "michael.chen@urgenttech.com",
		"company":  "UrgentTech Solutions",
		"phone":    "+1-555-0299",
		"services": []string{"migration", "assessment"},
		"message":  "URGENT: Our production systems are experiencing critical performance issues. We need immediate cloud migration assistance. Can we schedule an emergency call today? This is blocking our entire operation.",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(inquiry)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Create inquiry
	fmt.Println("Creating urgent test inquiry...")
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
	fmt.Printf("✓ Created urgent inquiry with ID: %s\n", inquiryID)
	fmt.Println("✓ Check your email (info@cloudpartner.pro) for the improved internal notification!")
	fmt.Println("✓ The email should now have properly formatted report content with clean HTML structure.")
	
	fmt.Println("\nImproved email template test completed!")
}