

all: clean proto

.PHONY: all clean build run proto

clean: 
	rm -rf main/proto build

proto: 
	protoc 	--go_out=.  \
    		--go-grpc_out=. \
		 	proto/auth-api.proto
	protoc 	-I=..\
			--go_out=.  \
    		--go-grpc_out=. \
		 	user-base/proto/user-api.proto





