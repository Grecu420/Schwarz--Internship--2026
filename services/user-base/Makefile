

all: clean build

clean: 
	rm -r main/proto 

main/proto/api.pb.go: proto/api.proto
	protoc --go_out=. ./proto/api.proto

build: main/service.go main/proto/api.pb.go
	go build -o service ./main

run: service
	./service

