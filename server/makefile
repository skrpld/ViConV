include .env

.PHONY: proto mongo-up mongo-down postgres-up postgres-down start

proto:
	@protoc -I./proto --go_out=./pkg/api/auth \
	--go-grpc_out=./pkg/api/auth --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative ./proto/auth.proto

mongo-up:
	@migrate-mongo up

mongo-down:
	@migrate-mongo down

postgres-up:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE} -path migrations up

postgres-down:
	@migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE} -path migrations down

start:
	@go run cmd/main.go