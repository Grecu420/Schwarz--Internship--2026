

all: clean build

clean: 
	rm -r main/proto 
	rm -r build

main/proto/api.pb.go: proto/api.proto
	protoc 	--go_out=.  \
    		--go-grpc_out=. \
		 	./proto/api.proto


build: main/service.go main/proto/api.pb.go
	go build -o ./build/service ./main

run: build/service
	./build/service

