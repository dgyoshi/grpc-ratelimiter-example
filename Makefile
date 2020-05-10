PROTO_DIR := internal/pkg/pb
protogen:
	mkdir -p $(PROTO_DIR)/greeting
	protoc --proto_path=./proto/greeting --go_out=plugins=grpc:$(PROTO_DIR)/greeting greeting.proto

run-server:
	go run github.com/dgyoshi/grpc-ratelimiter-example/cmd/greeting/

run-bob:
	CLIENT_NAME=Bob CALL_INTERVAL=1s go run github.com/dgyoshi/grpc-ratelimiter-example/cmd/client/
