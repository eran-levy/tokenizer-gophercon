// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package tokenizer

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TokenizerClient is the client API for Tokenizer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenizerClient interface {
	GetTokens(ctx context.Context, in *TokenizePayloadRequest, opts ...grpc.CallOption) (*TokenizePayloadReresponse, error)
}

type tokenizerClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenizerClient(cc grpc.ClientConnInterface) TokenizerClient {
	return &tokenizerClient{cc}
}

func (c *tokenizerClient) GetTokens(ctx context.Context, in *TokenizePayloadRequest, opts ...grpc.CallOption) (*TokenizePayloadReresponse, error) {
	out := new(TokenizePayloadReresponse)
	err := c.cc.Invoke(ctx, "/tokenizer.Tokenizer/GetTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenizerServer is the server API for Tokenizer service.
// All implementations must embed UnimplementedTokenizerServer
// for forward compatibility
type TokenizerServer interface {
	GetTokens(context.Context, *TokenizePayloadRequest) (*TokenizePayloadReresponse, error)
	mustEmbedUnimplementedTokenizerServer()
}

// UnimplementedTokenizerServer must be embedded to have forward compatible implementations.
type UnimplementedTokenizerServer struct {
}

func (UnimplementedTokenizerServer) GetTokens(context.Context, *TokenizePayloadRequest) (*TokenizePayloadReresponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokens not implemented")
}
func (UnimplementedTokenizerServer) mustEmbedUnimplementedTokenizerServer() {}

// UnsafeTokenizerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenizerServer will
// result in compilation errors.
type UnsafeTokenizerServer interface {
	mustEmbedUnimplementedTokenizerServer()
}

func RegisterTokenizerServer(s grpc.ServiceRegistrar, srv TokenizerServer) {
	s.RegisterService(&Tokenizer_ServiceDesc, srv)
}

func _Tokenizer_GetTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenizePayloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenizerServer).GetTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenizer.Tokenizer/GetTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenizerServer).GetTokens(ctx, req.(*TokenizePayloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Tokenizer_ServiceDesc is the grpc.ServiceDesc for Tokenizer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tokenizer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tokenizer.Tokenizer",
	HandlerType: (*TokenizerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTokens",
			Handler:    _Tokenizer_GetTokens_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokenizer.proto",
}
