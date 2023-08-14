package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	//查询 social 表中是否有这两个用户
	UserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)
	ToUserIdExist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.ToUserId)

	//如果不存在则注册进social表内并重试
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

	isSucceed := false

	if in.ActionType == 1 { //如果是执行关注对方行为
		//如果没有关注对方
		if resp.IsFollow == false {
			ok, err := l.UpdateUserSocialCount(in, 1)
			if err != nil || ok != true {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return &social.RelationActionResp{IsSucceed: false}, nil
			}
			isSucceed = true
		}

		ok, err := l.CheckFollowAndExecute(in, 1)
		if ok != true || err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}

	} else if in.ActionType == 2 { //如果是执行取关对方行为
		//如果关注了对方
		if resp.IsFollow == true {
			ok, err := l.UpdateUserSocialCount(in, 2)
			if err != nil || ok != true {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return &social.RelationActionResp{IsSucceed: false}, nil
			}
			isSucceed = true
		}

		ok, err := l.CheckFollowAndExecute(in, 2)
		if ok != true || err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		}
	}

	if err != nil {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		return &social.RelationActionResp{IsSucceed: false}, nil
	}

	return &social.RelationActionResp{IsSucceed: isSucceed}, nil
}

// UpdateUserSocialCount update user socialCount(followList and followerList) in social table operationType == 1 ==> add; operationType == 2 ==> remove
func (l *SendRelationActionLogic) UpdateUserSocialCount(in *social.RelationActionReq, operationType int8) (ok bool, err error) {
	resultChan := make(chan bool)
	defer close(resultChan)
	go func(in *social.RelationActionReq, operationType int8) {
		switch operationType {
		case 1:
			//开启事务
			err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				err := l.svcCtx.CustomDB.TransactionUpdateSocialCount(l.ctx, session, in.UserId, in.ToUserId, 1, 1)
				if err != sql.ErrNoRows {
					return err
				}
				return nil
			})

			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				resultChan <- false
				return
			}
			resultChan <- true
			return
		case 2:
			//开启事务
			err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				err := l.svcCtx.CustomDB.TransactionUpdateSocialCount(l.ctx, session, in.UserId, in.ToUserId, 0, -1)
				if err != sql.ErrNoRows {
					return err
				}
				return nil
			})

			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				resultChan <- false
				return
			}
			resultChan <- true
			return
		}
	}(in, operationType)

	select {
	case result := <-resultChan:
		if result == false {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
			return false, err
		}
		return true, nil
	}
}

// CheckFollowAndExecute 查询关注表对方是否有关注自己，有则执行对应操作（成为好友/解除好友） operationType == 1 ==> add; operationType == 2 ==> remove
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
