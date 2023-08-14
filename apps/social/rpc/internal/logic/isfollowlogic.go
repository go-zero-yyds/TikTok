package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
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

type IsFollowFunc func(*social.IsFollowReq) (*social.IsFollowResp, error)

// IsFollow 是否关注
func (l *IsFollowLogic) IsFollow(in *social.IsFollowReq) (*social.IsFollowResp, error) {

	//验证用户存在性并注册，否则返回失败
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId, in.ToUserId)
	if ok != true {
		return &social.IsFollowResp{IsFollow: false}, nil
	}

	//查询 user 在 follow 表中的字段
	followStruct, err := l.svcCtx.CustomDB.QueryRecordByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId)
	//如果 id 为 0 说明没查到
	if followStruct.Id == 0 {
		//插入该字段
		err = l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId)
		//直接返回未关注
		return &social.IsFollowResp{IsFollow: false}, nil
	}
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
	}

	//返回关注状态
	followStatus := false

	if followStruct.Status[0] == 1 {
		followStatus = true
	}
	return &social.IsFollowResp{IsFollow: followStatus}, nil
}
