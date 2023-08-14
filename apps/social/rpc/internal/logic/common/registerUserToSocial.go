/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/13 20:52
 */

package common

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUser struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterUser(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUser {
	return &RegisterUser{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RegisterUserToSocial todo 可优化
// RegisterUserToSocial 用户注册入库
func (re *RegisterUser) RegisterUserToSocial(userIdMap map[int64]bool) (err error) {
	for userId, value := range userIdMap {
		if value == false {
			err := re.svcCtx.CustomDB.InsertRecordByUserIdInSocial(re.ctx, userId, 0, 0)
			if err != nil {
				logc.Error(re.ctx, errors.RecordNotFound, userId, value)
			}
		}
	}
	return err
}
