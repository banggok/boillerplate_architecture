package main

// const minCoverage = 80.0
const (
	dir                 = "./internal/"
	maxLinesPerFile     = 500 // Default threshold for large files in SonarQube
	maxLinesPerFunction = 75  // SonarQube suggests 75 as the maximum number of lines per function
	complexityThreshold = 15  // Cyclomatic complexity threshold in SonarQube (suggested maximum: 15)
	maxStructFields     = 10  // SonarQube generally flags structs with more than 10 fields
	maxInterfaceMethods = 10  // Increase to align with common interface design patterns
	maxFunctionParams   = 4   // SonarQube flags functions with 4 or more parameters
)
