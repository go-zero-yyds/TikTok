// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: video.proto

package video

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

// VideoClient is the client API for Video service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoClient interface {
	GetWorkCountByUserId(ctx context.Context, in *WorkCountByUserIdReq, opts ...grpc.CallOption) (*WorkCountByUserIdResp, error)
	GetFeed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error)
	SendPublishAction(ctx context.Context, in *PublishActionReq, opts ...grpc.CallOption) (*PublishActionResp, error)
	GetPublishList(ctx context.Context, in *PublishListReq, opts ...grpc.CallOption) (*PublishListResp, error)
	Detail(ctx context.Context, in *BasicVideoInfoReq, opts ...grpc.CallOption) (*BasicVideoInfoResp, error)
	// 视频事件通知,如:点赞、评论、分享
	NotifyHotVideo(ctx context.Context, in *NotifyHotVideoReq, opts ...grpc.CallOption) (*NotifyHotVideoResp, error)
	// 校验热门视频
	CheckHotVideo(ctx context.Context, in *CheckHotVideoReq, opts ...grpc.CallOption) (*CheckHotVideoResp, error)
	// pusher推送视频id到consumer, consumer校验视频是否热门视频再做别的业务
	VideoPusher(ctx context.Context, in *VideoPusherReq, opts ...grpc.CallOption) (*VideoPusherResp, error)
}

type videoClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoClient(cc grpc.ClientConnInterface) VideoClient {
	return &videoClient{cc}
}

func (c *videoClient) GetWorkCountByUserId(ctx context.Context, in *WorkCountByUserIdReq, opts ...grpc.CallOption) (*WorkCountByUserIdResp, error) {
	out := new(WorkCountByUserIdResp)
	err := c.cc.Invoke(ctx, "/video.Video/GetWorkCountByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetFeed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedResp, error) {
	out := new(FeedResp)
	err := c.cc.Invoke(ctx, "/video.Video/GetFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) SendPublishAction(ctx context.Context, in *PublishActionReq, opts ...grpc.CallOption) (*PublishActionResp, error) {
	out := new(PublishActionResp)
	err := c.cc.Invoke(ctx, "/video.Video/SendPublishAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetPublishList(ctx context.Context, in *PublishListReq, opts ...grpc.CallOption) (*PublishListResp, error) {
	out := new(PublishListResp)
	err := c.cc.Invoke(ctx, "/video.Video/GetPublishList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) Detail(ctx context.Context, in *BasicVideoInfoReq, opts ...grpc.CallOption) (*BasicVideoInfoResp, error) {
	out := new(BasicVideoInfoResp)
	err := c.cc.Invoke(ctx, "/video.Video/Detail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) NotifyHotVideo(ctx context.Context, in *NotifyHotVideoReq, opts ...grpc.CallOption) (*NotifyHotVideoResp, error) {
	out := new(NotifyHotVideoResp)
	err := c.cc.Invoke(ctx, "/video.Video/NotifyHotVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) CheckHotVideo(ctx context.Context, in *CheckHotVideoReq, opts ...grpc.CallOption) (*CheckHotVideoResp, error) {
	out := new(CheckHotVideoResp)
	err := c.cc.Invoke(ctx, "/video.Video/CheckHotVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) VideoPusher(ctx context.Context, in *VideoPusherReq, opts ...grpc.CallOption) (*VideoPusherResp, error) {
	out := new(VideoPusherResp)
	err := c.cc.Invoke(ctx, "/video.Video/VideoPusher", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoServer is the server API for Video service.
// All implementations must embed UnimplementedVideoServer
// for forward compatibility
type VideoServer interface {
	GetWorkCountByUserId(context.Context, *WorkCountByUserIdReq) (*WorkCountByUserIdResp, error)
	GetFeed(context.Context, *FeedReq) (*FeedResp, error)
	SendPublishAction(context.Context, *PublishActionReq) (*PublishActionResp, error)
	GetPublishList(context.Context, *PublishListReq) (*PublishListResp, error)
	Detail(context.Context, *BasicVideoInfoReq) (*BasicVideoInfoResp, error)
	// 视频事件通知,如:点赞、评论、分享
	NotifyHotVideo(context.Context, *NotifyHotVideoReq) (*NotifyHotVideoResp, error)
	// 校验热门视频
	CheckHotVideo(context.Context, *CheckHotVideoReq) (*CheckHotVideoResp, error)
	// pusher推送视频id到consumer, consumer校验视频是否热门视频再做别的业务
	VideoPusher(context.Context, *VideoPusherReq) (*VideoPusherResp, error)
	mustEmbedUnimplementedVideoServer()
}

// UnimplementedVideoServer must be embedded to have forward compatible implementations.
type UnimplementedVideoServer struct {
}

func (UnimplementedVideoServer) GetWorkCountByUserId(context.Context, *WorkCountByUserIdReq) (*WorkCountByUserIdResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkCountByUserId not implemented")
}
func (UnimplementedVideoServer) GetFeed(context.Context, *FeedReq) (*FeedResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeed not implemented")
}
func (UnimplementedVideoServer) SendPublishAction(context.Context, *PublishActionReq) (*PublishActionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPublishAction not implemented")
}
func (UnimplementedVideoServer) GetPublishList(context.Context, *PublishListReq) (*PublishListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublishList not implemented")
}
func (UnimplementedVideoServer) Detail(context.Context, *BasicVideoInfoReq) (*BasicVideoInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detail not implemented")
}
func (UnimplementedVideoServer) NotifyHotVideo(context.Context, *NotifyHotVideoReq) (*NotifyHotVideoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyHotVideo not implemented")
}
func (UnimplementedVideoServer) CheckHotVideo(context.Context, *CheckHotVideoReq) (*CheckHotVideoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHotVideo not implemented")
}
func (UnimplementedVideoServer) VideoPusher(context.Context, *VideoPusherReq) (*VideoPusherResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoPusher not implemented")
}
func (UnimplementedVideoServer) mustEmbedUnimplementedVideoServer() {}

// UnsafeVideoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoServer will
// result in compilation errors.
type UnsafeVideoServer interface {
	mustEmbedUnimplementedVideoServer()
}

func RegisterVideoServer(s grpc.ServiceRegistrar, srv VideoServer) {
	s.RegisterService(&Video_ServiceDesc, srv)
}

func _Video_GetWorkCountByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkCountByUserIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetWorkCountByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/GetWorkCountByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetWorkCountByUserId(ctx, req.(*WorkCountByUserIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/GetFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetFeed(ctx, req.(*FeedReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_SendPublishAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishActionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).SendPublishAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/SendPublishAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).SendPublishAction(ctx, req.(*PublishActionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetPublishList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetPublishList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/GetPublishList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetPublishList(ctx, req.(*PublishListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BasicVideoInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/Detail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).Detail(ctx, req.(*BasicVideoInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_NotifyHotVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyHotVideoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).NotifyHotVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/NotifyHotVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).NotifyHotVideo(ctx, req.(*NotifyHotVideoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_CheckHotVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckHotVideoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).CheckHotVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/CheckHotVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).CheckHotVideo(ctx, req.(*CheckHotVideoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_VideoPusher_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoPusherReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).VideoPusher(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.Video/VideoPusher",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).VideoPusher(ctx, req.(*VideoPusherReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Video_ServiceDesc is the grpc.ServiceDesc for Video service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Video_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "video.Video",
	HandlerType: (*VideoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorkCountByUserId",
			Handler:    _Video_GetWorkCountByUserId_Handler,
		},
		{
			MethodName: "GetFeed",
			Handler:    _Video_GetFeed_Handler,
		},
		{
			MethodName: "SendPublishAction",
			Handler:    _Video_SendPublishAction_Handler,
		},
		{
			MethodName: "GetPublishList",
			Handler:    _Video_GetPublishList_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _Video_Detail_Handler,
		},
		{
			MethodName: "NotifyHotVideo",
			Handler:    _Video_NotifyHotVideo_Handler,
		},
		{
			MethodName: "CheckHotVideo",
			Handler:    _Video_CheckHotVideo_Handler,
		},
		{
			MethodName: "VideoPusher",
			Handler:    _Video_VideoPusher_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "video.proto",
}
