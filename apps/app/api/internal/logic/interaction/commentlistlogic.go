package interaction

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/mr"

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

	tokenID := int64(-1)
	if req.Token != "" {
		tokenID, err = l.svcCtx.JwtAuth.ParseToken(req.Token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = nil
			} else {
				return nil, err
			}

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
	e := apivars.Success
	size := len(list.CommentList)
	commentList, err := mr.MapReduce(func(source chan<- idxComment) {
		for i, bv := range list.CommentList {
			source <- idxComment{
				Comment: bv,
				idx:     i,
			}
		}
	}, func(item idxComment, writer mr.Writer[idxApiComment], cancel func(error)) {
		videoInfo, err := GetCommentInfo(item.Comment, tokenID, l.svcCtx, l.ctx)
		if err != nil {
			e = apivars.SomeDataErr
			if err != apivars.SomeDataErr {
				return
			}
		}
		writer.Write(idxApiComment{
			Comment: videoInfo,
			idx:     item.idx,
		})
	}, func(pipe <-chan idxApiComment, writer mr.Writer[[]types.Comment], cancel func(error)) {
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

type idxComment struct {
	*interaction.Comment
	idx int
}
type idxApiComment struct {
	*types.Comment
	idx int
}
