#!/bin/sh

#* variables
PROTO_PATH=/proto
IDL_PATH=/pb


#* clean
# rm -rf ${IDL_PATH}/*


protoc \
  ${PROTO_PATH}/*.proto \
  -I=/usr/local/include \
  --proto_path=${PROTO_PATH} \
  --go_out=:${IDL_PATH} \
  --go-grpc_out=:${IDL_PATH} \
  --grpc-gateway_out=:${IDL_PATH} \
  --http-client_out=:${IDL_PATH} 

