all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc -I=..\
		   --go_out=. \
	       property-base/proto/property-api.proto
	protoc -I=..\
		   --go_out=. \
	       --go-grpc_out=. \
	       property-list-base/proto/property-list-api.proto

build: proto
	go build -o ./build/service ./main

run: build
	./build/service