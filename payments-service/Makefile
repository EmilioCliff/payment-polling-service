sqlc:
	cd internal/postgres && sqlc generate

asynq:
	docker run --name asynq -p 8383:8383 --network payment-polling-service_default hibiken/asynqmon --redis-addr=redis:6379 --redis-db=1 --port 8383

mock:
	mockgen -package mockdb -destination ./internal/mock/store.go github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres Store

.PHONY: sqlc asynq mock