package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func max_params() {
	failed := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fileFailed := checkFileForParams(path)
			if fileFailed {
				failed = true
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", dir, err)
	}

	if failed {
		os.Exit(1)
	}
}

func checkFileForParams(filename string) bool {
	failed := false
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Failed to parse file %s: %v\n", filename, err)
		return false
	}

	// Walk through the AST to find function declarations
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Count the parameters
		numParams := countFunctionParams(funcDecl)
		if numParams > maxFunctionParams {
			failed = true
			pos := fset.Position(funcDecl.Pos()) // Get the file position
			fmt.Printf("Function %s in %s:%d has %d parameters (max allowed: %d)\n",
				funcDecl.Name.Name, pos.Filename, pos.Line, numParams, maxFunctionParams)
		}
	}

	return failed
}

func countFunctionParams(funcDecl *ast.FuncDecl) int {
	if funcDecl.Type.Params == nil {
		return 0
	}

	count := 0
	for _, field := range funcDecl.Type.Params.List {
		// Each field can represent one or more parameters (e.g., `a, b int`)
		if len(field.Names) > 0 {
			count += len(field.Names)
		} else {
			// Anonymous parameters still count as one each
			count++
		}
	}

	return count
}
