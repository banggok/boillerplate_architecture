package main

const (
	dir                 = "./internal/" // Directory to review
	minCoverage         = 80.0          // Minimum Testing Coverage, Unit and E2E Test combining
	maxLinesPerFile     = 500           // Maximum Line of Code per file
	maxLinesPerFunction = 75            // Maximum Line of Code per function
	complexityThreshold = 15            // Cyclomatic complexity threshold
	maxStructFields     = 10            // Maximum fields in a struct
	maxInterfaceMethods = 10            // Maximum methods in an interface
	maxFunctionParams   = 4             // Maximum parameter in a function
)
