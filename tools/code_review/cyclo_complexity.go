package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

var Analyzer = &analysis.Analyzer{
	Name: "cyclomatic",
	Doc:  "checks cyclomatic complexity of functions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			complexity := calculateCyclomaticComplexity(fn)
			if complexity > complexityThreshold {
				pass.Reportf(fn.Pos(), "function %q has high cyclomatic complexity: %d", fn.Name.Name, complexity)
			}
			return false
		})
	}
	return nil, nil
}

// calculateCyclomaticComplexity computes the cyclomatic complexity of a function
func calculateCyclomaticComplexity(fn *ast.FuncDecl) int {
	// Start with a base complexity of 1
	complexity := 1

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n := n.(type) { // Assigns type assertion to a variable
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complexity++
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complexity++
			}
		}
		return true
	})

	return complexity
}

func cyclomatic_complexity() {

	// Load the Go files in the specified directory
	fs := token.NewFileSet()
	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedTypes | packages.NeedName,
		Fset:  fs,
		Tests: false,
		Dir:   dir,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load packages: %v\n", err)
		os.Exit(1)
	}

	// Analyze each package and file
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				complexity := calculateCyclomaticComplexity(fn)
				if complexity > complexityThreshold {
					pos := fs.Position(fn.Pos())
					fmt.Printf("Function %q in %s:%d has high cyclomatic complexity: %d. Max: %d\n",
						fn.Name.Name, pos.Filename, pos.Line, complexity, complexityThreshold)
				}
				return false
			})
		}
	}
}
