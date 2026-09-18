all: clean build

.PHONY: all clean build run

clean: 
	rm -rf main/proto build

main/proto/reservation-api.pb.go: proto/reservation-api.proto
	protoc  --go_out=.  \
	        --go-grpc_out=. \
	        proto/reservation-api.proto

build: main/service.go main/proto/reservation-api.pb.go
	go build -o ./build/service ./main

run: build
	./build/service