package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/social/rpc/social"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

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
	// todo: add your logic here and delete this line
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &types.MessageActionResponse{RespStatus: types.RespStatus(apiVars.Success)}, nil
}
