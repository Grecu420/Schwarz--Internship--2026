#!/bin/zsh

OUT_DIR="./src/generated/proto"
PLUGIN="./node_modules/.bin/protoc-gen-ts_proto"
TS_OUT_OPTS="outputServices=grpc-web,env=browser,outputClientImpl=true:${OUT_DIR}"

# List of target proto files
PROTO_FILES=(
  "../services/user-base/proto/user-api.proto"
  "../services/auth-base/proto/auth-api.proto"
  "../services/property-base/proto/property-api.proto"
  "../services/property-list-base/proto/property-list-api.proto"
  "../services/message-base/proto/message-api.proto"
  "../services/conversation-base/proto/conversation-api.proto"
  "../services/reservation-base/proto/reservation-api.proto"
)

# Recreate target output directory
rm -rf src/generated
mkdir -p "$OUT_DIR"

# Generate code for each file
for FILE in "${PROTO_FILES[@]}"; do
  INCLUDE_DIR=$(dirname "$FILE")
  
  protoc \
    --plugin="$PLUGIN" \
    --ts_proto_out="$TS_OUT_OPTS" \
    -I="$INCLUDE_DIR" \
    -I="../services" \
    "$FILE"
done