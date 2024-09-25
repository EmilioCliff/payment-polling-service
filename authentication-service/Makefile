SERVICE=authentication-service

createMigrate:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

sqlc:
	cd .envs/configs && sqlc generate

mock:
	mockgen -package mockdb -destination ./internal/postgres/mock/store.go github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated Querier

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

.PHONY: createMigrate migrateup migratedown sqlc mock test race-test revive