package common

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Validator struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidator(ctx context.Context, svcCtx *svc.ServiceContext) *Validator {
	return &Validator{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ValidateUserInSocial todo 可优化
// ValidateUserInSocial 用户存在性检查
func (va *Validator) ValidateUserInSocial(userIdMap *map[int64]bool) (ok bool, err error) {

	ok = true
	for userId, _ := range *userIdMap {
		exist, err := va.svcCtx.CustomDB.QueryUserIdIsExistInSocial(va.ctx, userId)
		if err != nil || exist != true {
			logc.Error(va.ctx, errors.RecordNotFound, userId)
			(*userIdMap)[userId] = false
			//如果有任何一个不存在的则设置为false
			ok = false
			continue
		}
		(*userIdMap)[userId] = true
	}
	return ok, nil
}
