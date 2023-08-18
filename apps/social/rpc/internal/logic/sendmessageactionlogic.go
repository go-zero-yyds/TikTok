package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SendMessageActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageActionLogic {
	return &SendMessageActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendMessageAction 发送消息
func (l *SendMessageActionLogic) SendMessageAction(in *social.MessageActionReq) (*social.MessageActionResp, error) {

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId, in.ToUserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
	}

	//如果不等于1则return
	if in.ActionType != 1 {
		return &social.MessageActionResp{IsSucceed: false}, nil
	}

	isSucceed := true
	//开启事务
	err := l.svcCtx.CustomDB.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//查询他们是否为好友，否则不能发送消息
		exist, err := l.svcCtx.CustomDB.QueryIsExistByUserIdAndToUserIdInFriend(l.ctx, in.UserId, in.ToUserId)
		if exist != true || err != nil {
			logc.Error(l.ctx, errors.RecordNotFound, in.UserId, in.ToUserId)
			return err
		}

		// 如果为1则插入对应消息记录
		err = l.svcCtx.CustomDB.InsertRecordByUserIdAndToUserIdAndContentInMessage(l.ctx, in.UserId, in.ToUserId, in.Content)
		if err != nil {
			logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
			return err
		}

		return nil
	})
	//事务出错判断
	if err != nil {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId, in.ToUserId)
		isSucceed = false
	} else {
		isSucceed = true
	}

	return &social.MessageActionResp{IsSucceed: isSucceed}, nil
}
