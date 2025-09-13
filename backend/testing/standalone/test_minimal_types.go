package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fmt.Println("Parsing client_solution.go to check for syntax errors...")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "internal/interfaces/client_solution.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("File parsed successfully. Package: %s\n", node.Name.Name)

	// Count type declarations
	typeCount := 0
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					fmt.Printf("Found type: %s\n", typeSpec.Name.Name)
					typeCount++
				}
			}
		}
	}

	fmt.Printf("Total types found: %d\n", typeCount)
}
