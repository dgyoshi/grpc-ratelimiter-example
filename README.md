# grpc ratelimiter example

run server
```
make run-server
```

terminal for Bob
```
CLIENT_NAME=Bob CALL_INTERVAL=100ms CALL_TIMES=100 go run github.com/dgyoshi/grpc-ratelimiter-example/cmd/client/
```

terminal for Carol
```
CLIENT_NAME=Carol CALL_INTERVAL=1s CALL_TIMES=100 go run github.com/dgyoshi/grpc-ratelimiter-example/cmd/client/
```