package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"sync"
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

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId, in.ToUserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
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
	}

	return &social.RelationActionResp{IsSucceed: isSucceed}, nil
}

// UpdateUserSocialCount update user socialCount(followList and followerList) in social table operationType == 1 ==> add; operationType == 2 ==> remove
func (l *SendRelationActionLogic) UpdateUserSocialCount(in *social.RelationActionReq, operationType int8) (ok bool, err error) {
	switch operationType {
	case 1:
		//开启事务
		err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			err := l.svcCtx.CustomDB.TransactionUpdateSocialCount(l.ctx, session, in.UserId, in.ToUserId, 1, 1)
			if err != sql.ErrNoRows && err != nil {
				return err
			}
			ok, err := l.CheckFollowAndExecute(in, 1)
			if ok != true || err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return err
			}
			return nil
		})

		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
			return false, err
		}

		return true, err
	case 2:
		//开启事务
		err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			err := l.svcCtx.CustomDB.TransactionUpdateSocialCount(l.ctx, session, in.UserId, in.ToUserId, 0, -1)
			if err != sql.ErrNoRows && err != nil {
				return err
			}
			ok, err := l.CheckFollowAndExecute(in, 2)
			if ok != true || err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				return err
			}
			return nil
		})

		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
			return false, err
		}
		return true, err
	}

	return true, nil

}

// CheckFollowAndExecute 查询关注表对方是否有关注自己，有则执行对应操作（成为好友/解除好友） operationType == 1 ==> add; operationType == 2 ==> remove
func (l *SendRelationActionLogic) CheckFollowAndExecute(in *social.RelationActionReq, operationType int8) (ok bool, err error) {
	ok = true

	isFollowSelfChan := make(chan bool, 100)
	existChan := make(chan bool, 100)
	defer close(isFollowSelfChan)
	defer close(existChan)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		//查询关注表对方是否有关注自己
		follow, _ := l.svcCtx.CustomDB.QueryRecordByUserIdAndToUserIdAndStatusInFollow(l.ctx, in.ToUserId, in.UserId, 1)
		if follow.Id != 0 {
			isFollowSelfChan <- true
			return
		}
		isFollowSelfChan <- false
		return
	}()

	go func() {
		defer wg.Done()
		//查询 friend 表是否有我与对方
		exist, _ := l.svcCtx.CustomDB.QueryIsExistByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)
		existChan <- exist
		return
	}()

	wg.Wait()

	isFollowSelf := <-isFollowSelfChan
	exist := <-existChan

	switch operationType {
	case 1:
		//如果有关注自己且 friend 表中有相关记录则只更改 status 值成为好友关系，否则插入整条记录
		if isFollowSelf {
			switch exist {
			case true:
				//开启事务
				err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
					err := l.svcCtx.CustomDB.LockRecordBySessionInFriend(l.ctx, session, in.UserId, in.ToUserId)
					if err != nil {
						return err
					}
					err = l.svcCtx.CustomDB.UpdateRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 1)
					if err != sql.ErrNoRows && err != nil {
						logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
						ok = false
						return err
					}
					ok = true
					return nil
				})

				if err != nil {
					logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
					ok = false
					return
				}

			case false:
				err := l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 1)
				if err != nil {
					logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
					ok = false
				}
			}

			//成为好友自动发起一条消息聊天，以方便 getRelationFriendList 的调用
			err = l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdAndContentInMessage(l.ctx, in.UserId, in.ToUserId, "我们成为好友啦，快来聊天吧！")
			if err != nil {
				logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
				ok = false
				return
			}
		}

	case 2:
		//如果 friend 表有对方则更改 status 值解除好友
		if isFollowSelf {
			//开启事务
			err = l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				err := l.svcCtx.CustomDB.LockRecordBySessionInFriend(l.ctx, session, in.UserId, in.ToUserId)
				if err != nil {
					return err
				}
				err = l.svcCtx.CustomDB.UpdateRecordByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId, 0)
				if err != nil {
					logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
					ok = false
					return err
				}
				return nil
			})
		}
	}

	return ok, nil
}
