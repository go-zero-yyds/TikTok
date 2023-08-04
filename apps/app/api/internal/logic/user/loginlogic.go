package user

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/user/rpc/user"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	userID, err := l.svcCtx.UserRPC.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	token, err := l.svcCtx.JwtAuth.CreateToken(userID.GetUserId())
	if err != nil {
		return nil, err
	}
	return &types.UserLoginResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserID:     userID.UserId,
		Token:      token,
	}, nil
}
