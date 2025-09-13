package main

import (
	"fmt"

	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🔧 Simple Template Test...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create template service
	fmt.Println("Creating template service...")
	templateService := services.NewTemplateService("templates", logger)

	// Get available templates
	fmt.Println("Getting available templates...")
	templates := templateService.GetAvailableTemplates()

	fmt.Printf("Available templates: %v\n", templates)

	if len(templates) == 0 {
		fmt.Println("❌ No templates found!")
	} else {
		fmt.Printf("✅ Found %d templates\n", len(templates))
		for _, tmpl := range templates {
			fmt.Printf("   - %s\n", tmpl)
		}
	}
}
