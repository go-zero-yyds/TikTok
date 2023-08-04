package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
// 查看视频的所有评论，按发布时间倒序
func (l *GetCommentListLogic) GetCommentList(in *interaction.CommentListReq) (*interaction.CommentListResp, error) {
	commentlist , err := l.svcCtx.DBact.CommentList(in.UserId , in.VideoId)
	if err != nil{
		return nil , err
	}
	ret := make([]*interaction.Comment , 0 , len(commentlist))
	for _, v := range commentlist{
		ret = append(ret, &interaction.Comment{
			Id: v.CommentId,
			UserId: v.UserId,
			Content: v.Content,
			CreateDate: v.CreateDate.Format("01-02"),
		})
	}
	return &interaction.CommentListResp{
		CommentList: ret,
	}, nil
}
