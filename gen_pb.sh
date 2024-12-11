#!/usr/bin/env bash

PROTO_FILE=$1

if [ -z "$PROTO_FILE" ]; then
  echo "缺少proto路径参数"
  exit 1
fi

protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        "$PROTO_FILE"
