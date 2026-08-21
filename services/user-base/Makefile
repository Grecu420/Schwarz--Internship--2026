FRONTEND_DIR := ../../frontend/src/generated

all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build $(FRONTEND_DIR)

proto: 
	protoc  --go_out=.  \
	        --go-grpc_out=. \
	        proto/user-api.proto
	
	mkdir -p $(FRONTEND_DIR)
	protoc  --plugin=protoc-gen-ts_proto=../../frontend/node_modules/.bin/protoc-gen-ts_proto \
	        --ts_proto_out=outputServices=grpc-web,env=browser,outputClientImpl=true:$(FRONTEND_DIR) \
	        proto/user-api.proto

main/proto/user-api.pb.go: proto/user-api.proto
	protoc  -I=.. \
	        --go_out=.  \
	        --go-grpc_out=. \
	        ./user-base/proto/user-api.proto

build: main/service.go main/proto/user-api.pb.go
	go build -o ./build/service ./main

run: build/service
	./build/service