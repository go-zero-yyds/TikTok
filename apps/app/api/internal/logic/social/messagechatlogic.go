package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/social/rpc/social"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageChatLogic {
	return &MessageChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageChatLogic) MessageChat(req *types.MessageChatRequest) (resp *types.MessageChatResponse, err error) {
	// todo: add your logic here and delete this line

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	messages, err := l.svcCtx.SocialRPC.GetMessages(l.ctx, &social.MessageChatReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		PreMsgTime: req.PreMsgTime,
	})
	if err != nil {
		return nil, err
	}
	res := make([]types.Message, len(messages.MessageList))
	for i, message := range messages.MessageList {
		res[i] = types.Message{
			ID:         message.Id,
			ToUserID:   message.ToUserId,
			FromUserID: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.GetCreateTime(),
		}
	}
	return &types.MessageChatResponse{
		RespStatus:  types.RespStatus(apiVars.Success),
		MessageList: res,
	}, nil
}
