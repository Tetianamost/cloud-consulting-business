package main

import (
	"fmt"
	"reflect"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

func main() {
	fmt.Println("Testing direct type access...")

	// Try to get the type information
	var req *interfaces.WorkloadRequirements
	if req != nil {
		fmt.Printf("WorkloadRequirements type: %v\n", reflect.TypeOf(req))
	} else {
		fmt.Println("WorkloadRequirements is nil but type exists")
	}
}
