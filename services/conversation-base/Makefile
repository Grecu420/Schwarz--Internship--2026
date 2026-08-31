all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc --go_out=. \
	       --go-grpc_out=. \
	       proto/conversation-api.proto

build: proto
	go build -o ./build/service ./main

run: build
	./build/service