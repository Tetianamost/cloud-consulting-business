package main

import (
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🔧 Testing Email Service Fix...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Check SES configuration
	fmt.Printf("📧 SES Configuration Check:\n")
	fmt.Printf("   Access Key ID: %s\n", maskString(cfg.SES.AccessKeyID))
	fmt.Printf("   Secret Key: %s\n", maskString(cfg.SES.SecretAccessKey))
	fmt.Printf("   Region: %s\n", cfg.SES.Region)
	fmt.Printf("   Sender Email: %s\n", cfg.SES.SenderEmail)

	// Test if configuration is complete
	if cfg.SES.AccessKeyID == "" || cfg.SES.SecretAccessKey == "" || cfg.SES.SenderEmail == "" {
		fmt.Println("❌ SES configuration is incomplete!")
		return
	}

	fmt.Println("✅ SES configuration is complete!")

	// Test email service initialization
	fmt.Println("\n🚀 Testing Email Service Initialization...")
	emailService, err := services.NewEmailServiceWithSES(cfg.SES, logger)
	if err != nil {
		fmt.Printf("❌ Failed to initialize email service: %v\n", err)
		return
	}

	fmt.Println("✅ Email service initialized successfully!")

	// Test email service health
	fmt.Println("\n🏥 Testing Email Service Health...")
	if emailService.IsHealthy() {
		fmt.Println("✅ Email service is healthy!")
	} else {
		fmt.Println("❌ Email service is not healthy!")
	}

	fmt.Println("\n🎉 Email service fix verification completed!")
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
