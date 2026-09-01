all: clean build

.PHONY: all clean build run

clean: 
	rm -rf main/proto build

main/proto/message-api.pb.go: proto/message-api.proto
	protoc  --go_out=.  \
	        --go-grpc_out=. \
	        proto/message-api.proto

build: main/service.go main/proto/message-api.pb.go
	go build -o ./build/service ./main

run: build
	./build/service