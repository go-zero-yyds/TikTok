// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: social.proto

package social

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

// SocialClient is the client API for Social service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SocialClient interface {
	IsFollow(ctx context.Context, in *IsFollowReq, opts ...grpc.CallOption) (*IsFollowResp, error)
	GetFollowCount(ctx context.Context, in *FollowCountReq, opts ...grpc.CallOption) (*FollowCountResp, error)
	GetFollowerCount(ctx context.Context, in *FollowerCountReq, opts ...grpc.CallOption) (*FollowerCountResp, error)
	SendRelationAction(ctx context.Context, in *RelationActionReq, opts ...grpc.CallOption) (*RelationActionResp, error)
	GetRelationFollowList(ctx context.Context, in *RelationFollowListReq, opts ...grpc.CallOption) (*RelationFollowListResp, error)
	GetRelationFollowerList(ctx context.Context, in *RelationFollowerListReq, opts ...grpc.CallOption) (*RelationFollowerListResp, error)
	GetRelationFriendList(ctx context.Context, in *RelationFriendListReq, opts ...grpc.CallOption) (*RelationFriendListResp, error)
	GetMessages(ctx context.Context, in *MessageChatReq, opts ...grpc.CallOption) (*MessageChatResp, error)
	SendMessageAction(ctx context.Context, in *MessageActionReq, opts ...grpc.CallOption) (*MessageActionResp, error)
}

type socialClient struct {
	cc grpc.ClientConnInterface
}

func NewSocialClient(cc grpc.ClientConnInterface) SocialClient {
	return &socialClient{cc}
}

func (c *socialClient) IsFollow(ctx context.Context, in *IsFollowReq, opts ...grpc.CallOption) (*IsFollowResp, error) {
	out := new(IsFollowResp)
	err := c.cc.Invoke(ctx, "/social.Social/IsFollow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetFollowCount(ctx context.Context, in *FollowCountReq, opts ...grpc.CallOption) (*FollowCountResp, error) {
	out := new(FollowCountResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetFollowCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetFollowerCount(ctx context.Context, in *FollowerCountReq, opts ...grpc.CallOption) (*FollowerCountResp, error) {
	out := new(FollowerCountResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetFollowerCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) SendRelationAction(ctx context.Context, in *RelationActionReq, opts ...grpc.CallOption) (*RelationActionResp, error) {
	out := new(RelationActionResp)
	err := c.cc.Invoke(ctx, "/social.Social/SendRelationAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetRelationFollowList(ctx context.Context, in *RelationFollowListReq, opts ...grpc.CallOption) (*RelationFollowListResp, error) {
	out := new(RelationFollowListResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetRelationFollowList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetRelationFollowerList(ctx context.Context, in *RelationFollowerListReq, opts ...grpc.CallOption) (*RelationFollowerListResp, error) {
	out := new(RelationFollowerListResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetRelationFollowerList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetRelationFriendList(ctx context.Context, in *RelationFriendListReq, opts ...grpc.CallOption) (*RelationFriendListResp, error) {
	out := new(RelationFriendListResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetRelationFriendList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) GetMessages(ctx context.Context, in *MessageChatReq, opts ...grpc.CallOption) (*MessageChatResp, error) {
	out := new(MessageChatResp)
	err := c.cc.Invoke(ctx, "/social.Social/GetMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialClient) SendMessageAction(ctx context.Context, in *MessageActionReq, opts ...grpc.CallOption) (*MessageActionResp, error) {
	out := new(MessageActionResp)
	err := c.cc.Invoke(ctx, "/social.Social/SendMessageAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SocialServer is the server API for Social service.
// All implementations must embed UnimplementedSocialServer
// for forward compatibility
type SocialServer interface {
	IsFollow(context.Context, *IsFollowReq) (*IsFollowResp, error)
	GetFollowCount(context.Context, *FollowCountReq) (*FollowCountResp, error)
	GetFollowerCount(context.Context, *FollowerCountReq) (*FollowerCountResp, error)
	SendRelationAction(context.Context, *RelationActionReq) (*RelationActionResp, error)
	GetRelationFollowList(context.Context, *RelationFollowListReq) (*RelationFollowListResp, error)
	GetRelationFollowerList(context.Context, *RelationFollowerListReq) (*RelationFollowerListResp, error)
	GetRelationFriendList(context.Context, *RelationFriendListReq) (*RelationFriendListResp, error)
	GetMessages(context.Context, *MessageChatReq) (*MessageChatResp, error)
	SendMessageAction(context.Context, *MessageActionReq) (*MessageActionResp, error)
	mustEmbedUnimplementedSocialServer()
}

// UnimplementedSocialServer must be embedded to have forward compatible implementations.
type UnimplementedSocialServer struct {
}

func (UnimplementedSocialServer) IsFollow(context.Context, *IsFollowReq) (*IsFollowResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFollow not implemented")
}
func (UnimplementedSocialServer) GetFollowCount(context.Context, *FollowCountReq) (*FollowCountResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowCount not implemented")
}
func (UnimplementedSocialServer) GetFollowerCount(context.Context, *FollowerCountReq) (*FollowerCountResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowerCount not implemented")
}
func (UnimplementedSocialServer) SendRelationAction(context.Context, *RelationActionReq) (*RelationActionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRelationAction not implemented")
}
func (UnimplementedSocialServer) GetRelationFollowList(context.Context, *RelationFollowListReq) (*RelationFollowListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationFollowList not implemented")
}
func (UnimplementedSocialServer) GetRelationFollowerList(context.Context, *RelationFollowerListReq) (*RelationFollowerListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationFollowerList not implemented")
}
func (UnimplementedSocialServer) GetRelationFriendList(context.Context, *RelationFriendListReq) (*RelationFriendListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationFriendList not implemented")
}
func (UnimplementedSocialServer) GetMessages(context.Context, *MessageChatReq) (*MessageChatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedSocialServer) SendMessageAction(context.Context, *MessageActionReq) (*MessageActionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessageAction not implemented")
}
func (UnimplementedSocialServer) mustEmbedUnimplementedSocialServer() {}

// UnsafeSocialServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SocialServer will
// result in compilation errors.
type UnsafeSocialServer interface {
	mustEmbedUnimplementedSocialServer()
}

func RegisterSocialServer(s grpc.ServiceRegistrar, srv SocialServer) {
	s.RegisterService(&Social_ServiceDesc, srv)
}

func _Social_IsFollow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsFollowReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).IsFollow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/IsFollow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).IsFollow(ctx, req.(*IsFollowReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetFollowCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowCountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetFollowCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetFollowCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetFollowCount(ctx, req.(*FollowCountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetFollowerCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowerCountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetFollowerCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetFollowerCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetFollowerCount(ctx, req.(*FollowerCountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_SendRelationAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelationActionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).SendRelationAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/SendRelationAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).SendRelationAction(ctx, req.(*RelationActionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetRelationFollowList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelationFollowListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetRelationFollowList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetRelationFollowList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetRelationFollowList(ctx, req.(*RelationFollowListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetRelationFollowerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelationFollowerListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetRelationFollowerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetRelationFollowerList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetRelationFollowerList(ctx, req.(*RelationFollowerListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetRelationFriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelationFriendListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetRelationFriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetRelationFriendList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetRelationFriendList(ctx, req.(*RelationFriendListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_GetMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageChatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).GetMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/GetMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).GetMessages(ctx, req.(*MessageChatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Social_SendMessageAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageActionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialServer).SendMessageAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/social.Social/SendMessageAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialServer).SendMessageAction(ctx, req.(*MessageActionReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Social_ServiceDesc is the grpc.ServiceDesc for Social service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Social_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "social.Social",
	HandlerType: (*SocialServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsFollow",
			Handler:    _Social_IsFollow_Handler,
		},
		{
			MethodName: "GetFollowCount",
			Handler:    _Social_GetFollowCount_Handler,
		},
		{
			MethodName: "GetFollowerCount",
			Handler:    _Social_GetFollowerCount_Handler,
		},
		{
			MethodName: "SendRelationAction",
			Handler:    _Social_SendRelationAction_Handler,
		},
		{
			MethodName: "GetRelationFollowList",
			Handler:    _Social_GetRelationFollowList_Handler,
		},
		{
			MethodName: "GetRelationFollowerList",
			Handler:    _Social_GetRelationFollowerList_Handler,
		},
		{
			MethodName: "GetRelationFriendList",
			Handler:    _Social_GetRelationFriendList_Handler,
		},
		{
			MethodName: "GetMessages",
			Handler:    _Social_GetMessages_Handler,
		},
		{
			MethodName: "SendMessageAction",
			Handler:    _Social_SendMessageAction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "social.proto",
}