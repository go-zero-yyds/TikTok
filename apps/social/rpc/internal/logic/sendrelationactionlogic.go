package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

type SendRelationActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRelationActionLogic {
	return &SendRelationActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendRelationAction 执行关注/取关操作
func (l *SendRelationActionLogic) SendRelationAction(in *social.RelationActionReq) (*social.RelationActionResp, error) {
	//查询 social 表中是否有这两个用户
	UserIdExist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.UserId)
	ToUserIdExist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.ToUserId)

	//如果不存在则直接返回失败
	if UserIdExist == false || ToUserIdExist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
		return &social.RelationActionResp{IsSucceed: false}, nil
	}

	//调用同层 isFollowLogic 直接查询用户是否关注了对方
	isFollowLogic := NewIsFollowLogic(l.ctx, l.svcCtx)
	resp, _ := isFollowLogic.IsFollow(&social.IsFollowReq{
		UserId:   in.UserId,
		ToUserId: in.ToUserId,
	})

	isSuccesses := false

	//如果是执行关注对方行为
	if in.ActionType == 1 {
		//如果没有关注对方
		if resp.IsFollow == false {
			err = l.svcCtx.CustomDB.UpdateStatusByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId, 1) //修改字段为关注状态
			//user 在 social 表中的 follow_count +1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.UserId, "social", "follow_count", "1")
			//toUser 在 social 表中的 follower_count +1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.ToUserId, "social", "follower_count", "1")
			isSuccesses = true
		}

		//查询关注表对方是否有关注自己
		follow, err := l.svcCtx.CustomDB.QueryRecordByUserIdAndToUserIdInFollow(l.ctx, in.ToUserId, in.UserId)
		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}

		//如果有关注自己则成为好友关系
		if follow.Id != 0 {
			err := l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)
			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return &social.RelationActionResp{IsSucceed: false}, nil
			}
		}

	}
	//如果是执行取关对方行为
	if in.ActionType == 2 {
		//如果关注了对方
		if resp.IsFollow == true {
			err = l.svcCtx.CustomDB.UpdateStatusByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId, 0) //修改字段为未关注状态
			//user 在 social 表中的 follow_count -1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.UserId, "social", "follow_count", "-1")
			//toUser 在 social 表中的 follower_count -1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.ToUserId, "social", "follower_count", "-1")
			isSuccesses = true
		}

		//查询关注表对方是否有关注自己
		follow, err := l.svcCtx.CustomDB.QueryRecordByUserIdAndToUserIdInFollow(l.ctx, in.ToUserId, in.UserId)
		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}

		log.Printf("+%v", follow)

		//有则解除好友
		if follow.Id != 0 {
			err := l.svcCtx.CustomDB.DeleteRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)
			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return &social.RelationActionResp{IsSucceed: false}, nil
			}
		}
	}

	if err != nil {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
	}

	return &social.RelationActionResp{IsSucceed: isSuccesses}, nil
}
