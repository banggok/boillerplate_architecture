.PHONY: code-review run

code-review:
	go test ./internal/... -race -coverprofile=coverage.out && go tool cover -func=coverage.out > coverage.txt
	
	@echo "Running nestif to check for nested if statements with complexity > 3..."
	@output=$$(nestif --min 4 ./internal/...); \
	if [ -n "$$output" ]; then \
		echo "$$output"; \
		echo "Error: Detected nested if statements with complexity greater than 3."; \
		exit 1; \
	fi
	
	go run ./cmd/code_review/

run:
	go run ./cmd/app/
