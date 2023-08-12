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

// SendRelationAction 执行关注/取关操作（好友相关：只要双方互相关注了自动变成好友，否则一方取关则自动解除好友）
func (l *SendRelationActionLogic) SendRelationAction(in *social.RelationActionReq) (*social.RelationActionResp, error) {
	isSucceed := false
	//查询 social 表中是否有这两个用户
	UserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)
	ToUserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.ToUserId)

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

	if in.ActionType == 1 { //如果是执行关注对方行为
		//如果没有关注对方
		if resp.IsFollow == false {
			err = l.svcCtx.CustomDB.UpdateStatusByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId, 1) //修改字段为关注状态
			//user 在 social 表中的 follow_count +1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.UserId, "social", "follow_count", "1")
			//toUser 在 social 表中的 follower_count +1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.ToUserId, "social", "follower_count", "1")
			isSucceed = true
		}

		_, err = l.CheckFollowAndExecute(in, 1)
		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}

	} else if in.ActionType == 2 { //如果是执行取关对方行为
		//如果关注了对方
		if resp.IsFollow == true {
			err = l.svcCtx.CustomDB.UpdateStatusByUserIdAndToUserIdInFollow(l.ctx, in.UserId, in.ToUserId, 0) //修改字段为未关注状态
			//user 在 social 表中的 follow_count -1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.UserId, "social", "follow_count", "-1")
			//toUser 在 social 表中的 follower_count -1
			err = l.svcCtx.CustomDB.AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(l.ctx, in.ToUserId, "social", "follower_count", "-1")
			isSucceed = true
		}

		_, err = l.CheckFollowAndExecute(in, 2)
		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}
	}

	log.Println("result::", isSucceed)

	if err != nil {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		return &social.RelationActionResp{IsSucceed: false}, nil
	}

	return &social.RelationActionResp{IsSucceed: isSucceed}, nil
}

// CheckFollowAndExecute 查询关注表对方是否有关注自己，有则执行对应操作 operationType == 1 ==> add; operationType == 2 ==> remove
func (l *SendRelationActionLogic) CheckFollowAndExecute(in *social.RelationActionReq, operationType int8) (ok bool, err error) {
	ok = true
	//查询关注表对方是否有关注自己
	follow, _ := l.svcCtx.CustomDB.QueryRecordByUserIdAndToUserIdAndStatusInFollow(l.ctx, in.ToUserId, in.UserId, 1)

	isFollowSelf := false
	if follow.Id != 0 {
		isFollowSelf = true
	}

	//查询 friend 表是否有我与对方
	exist, _ := l.svcCtx.CustomDB.QueryIsExistByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)

	switch operationType {
	case 1:

		//如果有关注自己且 friend 表中有相关记录则只更改 status 值成为好友关系，否则插入整条记录
		if isFollowSelf {
			switch exist {
			case true:
				err := l.svcCtx.CustomDB.UpdateRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 1)
				if err != nil {
					logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
					ok = false
				}
			case false:
				err := l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 1)
				if err != nil {
					logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
					ok = false
				}
			}
		}

	case 2:
		//如果 friend 表有对方则更改 status 值解除好友
		if isFollowSelf {
			err := l.svcCtx.CustomDB.UpdateRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 0)
			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				ok = false
			}
		}

	}

	return ok, nil

}
