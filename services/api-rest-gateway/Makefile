

all: clean proto

.PHONY: all clean proto

clean: 
	rm -rf main/proto build

proto: 
	protoc 	-I=.. \
			--go_out=.  \
    		--go-grpc_out=. \
			--grpc-gateway_out ./ \
		 	api-rest-gateway/proto/gateway-api.proto
	protoc 	-I=.. \
			--go_out=.  \
			--go-grpc_out=. \
		 	user-base/proto/user-api.proto
	protoc 	-I=.. \
			--go_out=.  \
			--go-grpc_out=. \
		 	friend-request-base/proto/friend-request-api.proto
	protoc 	-I=.. \
			--go_out=.  \
			--go-grpc_out=. \
		 	auth-base/proto/auth-api.proto


