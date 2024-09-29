SERVICE=payment-service

sqlc:
	cd .envs/configs && sqlc generate

createMigrate:
	migrate create -ext sql -dir internal/postgres/migrations -seq transaction_message

asynq:
	docker run --name asynq -p 8383:8383 --network payment-polling-service_default hibiken/asynqmon --redis-addr=redis:6379 --redis-db=1 --port 8383

mock:
	mockgen -package mock -destination ./internal/postgres/mock/store.go github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/generated Querier

test:
	@echo "Running tests for $(SERVICE)..."

	go test -v ./...

race-test:
	@echo "Running race-tests for $(SERVICE)..."

	go test -v -race ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out 
	go tool cover -html=coverage.out -o coverage.html

revive:
	@echo "Running revive for $(SERVICE)..."

	revive -formatter friendly ./...

.PHONY: sqlc asynq mock createMigrate test race-test coverage revive