package main

import (
	"fmt"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

func main() {
	fmt.Println("Testing client solution types...")

	// Test WorkloadRequirements
	req := &interfaces.WorkloadRequirements{
		WorkloadType: "web-application",
	}

	// Test WorkloadOptimization
	opt := &interfaces.WorkloadOptimization{
		WorkloadType: "web-application",
	}

	// Test MigrationRequest
	migReq := &interfaces.MigrationRequest{
		SourceType: "on-premises",
		TargetType: "aws",
	}

	// Test MigrationStrategy
	migStrat := &interfaces.MigrationStrategy{
		MigrationName: "Test Migration",
	}

	// Test DRRequirements
	drReq := &interfaces.DRRequirements{
		BusinessCriticality: "high",
	}

	// Test BCPRequirements
	bcpReq := &interfaces.BCPRequirements{
		BusinessFunctions: []string{"core-operations"},
	}

	fmt.Printf("✓ WorkloadRequirements: %s\n", req.WorkloadType)
	fmt.Printf("✓ WorkloadOptimization: %s\n", opt.WorkloadType)
	fmt.Printf("✓ MigrationRequest: %s -> %s\n", migReq.SourceType, migReq.TargetType)
	fmt.Printf("✓ MigrationStrategy: %s\n", migStrat.MigrationName)
	fmt.Printf("✓ DRRequirements: %s\n", drReq.BusinessCriticality)
	fmt.Printf("✓ BCPRequirements: %d functions\n", len(bcpReq.BusinessFunctions))

	fmt.Println("All types are accessible!")
}
