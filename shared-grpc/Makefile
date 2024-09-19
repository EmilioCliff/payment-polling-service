protoc:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

mock:
	mockgen -package mock -destination ./mock/mock.go github.com/EmilioCliff/payment-polling-service/shared-grpc/pb AuthenticationServiceClient

.PHONY: protoc mock