#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$( cd "${SCRIPT_DIR}/../.." && pwd )"
THIRD_PARTY_PATH="${SCRIPT_DIR}/third_party"

cd "${REPO_ROOT}"

echo "Generating ShardingTable proto..."

mkdir -p "${THIRD_PARTY_PATH}"
cd "${THIRD_PARTY_PATH}"

if [ ! -d "xds" ]; then
    echo "Cloning xds..."
    git clone https://github.com/cncf/xds.git
fi

if [ ! -d "udpa" ]; then
    echo "Cloning udpa..."
    git clone https://github.com/cncf/udpa.git
fi

if [ ! -d "protoc-gen-validate" ]; then
    echo "Cloning protoc-gen-validate..."
    git clone https://github.com/envoyproxy/protoc-gen-validate.git
fi

if [ ! -d "googleapis" ]; then
    echo "Cloning googleapis..."
    git clone https://github.com/googleapis/googleapis.git
fi

cd "${REPO_ROOT}"

rm -f examples/shardingtable/resource/sharding_table.pb.go
rm -f examples/shardingtable/resource/sds_grpc.pb.go
rm -f examples/shardingtable/resource/sds.pb.go

ENVOY_PROTO_PATH="/Users/zhiyan.foo/src/envoy/api"

protoc \
  -I. \
  -I"${ENVOY_PROTO_PATH}" \
  -I"${THIRD_PARTY_PATH}/xds" \
  -I"${THIRD_PARTY_PATH}/udpa" \
  -I"${THIRD_PARTY_PATH}/protoc-gen-validate" \
  -I"${THIRD_PARTY_PATH}/googleapis" \
  --go_out=. \
  --go_opt=module=github.com/envoyproxy/go-control-plane \
  examples/shardingtable/proto/sharding_table.proto

protoc \
  -I. \
  -I"${ENVOY_PROTO_PATH}" \
  -I"${THIRD_PARTY_PATH}/xds" \
  -I"${THIRD_PARTY_PATH}/udpa" \
  -I"${THIRD_PARTY_PATH}/protoc-gen-validate" \
  -I"${THIRD_PARTY_PATH}/googleapis" \
  --go_out=. \
  --go_opt=module=github.com/envoyproxy/go-control-plane \
  --go-grpc_out=. \
  --go-grpc_opt=module=github.com/envoyproxy/go-control-plane \
  examples/shardingtable/proto/sds.proto

echo "Done!"
