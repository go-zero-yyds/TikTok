package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
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

// IsFollow 是否关注
func (l *IsFollowLogic) IsFollow(in *social.IsFollowReq) (*social.IsFollowResp, error) {

	////校验参数是否合规
	//ok, err2 := utils.MatchUID(in.UserId)
	//if err2 != nil || ok != true {
	//	logc.Error(l.ctx, errors.ParamsError, in.UserId, in.ToUserId)
	//	return &social.IsFollowResp{IsFollow: false}, nil
	//}

	//查询 social 表中是否有这两个用户
	UserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)
	ToUserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.ToUserId)

	//如果不存在则直接返回失败
	if UserIdExist == false || ToUserIdExist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
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
