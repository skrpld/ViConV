include .env
export

.PHONY: proto mongodb-up mongodb-down postgres-up postgres-down run docker-up docker-down

proto:
	@protoc -I./proto --go_out=./pkg/api/posts \
	--go-grpc_out=./pkg/api/posts --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative ./proto/posts.proto

mongodb-up:
	@migrate -database 'mongodb://${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DB}' -path migrations/mongodb up

mongodb-down:
	@migrate -database 'mongodb://${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DB}' -path migrations/mongodb down

postgres-up:
	@migrate -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}' -path migrations/postgres up

postgres-down:
	@migrate -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}' -path migrations/postgres down

run:
	@go run cmd/main.go

docker-up:
	@docker-compose up --build

docker-down:
	@docker-compose down -v