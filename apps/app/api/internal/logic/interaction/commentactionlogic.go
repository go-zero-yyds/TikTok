package interaction

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/interactionclient"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionRequest) (resp *types.CommentActionResponse, err error) {
	// todo: add your logic here and delete this line
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	rpcReq := &interaction.CommentActionReq{
		UserId:     tokenID,
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
	}
	if req.CommentID == 0 {
		rpcReq.CommentText = &req.CommentText
	} else {
		rpcReq.CommentId = &req.CommentID
	}
	sendCommentAction, err := l.svcCtx.InteractionRPC.SendCommentAction(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	apiResp := apiVars.Success

	var commentInfo *types.Comment
	if req.CommentID == 0 {
		commentInfo, err = GetCommentInfo(sendCommentAction.Comment, tokenID, l.svcCtx, l.ctx)

		if err == apiVars.SomeDataErr {
			return &types.CommentActionResponse{
				RespStatus: types.RespStatus(apiResp),
				Comment:    *commentInfo,
			}, nil
		}

		if err != nil && err != apiVars.SomeDataErr {
			return nil, err
		}

		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiResp),
			Comment:    *commentInfo,
		}, nil
	}
	return &types.CommentActionResponse{
		RespStatus: types.RespStatus(apiResp),
	}, nil

}

func GetCommentInfo(comment *interactionclient.Comment, tokenID int64, svcCtx *svc.ServiceContext, ctx context.Context) (*types.Comment, error) {
	res := &types.Comment{}
	res.CreateDate = comment.CreateDate
	res.ID = comment.Id
	res.Content = comment.Content
	userInfo, err := user.TryGetUserInfo(tokenID, comment.UserId, svcCtx, ctx)
	if err != nil && err != apiVars.SomeDataErr {
		return nil, err
	}
	res.User = *userInfo
	if err == apiVars.SomeDataErr {
		return res, apiVars.SomeDataErr
	}

	return res, nil
}
