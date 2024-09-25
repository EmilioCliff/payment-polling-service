SERVICE=gateway-service

test:
	@echo "Running tests for $(SERVICE)..."

	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out 
	go tool cover -html=coverage.out -o coverage.html

race-test:
	@echo "Running race-tests for $(SERVICE)..."

	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

revive:
	@echo "Running revive for $(SERVICE)..."

	revive -formatter friendly ./...

.PHONY: test race-test lint revive