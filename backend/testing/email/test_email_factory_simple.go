package main

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("=== Email Service Factory Simple Test ===")

	// Test configuration
	cfg := config.SESConfig{
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@cloudpartner.pro",
		ReplyToEmail:    "reply@cloudpartner.pro",
		Timeout:         30,
	}

	// Test 1: Create email service without event recorder
	fmt.Println("1. Testing email service creation without event recorder...")
	emailService, err := services.NewEmailServiceFactory(cfg, nil, logger)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}

	if emailService == nil {
		log.Fatal("Email service is nil")
	}

	fmt.Println("✅ Email service created successfully")

	// Test 2: Check health
	fmt.Println("2. Testing email service health...")
	if !emailService.IsHealthy() {
		log.Fatal("Email service health check failed")
	}

	fmt.Println("✅ Email service health check passed")

	// Test 3: Test configuration validation
	fmt.Println("3. Testing configuration validation...")
	invalidCfg := config.SESConfig{
		AccessKeyID:     "", // Missing required field
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@cloudpartner.pro",
		Timeout:         30,
	}

	_, err = services.NewEmailServiceFactory(invalidCfg, nil, logger)
	if err == nil {
		log.Fatal("Expected error with invalid configuration")
	}

	fmt.Printf("✅ Invalid configuration correctly rejected: %v\n", err)

	fmt.Println("\n=== Email Service Factory Test Completed Successfully! ===")
	fmt.Println("✅ Email service factory is working correctly")
	fmt.Println("✅ Configuration validation is working")
	fmt.Println("✅ Health checks are functional")
	fmt.Println("✅ Task 13 implementation is complete!")
}
