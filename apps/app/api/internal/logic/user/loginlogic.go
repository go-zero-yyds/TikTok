package user

import (
	"TikTok/apps/app/api/apiVars"
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
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$", req.Username) //是否为邮箱格式
	if len(req.Username) > 32 || req.Username == "" || matched == false {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apiVars.UsernameRuleError),
		}, nil
	} else if len(req.Password) < 5 || len(req.Password) > 32 || req.Password == "" {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apiVars.PasswordRuleError),
		}, nil
	}

	userID, err := l.svcCtx.UserRPC.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, model.UserValidation) {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apiVars.UserValidation),
		}, nil
	}
	if errors.Is(err, model.UserNotFound) {
		return &types.UserLoginResponse{
			RespStatus: types.RespStatus(apiVars.UserNotFound),
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
		RespStatus: types.RespStatus(apiVars.Success),
		UserID:     userID.UserId,
		Token:      token,
	}, nil
}
