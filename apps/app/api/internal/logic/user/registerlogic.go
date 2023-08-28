package user

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"
	"errors"
	"regexp"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {

	// 参数检查
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,32}$", req.Username) //是否符合用户名格式
	if matched == false {
		return &types.UserRegisterResponse{
			RespStatus: types.RespStatus(apiVars.UsernameRuleError),
		}, nil
	} else if len(req.Password) < 5 || len(req.Password) > 32 {
		return &types.UserRegisterResponse{
			RespStatus: types.RespStatus(apiVars.PasswordRuleError),
		}, nil
	}

	userID, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, model.DuplicateUsername) {
		return &types.UserRegisterResponse{
			RespStatus: types.RespStatus(apiVars.DuplicateUsername),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	token, err := l.svcCtx.JwtAuth.CreateToken(userID.GetUserId())
	if err != nil {
		return nil, err
	}
	return &types.UserRegisterResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserID:     userID.UserId,
		Token:      token,
	}, nil
}
