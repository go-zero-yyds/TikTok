package common

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateAndRegisterStruct struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateAndRegisterStruct(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateAndRegisterStruct {
	return &ValidateAndRegisterStruct{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ValidateAndRegister 检测用户存在性并注册
func (v *ValidateAndRegisterStruct) ValidateAndRegister(userIdList ...int64) (ok bool) {

	////配置验证的分布式锁
	//validatorLock := redis.NewRedisLock(v.svcCtx.Redis, "validator_Lock")
	//validatorLock.SetExpire(5)
	//
	////如果获取到锁
	//for {
	//	acquire, _ := validatorLock.AcquireCtx(v.ctx)
	//	if acquire {
	//		break
	//	}
	//}
	//defer validatorLock.ReleaseCtx(v.ctx)

	ok = true
	userIdMap := make(map[int64]bool)

	//将userId导入map
	for _, v := range userIdList {
		userIdMap[v] = false
	}

	//使用校验器校验用户是否在 social 表内
	validator := NewValidator(v.ctx, v.svcCtx)
	ok, err := validator.ValidateUserInSocial(&userIdMap)
	//如果有任何一个不存在则自动注册进表
	if ok != true || err != nil {
		logc.Error(v.ctx, errors.RecordNotFound, userIdMap)
		newRegisterUser := NewRegisterUser(v.ctx, v.svcCtx)
		err := newRegisterUser.RegisterUserToSocial(userIdMap)
		//如果注册过程出错传输失败
		if err != nil {
			logc.Error(v.ctx, errors.SQLOperateFailed, userIdMap)
			ok = false
			return
		}
		ok = true
	}

	return ok
}
