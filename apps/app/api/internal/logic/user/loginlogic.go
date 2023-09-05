package user

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"
	"errors"
	"regexp"

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

	// 参数检查
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,32}$", req.Username) //是否符合用户名格式
	if matched == false {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apivars.ErrUsernameRule),
		}, nil
	} else if len(req.Password) < 5 || len(req.Password) > 32 {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apivars.ErrPasswordRule),
		}, nil
	}

	userID, err := l.svcCtx.UserRPC.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, model.UserValidation) {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apivars.ErrUserValidation),
		}, nil
	}
	if errors.Is(err, model.UserNotFound) {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apivars.ErrUserNotFound),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	token, err := l.svcCtx.JwtAuth.CreateToken(userID.GetUserId())
	if err != nil {
		return nil, err
	}
	return &types.UserLoginResponse{
		RespStatus: types.RespStatus(apivars.Success),
		UserID:     userID.UserId,
		Token:      token,
	}, nil
}
