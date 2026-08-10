

all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc 	--go_out=.  \
    		--go-grpc_out=. \
		 	proto/user-api.proto

main/proto/user-api.pb.go: proto/user-api.proto
	protoc 	-I=.. \
			--go_out=.  \
    		--go-grpc_out=. \
		 	./user-base/proto/user-api.proto


build: main/service.go main/proto/user-api.pb.go
	go build -o ./build/service ./main

run: build/service
	./build/service

