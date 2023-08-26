package logic

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"TikTok/pkg/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	//测试
	mqMap := make(map[string][]string, 10)
	mqMap["1693651455589748736"] = []string{"signature", "www.baidu1234.com"}
	marshal, _ := json.Marshal(mqMap)
	s := string(marshal)
	log.Println("s:", s)
	if err := l.svcCtx.KqPusherClient.Push(s); err != nil {
		logx.Errorf("KqPusherClient Push Error , err :%v", err)
	}
	log.Println("mq发过去了")

	res, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.UserNotFound
		}
		return nil, err
	}

	// 验证密码
	isCorrect := tool.ComparePasswords(res.Password, in.Password)
	if !isCorrect {
		return nil, model.UserValidation
	}
	return &user.LoginResp{
		UserId: res.UserId,
	}, nil
	// return &user.LoginResp{}, nil
}
