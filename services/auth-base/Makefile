

all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc 	--go_out=.  \
    		--go-grpc_out=. \
		 	proto/auth-api.proto

main/proto/auth-api.pb.go: proto/auth-api.proto
	protoc 	-I=.. \
			--go_out=.  \
    		--go-grpc_out=. \
		 	./auth-base/proto/auth-api.proto


build: main/service.go main/proto/auth-api.pb.go
	go build -o ./build/service ./main

run: build/service
	./build/service

