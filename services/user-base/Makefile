

all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc 	-I=.. \
			--go_out=.  \
    		--go-grpc_out=. \
		 	user-base/proto/api.proto

main/proto/api.pb.go: proto/api.proto
	protoc 	-I=.. \
			--go_out=.  \
    		--go-grpc_out=. \
		 	./user-base/proto/api.proto


build: main/service.go main/proto/api.pb.go
	go build -o ./build/service ./main

run: build/service
	./build/service

