package main

// func enforce_coverage() {
// 	file, err := os.Open("coverage.txt")
// 	if err != nil {
// 		fmt.Printf("Failed to open coverage.txt: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	var totalCoverage float64
// 	var lowCoverageFunctions []string

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		parts := strings.Fields(line)

// 		if len(parts) < 2 {
// 			continue // Skip malformed lines
// 		}

// 		if strings.HasPrefix(line, "total:") {
// 			// Extract total coverage
// 			if len(parts) >= 3 {
// 				totalCoverage, err = strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
// 				if err != nil {
// 					fmt.Printf("Failed to parse total coverage: %v\n", err)
// 					os.Exit(1)
// 				}
// 			}
// 			continue
// 		}

// 		// Extract function name and coverage
// 		functionName := parts[1]
// 		path := parts[0]
// 		coveragePercentStr := strings.TrimSuffix(parts[len(parts)-1], "%")
// 		coverage, err := strconv.ParseFloat(coveragePercentStr, 64)
// 		if err != nil {
// 			fmt.Printf("Failed to parse coverage for %s: %v\n", functionName, err)
// 			continue
// 		}

// 		if coverage < minCoverage {
// 			lowCoverageFunctions = append(lowCoverageFunctions, fmt.Sprintf("%s %s (%.2f%%)", path, functionName, coverage))
// 		}
// 	}

// 	if len(lowCoverageFunctions) > 0 {
// 		fmt.Println("Functions with low coverage:")
// 		for _, fn := range lowCoverageFunctions {
// 			// Remove module name
// 			parts := strings.SplitN(fn, "/", 2)
// 			if len(parts) > 1 {
// 				fn = parts[1]
// 			}
// 			fmt.Println(" -", fn)
// 		}
// 	}

// 	if totalCoverage < minCoverage {
// 		fmt.Printf("Test coverage is %.2f%%, which is below the required %.2f%%\n", totalCoverage, minCoverage)
// 		os.Exit(1)
// 	} else {
// 		fmt.Printf("Test coverage is %.2f%%, which meets the required %.2f%%\n", totalCoverage, minCoverage)
// 	}
// }
