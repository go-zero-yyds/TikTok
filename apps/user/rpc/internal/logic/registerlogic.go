package logic

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"TikTok/common/tool"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	// 用户已注册
	if err == nil {
		return nil, model.DuplicateUsername
	} else if !errors.Is(err, model.ErrNotFound) { // 错误
		return nil, err
	} else { // 注册
		snowId := l.svcCtx.Snowflake.Generate().Int64() // 雪花算法生成id
		pwdHash, err := tool.HashAndSalt(in.Password)   // 加盐加密
		if err != nil {
			return nil, err
		}

		_, errInsert := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
			UserId:   snowId,
			Username: in.Username,
			Password: pwdHash,
		})

		if errInsert != nil {
			return nil, err
		}
		return &user.RegisterResp{
			UserId: snowId,
		}, nil
	}
}
