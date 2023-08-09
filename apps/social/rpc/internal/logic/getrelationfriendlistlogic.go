package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFriendListLogic {
	return &GetRelationFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRelationFriendList 获取好友列表（只要双方互相关注了自动变成好友）
func (l *GetRelationFriendListLogic) GetRelationFriendList(in *social.RelationFriendListReq) (*social.RelationFriendListResp, error) {
	//查询 social 表中是否有该 user_id
	exist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回空
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	//查询 friend 表中 userId 对应的 toUserId 的 social 数据
	socialsList, err := l.svcCtx.CustomDB.QueryUsersOfFriendListByUserId(l.ctx, in.UserId)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}
	log.Println(socialsList)
	//获取最新消息

	//获取到social
	return &social.RelationFriendListResp{UserList: nil}, nil
}
