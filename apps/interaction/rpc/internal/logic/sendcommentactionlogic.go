package logic

import (
	"context"
	"fmt"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCommentActionLogic {
	return &SendCommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendCommentAction 调用dB接口中函数 执行评论/取消操作
// 成功返回comment结构体，（评论成功 赋值，取消成功 nil）
func (l *SendCommentActionLogic) SendCommentAction(in *interaction.CommentActionReq) (*interaction.CommentActionResp, error) {
	//植入雪花id
	if in.ActionType == 1 {
		in.CommentId = new(int64)
		*in.CommentId = l.svcCtx.Snowflake.Generate().Int64()
	}
	if in.ActionType > 2 {
		return nil, nil
	}
	resp, err := l.svcCtx.DBAction.CommentAction(l.ctx, in.UserId, in.VideoId, in.ActionType, in.CommentText, in.CommentId)
	if err != nil {
		return nil, err
	}

	if in.ActionType == 1 {
		return &interaction.CommentActionResp{
			Comment: &interaction.Comment{
				Id:         resp.CommentId,
				UserId:     resp.UserId,
				Content:    resp.Content,
				CreateDate: fmt.Sprintf("%v", resp.CreateDate.Unix()),
			},
		}, nil
	}
	return &interaction.CommentActionResp{
		Comment: nil,
	}, nil
}
