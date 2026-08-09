

all: clean proto

.PHONY: all clean proto

clean: 
	rm -rf user-base/proto build

proto: userProtoGateway userProto friendProtoGateway friendProto



.PHONY: userProtoGateway
userProtoGateway:
	protoc -I=".." --grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative ../user-base/proto/api.proto
.PHONY: userProto
userProto:
	protoc -I=.. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../user-base/proto/api.proto


.PHONY: friendProtoGateway
friendProtoGateway:
	protoc -I=".." --grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative ../friend-request-base/proto/api.proto
.PHONY: friendProto
friendProto:
	protoc -I=.. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../friend-request-base/proto/api.proto
