package logic

import (
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *social.IsFollowReq) (*social.IsFollowResp, error) {
	//接收关注者ID和被关注者ID
	//查询关注表是否有该记录且status = 1
	exist, err := l.svcCtx.CustomDB.IsFollowByUserIdAndFolloweeId(in.UserId, in.ToUserId)
	log.Println("exist::", exist, in.UserId, in.ToUserId)
	if err != nil {
		//return nil, err
	}

	////返回是否关注 0：未关注 1：已关注
	//
	return &social.IsFollowResp{}, nil
}
