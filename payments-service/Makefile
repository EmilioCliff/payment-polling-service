sqlc:
	cd internal/postgres && sqlc generate

asynq:
	docker run --name asynq -p 8383:8383 --network payment-polling-service_default hibiken/asynqmon --redis-addr=redis:6379 --redis-db=1 --port 8383

.PHONY: sqlc asynq