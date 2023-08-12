package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageActionLogic {
	return &SendMessageActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendMessageAction 发送消息
func (l *SendMessageActionLogic) SendMessageAction(in *social.MessageActionReq) (*social.MessageActionResp, error) {
	//查询 social 表中是否有这两个用户
	UserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)
	ToUserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.ToUserId)

	//如果不存在则直接返回失败
	if UserIdExist == false || ToUserIdExist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
		return &social.MessageActionResp{IsSucceed: false}, nil
	}

	//查询他们是否为好友，否则不能发送消息
	exist, _ := l.svcCtx.CustomDB.QueryIsExistByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)
	if exist != true || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
		return &social.MessageActionResp{IsSucceed: false}, nil
	}

	//如果不等于1则return
	if in.ActionType != 1 {
		return &social.MessageActionResp{IsSucceed: false}, nil
	}

	//如果为1则插入对应消息记录
	err = l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdAndContentInMessage(l.ctx, in.UserId, in.ToUserId, in.Content)

	return &social.MessageActionResp{IsSucceed: true}, nil
}
