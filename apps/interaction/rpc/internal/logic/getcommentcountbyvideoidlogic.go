package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
//调用dB接口中函数，获取视频评论数量
//只有在底层数据库出现未知错误会返回err
func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *interaction.CommentCountByVideoIdReq) (*interaction.CommentCountByVideoIdResp, error) {
	count , err := l.svcCtx.DBact.CommentCountByVideoId(in.VideoId)
	if err != nil{
		return nil , err;
	}
	return &interaction.CommentCountByVideoIdResp{
		CommentCount: count,
	}, nil
}
