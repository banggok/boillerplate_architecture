package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func max_props() {
	failed := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			failed = checkFile(path)
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

func checkFile(filename string) bool {
	failed := false
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse file %s: %v\n", filename, err)
		return false
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			switch ts := spec.(type) {
			case *ast.TypeSpec:
				switch t := ts.Type.(type) {
				case *ast.StructType:
					numFields := len(t.Fields.List)
					if numFields > maxStructFields {
						failed = true
						pos := fset.Position(ts.Pos()) // Get the file position
						fmt.Printf("Struct %s in %s:%d has %d fields (max allowed: %d)\n",
							ts.Name.Name, pos.Filename, pos.Line, numFields, maxStructFields)
					}
				case *ast.InterfaceType:
					numMethods := len(t.Methods.List)
					if numMethods > maxInterfaceMethods {
						failed = true
						pos := fset.Position(ts.Pos()) // Get the file position
						fmt.Printf("Interface %s in %s:%d has %d methods (max allowed: %d)\n",
							ts.Name.Name, pos.Filename, pos.Line, numMethods, maxInterfaceMethods)
					}
				}
			}
		}
	}

	return failed
}
