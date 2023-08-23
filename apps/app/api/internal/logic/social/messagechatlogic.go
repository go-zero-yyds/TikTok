package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/social/rpc/social"
	"context"
	"regexp"
	"strconv"

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

	// 参数检查
	matched, err := regexp.MatchString("^\\d{19}$", strconv.FormatInt(req.ToUserID, 10)) //是否为19位纯数字
	if strconv.FormatInt(req.ToUserID, 10) == "" || matched == false {
		return &types.MessageChatResponse{
			RespStatus: types.RespStatus(apiVars.UserIdRuleError),
		}, nil
	}

	timestamp, err := strconv.ParseInt(strconv.FormatInt(req.PreMsgTime, 10), 10, 64) //检查时间戳格式是否正确
	timestampStr := strconv.FormatInt(timestamp, 10)
	if (req.PreMsgTime != 0) && (len(timestampStr) != 13) == true {
		return &types.MessageChatResponse{
			RespStatus: types.RespStatus(apiVars.TimestampRuleError),
		}, nil
	}

	if req.Token == "" {
		return &types.MessageChatResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

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
			CreateTime: message.CreateTime,
		}
	}
	return &types.MessageChatResponse{
		RespStatus:  types.RespStatus(apiVars.Success),
		MessageList: res,
	}, nil
}
