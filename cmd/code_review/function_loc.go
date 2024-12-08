package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func function_loc() {
	failed := false
	// Parse all Go files in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the file
		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Failed to parse file %s: %v\n", path, err)
			return nil
		}

		// Inspect the AST
		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				// Calculate the number of lines in the function
				start := fs.Position(fn.Pos()).Line
				end := fs.Position(fn.End()).Line
				lines := end - start + 1

				// Check if the function has more than 50 lines
				if lines > maxLinesPerFunction {
					failed = true
					fmt.Printf("Function %s in %s has %d lines\n", fn.Name.Name, path, lines)
				}
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
