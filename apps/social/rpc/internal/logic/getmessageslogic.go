package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"time"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMessages 获取消息列表
func (l *GetMessagesLogic) GetMessages(in *social.MessageChatReq) (messageChatResp *social.MessageChatResp, err error) {
	//查询 social 表中是否有该用户
	exist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回空
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.MessageChatResp{MessageList: nil}, nil
	}

	//根据 preMsgTime, 查询 friend 表中 userId 与 toUserId 的 created_time 前的所有聊天记录
	//preMsgTime转换为时间
	t := time.Unix(in.PreMsgTime, 0)
	formatStr := "2006-01-02 15:04:05"
	newTime := t.Format(formatStr)

	messageList, err := l.svcCtx.CustomDB.QueryMessageByUserIdAndToUserIdInMessage(l.ctx, in.UserId, in.ToUserId, newTime)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.MessageChatResp{MessageList: nil}, nil
	}

	//返回正确的时间（因查出来的时间字符串会有多余字符，所以直接返回传参时间更简单粗暴）
	for _, v := range messageList {
		v.CreateTime = &newTime
	}

	return &social.MessageChatResp{MessageList: messageList}, nil
}
