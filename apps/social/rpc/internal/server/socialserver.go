// Code generated by goctl. DO NOT EDIT.
// Source: social.proto

package server

import (
	"context"

	"TikTok/apps/social/rpc/internal/logic"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
)

type SocialServer struct {
	svcCtx *svc.ServiceContext
	social.UnimplementedSocialServer
}

func NewSocialServer(svcCtx *svc.ServiceContext) *SocialServer {
	return &SocialServer{
		svcCtx: svcCtx,
	}
}

func (s *SocialServer) IsFollow(ctx context.Context, in *social.IsFollowReq) (*social.IsFollowResp, error) {
	l := logic.NewIsFollowLogic(ctx, s.svcCtx)
	return l.IsFollow(in)
}

func (s *SocialServer) GetFollowCount(ctx context.Context, in *social.FollowCountReq) (*social.FollowCountResp, error) {
	l := logic.NewGetFollowCountLogic(ctx, s.svcCtx)
	return l.GetFollowCount(in)
}

func (s *SocialServer) GetFollowerCount(ctx context.Context, in *social.FollowerCountReq) (*social.FollowerCountResp, error) {
	l := logic.NewGetFollowerCountLogic(ctx, s.svcCtx)
	return l.GetFollowerCount(in)
}

func (s *SocialServer) SendRelationAction(ctx context.Context, in *social.RelationActionReq) (*social.RelationActionResp, error) {
	l := logic.NewSendRelationActionLogic(ctx, s.svcCtx)
	return l.SendRelationAction(in)
}

func (s *SocialServer) GetRelationFollowList(ctx context.Context, in *social.RelationFollowListReq) (*social.RelationFollowListResp, error) {
	l := logic.NewGetRelationFollowListLogic(ctx, s.svcCtx)
	return l.GetRelationFollowList(in)
}

func (s *SocialServer) GetRelationFollowerList(ctx context.Context, in *social.RelationFollowerListReq) (*social.RelationFollowerListResp, error) {
	l := logic.NewGetRelationFollowerListLogic(ctx, s.svcCtx)
	return l.GetRelationFollowerList(in)
}

func (s *SocialServer) GetRelationFriendList(ctx context.Context, in *social.RelationFriendListReq) (*social.RelationFriendListResp, error) {
	l := logic.NewGetRelationFriendListLogic(ctx, s.svcCtx)
	return l.GetRelationFriendList(in)
}

func (s *SocialServer) GetMessages(ctx context.Context, in *social.MessageChatReq) (*social.MessageChatResp, error) {
	l := logic.NewGetMessagesLogic(ctx, s.svcCtx)
	return l.GetMessages(in)
}

func (s *SocialServer) SendMessageAction(ctx context.Context, in *social.MessageActionReq) (*social.MessageActionResp, error) {
	l := logic.NewSendMessageActionLogic(ctx, s.svcCtx)
	return l.SendMessageAction(in)
}
