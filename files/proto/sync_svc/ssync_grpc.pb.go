//
// Copyright (c) 2022.
// Created by Andy Pangaribuan. All Rights Reserved.
//
// This product is protected by copyright and distributed under
// licenses restricting copying, distribution and decompilation.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: files/proto/ssync.proto

package sync_svc

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
	SyncService_KeyLock_FullMethodName = "/sync_svc.SyncService/KeyLock"
)

// SyncServiceClient is the client API for SyncService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SyncServiceClient interface {
	KeyLock(ctx context.Context, opts ...grpc.CallOption) (SyncService_KeyLockClient, error)
}

type syncServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSyncServiceClient(cc grpc.ClientConnInterface) SyncServiceClient {
	return &syncServiceClient{cc}
}

func (c *syncServiceClient) KeyLock(ctx context.Context, opts ...grpc.CallOption) (SyncService_KeyLockClient, error) {
	stream, err := c.cc.NewStream(ctx, &SyncService_ServiceDesc.Streams[0], SyncService_KeyLock_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &syncServiceKeyLockClient{stream}
	return x, nil
}

type SyncService_KeyLockClient interface {
	Send(*KeyLockRequest) error
	Recv() (*KeyLockResponse, error)
	grpc.ClientStream
}

type syncServiceKeyLockClient struct {
	grpc.ClientStream
}

func (x *syncServiceKeyLockClient) Send(m *KeyLockRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *syncServiceKeyLockClient) Recv() (*KeyLockResponse, error) {
	m := new(KeyLockResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SyncServiceServer is the server API for SyncService service.
// All implementations must embed UnimplementedSyncServiceServer
// for forward compatibility
type SyncServiceServer interface {
	KeyLock(SyncService_KeyLockServer) error
	mustEmbedUnimplementedSyncServiceServer()
}

// UnimplementedSyncServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSyncServiceServer struct {
}

func (UnimplementedSyncServiceServer) KeyLock(SyncService_KeyLockServer) error {
	return status.Errorf(codes.Unimplemented, "method KeyLock not implemented")
}
func (UnimplementedSyncServiceServer) mustEmbedUnimplementedSyncServiceServer() {}

// UnsafeSyncServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SyncServiceServer will
// result in compilation errors.
type UnsafeSyncServiceServer interface {
	mustEmbedUnimplementedSyncServiceServer()
}

func RegisterSyncServiceServer(s grpc.ServiceRegistrar, srv SyncServiceServer) {
	s.RegisterService(&SyncService_ServiceDesc, srv)
}

func _SyncService_KeyLock_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SyncServiceServer).KeyLock(&syncServiceKeyLockServer{stream})
}

type SyncService_KeyLockServer interface {
	Send(*KeyLockResponse) error
	Recv() (*KeyLockRequest, error)
	grpc.ServerStream
}

type syncServiceKeyLockServer struct {
	grpc.ServerStream
}

func (x *syncServiceKeyLockServer) Send(m *KeyLockResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *syncServiceKeyLockServer) Recv() (*KeyLockRequest, error) {
	m := new(KeyLockRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SyncService_ServiceDesc is the grpc.ServiceDesc for SyncService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SyncService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sync_svc.SyncService",
	HandlerType: (*SyncServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "KeyLock",
			Handler:       _SyncService_KeyLock_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "files/proto/ssync.proto",
}
