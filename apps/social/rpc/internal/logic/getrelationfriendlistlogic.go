package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
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
	//查询 social 表中是否有该用户
	exist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回空
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	//查询 friend 表中 userId/toUserId 对应的 userId/toUserId
	friendIdList, err := l.svcCtx.CustomDB.QueryFriendIdListByUserIdInFriend(l.ctx, in.UserId)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	//获取到每条消息
	messageList, err := l.svcCtx.CustomDB.QueryMessageByUserIdAndUserListInMessage(l.ctx, in.UserId, friendIdList)
	if err != nil {
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	log.Println(messageList)
	//拼接friendUserList
	FriendUserList := make([]*social.FriendUser, len(messageList))
	for i, v := range messageList {
		FriendUserList[i] = &social.FriendUser{}
		//拼接userId和消息类型
		if v.FromUserId != in.UserId {
			FriendUserList[i].UserId = v.FromUserId
			FriendUserList[i].MsgType = 0
		} else {
			FriendUserList[i].UserId = v.ToUserId
			FriendUserList[i].MsgType = 1
		}

		//拼接content
		FriendUserList[i].Message = v.Content

	}

	return &social.RelationFriendListResp{UserList: FriendUserList}, nil
}
