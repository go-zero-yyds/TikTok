// Code generated by goctl. DO NOT EDIT.
// Source: interaction.proto

package interactionclient

import (
	"context"

	"rpc/apps/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Comment                    = interaction.Comment
	CommentActionReq           = interaction.CommentActionReq
	CommentActionResp          = interaction.CommentActionResp
	CommentCountByVideoIdReq   = interaction.CommentCountByVideoIdReq
	CommentCountByVideoIdResp  = interaction.CommentCountByVideoIdResp
	CommentListReq             = interaction.CommentListReq
	CommentListResp            = interaction.CommentListResp
	FavoriteActionReq          = interaction.FavoriteActionReq
	FavoriteActionResp         = interaction.FavoriteActionResp
	FavoriteCountByUserIdReq   = interaction.FavoriteCountByUserIdReq
	FavoriteCountByUserIdResp  = interaction.FavoriteCountByUserIdResp
	FavoriteCountByVideoIdReq  = interaction.FavoriteCountByVideoIdReq
	FavoriteCountByVideoIdResp = interaction.FavoriteCountByVideoIdResp
	FavoriteListReq            = interaction.FavoriteListReq
	FavoriteListResp           = interaction.FavoriteListResp
	IsFavoriteReq              = interaction.IsFavoriteReq
	IsFavoriteResp             = interaction.IsFavoriteResp

	Interaction interface {
		IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error)
		GetFavoriteCountByUserId(ctx context.Context, in *FavoriteCountByUserIdReq, opts ...grpc.CallOption) (*FavoriteCountByUserIdResp, error)
		GetFavoriteCountByVideoId(ctx context.Context, in *FavoriteCountByVideoIdReq, opts ...grpc.CallOption) (*FavoriteCountByVideoIdResp, error)
		GetCommentCountByVideoId(ctx context.Context, in *CommentCountByVideoIdReq, opts ...grpc.CallOption) (*CommentCountByVideoIdResp, error)
		SendFavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error)
		GetFavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error)
		SendCommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error)
		GetCommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error)
	}

	defaultInteraction struct {
		cli zrpc.Client
	}
)

func NewInteraction(cli zrpc.Client) Interaction {
	return &defaultInteraction{
		cli: cli,
	}
}

func (m *defaultInteraction) IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.IsFavorite(ctx, in, opts...)
}

func (m *defaultInteraction) GetFavoriteCountByUserId(ctx context.Context, in *FavoriteCountByUserIdReq, opts ...grpc.CallOption) (*FavoriteCountByUserIdResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.GetFavoriteCountByUserId(ctx, in, opts...)
}

func (m *defaultInteraction) GetFavoriteCountByVideoId(ctx context.Context, in *FavoriteCountByVideoIdReq, opts ...grpc.CallOption) (*FavoriteCountByVideoIdResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.GetFavoriteCountByVideoId(ctx, in, opts...)
}

func (m *defaultInteraction) GetCommentCountByVideoId(ctx context.Context, in *CommentCountByVideoIdReq, opts ...grpc.CallOption) (*CommentCountByVideoIdResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.GetCommentCountByVideoId(ctx, in, opts...)
}

func (m *defaultInteraction) SendFavoriteAction(ctx context.Context, in *FavoriteActionReq, opts ...grpc.CallOption) (*FavoriteActionResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.SendFavoriteAction(ctx, in, opts...)
}

func (m *defaultInteraction) GetFavoriteList(ctx context.Context, in *FavoriteListReq, opts ...grpc.CallOption) (*FavoriteListResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.GetFavoriteList(ctx, in, opts...)
}

func (m *defaultInteraction) SendCommentAction(ctx context.Context, in *CommentActionReq, opts ...grpc.CallOption) (*CommentActionResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.SendCommentAction(ctx, in, opts...)
}

func (m *defaultInteraction) GetCommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error) {
	client := interaction.NewInteractionClient(m.cli.Conn())
	return client.GetCommentList(ctx, in, opts...)
}
