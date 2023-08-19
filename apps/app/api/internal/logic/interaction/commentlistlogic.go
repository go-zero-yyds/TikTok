package interaction

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"github.com/zeromicro/go-zero/core/mr"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListRequest) (resp *types.CommentListResponse, err error) {
	tokenID := int64(0)
	if req.Token != "" {
		tokenID, err = l.svcCtx.JwtAuth.ParseToken(req.Token)
		if err != nil {
			return nil, err
		}
	}
	list, err := l.svcCtx.InteractionRPC.GetCommentList(l.ctx, &interaction.CommentListReq{
		UserId:  tokenID,
		VideoId: req.VideoID,
	})
	if err != nil {
		return nil, err
	}
	if list.CommentList == nil {
		list.CommentList = make([]*interaction.Comment, 0)
	}
	e := apiVars.Success
	commentList, err := mr.MapReduce(func(source chan<- *interaction.Comment) {
		for _, bv := range list.CommentList {
			source <- bv
		}
	}, func(item *interaction.Comment, writer mr.Writer[*types.Comment], cancel func(error)) {
		videoInfo, err := GetCommentInfo(item, tokenID, l.svcCtx, l.ctx)
		if err != nil {
			e = apiVars.SomeDataErr
			if err != apiVars.SomeDataErr {
				return
			}
		}
		writer.Write(videoInfo)
	}, func(pipe <-chan *types.Comment, writer mr.Writer[[]types.Comment], cancel func(error)) {
		var vs []types.Comment
		for item := range pipe {
			v := item
			vs = append(vs, *v)
		}
		writer.Write(vs)
	})

	return &types.CommentListResponse{
		RespStatus:  types.RespStatus(e),
		CommentList: commentList,
	}, nil
}
