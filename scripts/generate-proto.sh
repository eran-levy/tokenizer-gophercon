#!/bin/sh
# you may install the protoc compiler and right-after that install the go plugin for code generation - https://grpc.io/docs/protoc-installation/
# brew install protobuf // apt install -y protobuf-compiler
# https://developers.google.com/protocol-buffers/docs/gotutorial
# install protoc - https://developers.google.com/protocol-buffers/docs/reference/go-generated
# go install google.golang.org/protobuf/cmd/protoc-gen-go
# run protoc in ../pkg/proto/tokenizer
# possible to use --go-grpc_opt=requireUnimplementedServers=false
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative tokenizer.proto
