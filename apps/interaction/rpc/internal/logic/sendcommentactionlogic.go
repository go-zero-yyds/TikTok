package logic

import (
	"TikTok/apps/interaction/rpc/model"
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
	if in.ActionType > 2 {
		return nil, nil
	}
	var c model.Comment
	c.UserId = in.UserId
	c.VideoId = in.VideoId
	if in.ActionType == 1 {
		c.CommentId = l.svcCtx.Snowflake.Generate().Int64() // 雪花算法生成id
		c.Content = *in.CommentText
		if in.IPAddr != nil {
			c.IpAddress = *in.IPAddr
		} else {
			c.IpAddress = ""
		}
		if in.IPAttr != nil {
			c.Location = *in.IPAttr
		} else {
			c.Location = ""
		}
	} else {
		c.CommentId = *in.CommentId
	}
	err := l.svcCtx.DBAction.CommentAction(l.ctx, &c, in.ActionType)
	if err != nil {
		return nil, err
	}

	if in.ActionType == 1 {
		return &interaction.CommentActionResp{
			Comment: &interaction.Comment{
				Id:         c.CommentId,
				UserId:     c.UserId,
				Content:    c.Content,
				CreateDate: fmt.Sprintf("%v", c.CreateDate.Unix()),
				IpAddress:  c.IpAddress,
				Location:   c.Location,
			},
		}, nil
	}
	return &interaction.CommentActionResp{
		Comment: nil,
	}, nil
}
