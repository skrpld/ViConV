.PHONY: proto mongo-up mongo-down start

proto:
	@protoc -I./proto --go_out=./pkg/api/auth \
	--go-grpc_out=./pkg/api/auth --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative ./proto/auth.proto

mongo-up:
	@migrate-mongo up

mongo-down:
	@migrate-mongo down

start:
	@go run cmd/main.go