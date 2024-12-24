package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func file_loc() {
	failed := false
	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is a Go file
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Count the lines in the file
		lines, err := countLines(path)
		if err != nil {
			fmt.Printf("Error counting lines in file %s: %v\n", path, err)
			return nil
		}

		// Print the file if it exceeds the line threshold
		if lines > maxLinesPerFile {
			failed = true
			fmt.Printf("File: %s has %d lines. Max %d\n", path, lines, maxLinesPerFile)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	if failed {
		os.Exit(1)
	}
}

// countLines counts the number of lines in a file
func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}
