package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type lowCover struct {
	Path     string
	Coverage float64
}

func enforce_coverage() {
	file, err := os.Open("coverage.txt")
	if err != nil {
		fmt.Printf("Failed to open coverage.txt: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var totalCoverage float64
	var lowCoverageFunctions []lowCover

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue // Skip malformed lines
		}

		if strings.HasPrefix(line, "total:") {
			// Extract total coverage
			if len(parts) >= 3 {
				totalCoverage, err = strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
				if err != nil {
					fmt.Printf("Failed to parse total coverage: %v\n", err)
					os.Exit(1)
				}
			}
			continue
		}

		// Extract function name and coverage
		functionName := parts[1]
		path := parts[0]
		coveragePercentStr := strings.TrimSuffix(parts[len(parts)-1], "%")
		coverage, err := strconv.ParseFloat(coveragePercentStr, 64)
		if err != nil {
			fmt.Printf("Failed to parse coverage for %s: %v\n", functionName, err)
			continue
		}

		if coverage < minCoverage {
			lowCoverageFunctions = append(lowCoverageFunctions, lowCover{
				Path:     fmt.Sprintf("%s %s (%.2f%%)", path, functionName, coverage),
				Coverage: coverage,
			})
		}
	}

	if len(lowCoverageFunctions) > 0 {
		sortLowCoverageAsc(lowCoverageFunctions)
		fmt.Println("Functions with low coverage:")
		for _, fn := range lowCoverageFunctions {
			// Remove module name
			parts := strings.SplitN(fn.Path, "/", 2)
			if len(parts) > 1 {
				fn.Path = parts[1]
			}
			fmt.Println(" -", fn.Path)
		}
	}

	if totalCoverage < minCoverage {
		fmt.Printf("Test coverage is %.2f%%, which is below the required %.2f%%\n", totalCoverage, minCoverage)
		os.Exit(1)
	} else {
		fmt.Printf("Test coverage is %.2f%%, which meets the required %.2f%%\n", totalCoverage, minCoverage)
	}
}

func sortLowCoverageAsc(lowCoverageFunctions []lowCover) {
	sort.Slice(lowCoverageFunctions, func(i, j int) bool {
		return lowCoverageFunctions[i].Coverage < lowCoverageFunctions[j].Coverage
	})
}
