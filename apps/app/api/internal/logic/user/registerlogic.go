package user

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/user/rpc/user"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

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
	// todo: add your logic here and delete this line
	userID, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterReq{
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
	return &types.UserRegisterResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserID:     userID.UserId,
		Token:      token,
	}, nil
}
