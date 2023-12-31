package social

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/social/rpc/model"
	"TikTok/apps/social/rpc/social"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {

	// 参数检查
	if len(req.Content) > 15000 { //内容是否为符合规范
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apivars.ErrTextRuleError),
		}, nil
	}

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)

	_, err = l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
		Content:    req.Content,
	})
	if errors.Is(err, model.ErrNotFriend) {
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apivars.ErrNotFriend),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &types.MessageActionResponse{RespStatus: types.RespStatus(apivars.Success)}, nil
}
