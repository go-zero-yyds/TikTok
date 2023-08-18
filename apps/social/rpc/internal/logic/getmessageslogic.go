package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"strconv"
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

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId, in.ToUserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
	}

	//根据 preMsgTime, 查询 friend 表中 userId 与 toUserId 的 created_time 前的所有聊天记录
	//preMsgTime转换为时间
	t := time.Unix(0, in.PreMsgTime*int64(time.Millisecond))
	formatStr := "2006-01-02 15:04:05"
	newTime := t.Format(formatStr)

	messageList, err := l.svcCtx.CustomDB.QueryMessageByUserIdAndToUserIdInMessage(l.ctx, in.UserId, in.ToUserId, newTime)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.MessageChatResp{MessageList: nil}, nil
	}

	//返回时间戳
	for _, v := range messageList {
		//将其转换为ms级别的时间戳
		layout := "2006-01-02T15:04:05-07:00"
		tempTimeStamp, _ := time.Parse(layout, *v.CreateTime)
		timestampStr := strconv.FormatInt(tempTimeStamp.UnixNano()/int64(time.Millisecond), 10)
		v.CreateTime = &timestampStr
	}

	return &social.MessageChatResp{MessageList: messageList}, nil
}
