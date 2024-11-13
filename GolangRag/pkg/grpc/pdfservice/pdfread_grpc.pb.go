// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: pdfread.proto

package grpcclient

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

const (
	PdfService_ExtractText_FullMethodName   = "/pdfservice.PdfService/ExtractText"
	PdfService_VectorizeText_FullMethodName = "/pdfservice.PdfService/VectorizeText"
)

// PdfServiceClient is the client API for PdfService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PdfServiceClient interface {
	ExtractText(ctx context.Context, in *PdfRequest, opts ...grpc.CallOption) (PdfService_ExtractTextClient, error)
	VectorizeText(ctx context.Context, in *VectorizeRequest, opts ...grpc.CallOption) (*VectorizeResponse, error)
}

type pdfServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPdfServiceClient(cc grpc.ClientConnInterface) PdfServiceClient {
	return &pdfServiceClient{cc}
}

func (c *pdfServiceClient) ExtractText(ctx context.Context, in *PdfRequest, opts ...grpc.CallOption) (PdfService_ExtractTextClient, error) {
	stream, err := c.cc.NewStream(ctx, &PdfService_ServiceDesc.Streams[0], PdfService_ExtractText_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &pdfServiceExtractTextClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PdfService_ExtractTextClient interface {
	Recv() (*PdfResponse, error)
	grpc.ClientStream
}

type pdfServiceExtractTextClient struct {
	grpc.ClientStream
}

func (x *pdfServiceExtractTextClient) Recv() (*PdfResponse, error) {
	m := new(PdfResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pdfServiceClient) VectorizeText(ctx context.Context, in *VectorizeRequest, opts ...grpc.CallOption) (*VectorizeResponse, error) {
	out := new(VectorizeResponse)
	err := c.cc.Invoke(ctx, PdfService_VectorizeText_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PdfServiceServer is the server API for PdfService service.
// All implementations must embed UnimplementedPdfServiceServer
// for forward compatibility
type PdfServiceServer interface {
	ExtractText(*PdfRequest, PdfService_ExtractTextServer) error
	VectorizeText(context.Context, *VectorizeRequest) (*VectorizeResponse, error)
	mustEmbedUnimplementedPdfServiceServer()
}

// UnimplementedPdfServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPdfServiceServer struct {
}

func (UnimplementedPdfServiceServer) ExtractText(*PdfRequest, PdfService_ExtractTextServer) error {
	return status.Errorf(codes.Unimplemented, "method ExtractText not implemented")
}
func (UnimplementedPdfServiceServer) VectorizeText(context.Context, *VectorizeRequest) (*VectorizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VectorizeText not implemented")
}
func (UnimplementedPdfServiceServer) mustEmbedUnimplementedPdfServiceServer() {}

// UnsafePdfServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PdfServiceServer will
// result in compilation errors.
type UnsafePdfServiceServer interface {
	mustEmbedUnimplementedPdfServiceServer()
}

func RegisterPdfServiceServer(s grpc.ServiceRegistrar, srv PdfServiceServer) {
	s.RegisterService(&PdfService_ServiceDesc, srv)
}

func _PdfService_ExtractText_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PdfRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PdfServiceServer).ExtractText(m, &pdfServiceExtractTextServer{stream})
}

type PdfService_ExtractTextServer interface {
	Send(*PdfResponse) error
	grpc.ServerStream
}

type pdfServiceExtractTextServer struct {
	grpc.ServerStream
}

func (x *pdfServiceExtractTextServer) Send(m *PdfResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PdfService_VectorizeText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VectorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PdfServiceServer).VectorizeText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PdfService_VectorizeText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PdfServiceServer).VectorizeText(ctx, req.(*VectorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PdfService_ServiceDesc is the grpc.ServiceDesc for PdfService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PdfService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pdfservice.PdfService",
	HandlerType: (*PdfServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VectorizeText",
			Handler:    _PdfService_VectorizeText_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExtractText",
			Handler:       _PdfService_ExtractText_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pdfread.proto",
}
