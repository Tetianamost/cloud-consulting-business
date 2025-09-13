package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

func main() {
	fmt.Println("🔍 SIMPLE AWS SES CONNECTION TEST")
	fmt.Println("=================================")

	// Get credentials from environment
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_SES_REGION")
	senderEmail := os.Getenv("SES_SENDER_EMAIL")

	if accessKey == "" {
		fmt.Println("❌ AWS_ACCESS_KEY_ID not set")
		return
	}
	if secretKey == "" {
		fmt.Println("❌ AWS_SECRET_ACCESS_KEY not set")
		return
	}
	if region == "" {
		region = "us-east-1"
	}
	if senderEmail == "" {
		senderEmail = "info@cloudpartner.pro"
	}

	fmt.Printf("📧 Testing SES in region: %s\n", region)
	fmt.Printf("📧 Sender email: %s\n", senderEmail)
	fmt.Printf("🔑 Access Key: %s...\n", accessKey[:10])
	fmt.Println()

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		return
	}

	// Create SES client
	client := ses.NewFromConfig(cfg)

	// Test connection by getting sending quota
	ctx := context.Background()
	result, err := client.GetSendQuota(ctx, &ses.GetSendQuotaInput{})
	if err != nil {
		fmt.Printf("❌ SES connection failed: %v\n", err)
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("   - Invalid AWS credentials")
		fmt.Println("   - SES not available in this region")
		fmt.Println("   - Network connectivity issues")
		fmt.Println("   - AWS account doesn't have SES access")
		return
	}

	fmt.Printf("✅ SES connection successful!\n")
	fmt.Printf("   Max 24h send: %.0f emails\n", result.Max24HourSend)
	fmt.Printf("   Max send rate: %.2f emails/sec\n", result.MaxSendRate)
	fmt.Printf("   Sent last 24h: %.0f emails\n", result.SentLast24Hours)
	fmt.Println()

	// Check if we're in sandbox mode
	if result.Max24HourSend <= 200 {
		fmt.Println("⚠️  SES appears to be in SANDBOX MODE")
		fmt.Println("   - Can only send to verified email addresses")
		fmt.Println("   - Limited to 200 emails per 24 hours")
		fmt.Println("   - Need to request production access for real use")
	} else {
		fmt.Println("✅ SES is in PRODUCTION MODE")
		fmt.Println("   - Can send to any email address")
		fmt.Println("   - Higher sending limits")
	}

	fmt.Println()
	fmt.Println("🎯 NEXT STEPS:")
	fmt.Println("1. Verify sender email in AWS SES Console")
	fmt.Println("2. Test sending a real email")
	fmt.Println("3. Check email deliverability")
}
