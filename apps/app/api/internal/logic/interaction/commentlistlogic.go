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
	size := len(list.CommentList)
	commentList, err := mr.MapReduce(func(source chan<- IdxComment) {
		for i, bv := range list.CommentList {
			source <- IdxComment{
				Comment: bv,
				idx:     i,
			}
		}
	}, func(item IdxComment, writer mr.Writer[IdxApiComment], cancel func(error)) {
		videoInfo, err := GetCommentInfo(item.Comment, tokenID, l.svcCtx, l.ctx)
		if err != nil {
			e = apiVars.SomeDataErr
			if err != apiVars.SomeDataErr {
				return
			}
		}
		writer.Write(IdxApiComment{
			Comment: videoInfo,
			idx:     item.idx,
		})
	}, func(pipe <-chan IdxApiComment, writer mr.Writer[[]types.Comment], cancel func(error)) {
		vs := make([]types.Comment, size)
		for item := range pipe {
			v := item
			vs[v.idx] = *v.Comment
		}
		writer.Write(vs)
	})

	return &types.CommentListResponse{
		RespStatus:  types.RespStatus(e),
		CommentList: commentList,
	}, nil
}

type IdxComment struct {
	*interaction.Comment
	idx int
}
type IdxApiComment struct {
	*types.Comment
	idx int
}
