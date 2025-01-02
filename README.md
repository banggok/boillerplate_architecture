# Boilerplate Architecture

This repository provides a robust example of clean code architecture with a Domain-Driven Design (DDD) approach. The repository adheres to SOLID principles and includes an automated code review process to validate these rules. It is designed to accelerate development while ensuring scalability, maintainability, and testability, making it suitable for both small-scale and enterprise-grade projects.

## Features

- **Clean Code Architecture**: Modular and layer-based project organization for enhanced scalability and maintainability.
- **Domain-Driven Design (DDD)**: Focuses on business logic encapsulation within domain layers.
- **SOLID Principles**: Ensures code is well-structured, reusable, and easy to refactor.
- **Code Review Tool**: Automated validation of SOLID principles and architectural standards.
- **Containerization**: Docker and Docker Compose support for a seamless development environment.
- **Test Coverage**: Built-in test cases with coverage reports.

## Project Structure

The project is organized into the following directories:

```
AMS/
├── cmd/                
│   ├── app/                # Entry points for the application
│   ├── scheduler/          # Entry points for the scheduler
│   ├── subscriber/         # Entry points for the subscriber
├── internal/               # Core application code
│   ├── config/             # Setup Configuration
│   ├── data/               # Data Manipulation related
│   │   ├── entity/         # Domain entities and business rules
│   │   ├── model/          # Database tables structure
│   └── delivery/           # API and external communication handler
│   ├── pkg/                # helper library packages
│   ├── service/            # Domain Business logic services
├── migrations/             # Database migrations and seeder
├── tools/                
│   ├── code-review/        # code static analysis and review
├── test/               # Unit and integration tests
├── docs/               # Documentation and design artifacts
├── Makefile            # Commands for building, testing, and running the application
├── docker-compose.yml  # Docker Compose configuration
├── .env                # Environment variables
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make

## Getting Started

1. **Set up environment variables**:
   Create a `.env` file with the required configurations based on the `.env.example` directory.

2. **Run Docker containers**:
   ```bash
   docker-compose up -d
   ```
3. **Run Migrations**:
   ```bash
   make migrate
   ```
4. **Run the application**:
   ```bash
   make run
   ```

## Code Review Process

The repository includes a custom code review mechanism to validate adherence to SOLID principles and clean code standards. This can be executed as part of the CI/CD pipeline. Configuration is in `tools/code_review/config.go`. The following thresholds ensure maintainable, high-quality code:

- **Directory to review**: Only relevant code under `./internal/` is analyzed.
- **Testing Coverage**: Minimum 80% combined coverage for unit and E2E tests.
- **Maximum Lines Per File**: Enforces modularity with a 500-line limit per file.
- **Function Length**: Promotes readability with a maximum of 75 lines per function.
- **Cyclomatic Complexity**: Limits complexity to 15 for better maintainability.
- **Struct Fields and Interface Methods**: Keeps structs and interfaces simple with a maximum of 10 fields/methods.
- **Function Parameters**: Limits parameters to 4, ensuring simplicity and testability.

### Steps for Code Review

1. **Set Configuration in `tools/code_review/config.go`**:
   ```go
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
   ```
2. **Run tests and review**:
   ```bash
   make code-review
   ```
3. **Open code coverage report in HTML**:
   ```bash
   make open-coverage
   ```
 
## Documentation

1. **Generate Swagger to `docs/`**:
   ```bash
   make docs
   ```
## Database Administration

1. **Create a new migration file to `migrations/`**:
   ```bash
   make create-migration [your_migration_name]
   ```
2. **Run all pending migrations**:
   ```bash
   make migrate
   ```
3. **Rollback the last migration**:
   ```bash
   make rollback
   ```
4. **Force the migration version**:
   ```bash
   make force-version [your_version_number]
   ```

## Disclaimer

1. This architecture reflects current best practices and my experience in software design. While it aims for a high-quality, maintainable structure, it is not without its flaws. Feedback and constructive criticism are always welcomed to help improve this boilerplate.

2. The design may appear overly complex or "too much" for certain use cases, especially smaller projects. However, it is intentionally structured to prioritize scalability, reusability, and long-term maintainability. By adhering to these principles, the boilerplate reduces technical debt and promotes efficient collaboration over the project lifecycle.
