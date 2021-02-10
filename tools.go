//+build tools

package tools

//https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
import (
	_ "github.com/golang/mock/mockgen"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
