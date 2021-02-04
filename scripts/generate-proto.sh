#!/bin/sh
# run protoc in ../pkg/proto/tokenizer
# possible to use --go-grpc_opt=requireUnimplementedServers=false
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative tokenizer.proto
