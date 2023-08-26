package mqs

import (
	"TikTok/apps/user/rpc/internal/logic"
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type PersonalSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPersonalSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *PersonalSuccess {
	return &PersonalSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo 记得创建api和三个新函数的参数检查230824
func (l *PersonalSuccess) Consume(key, val string) error {
	var mqMap map[string][]string

	err := json.Unmarshal([]byte(val), &mqMap)
	if err != nil {
		if errors.Is(err, model.UnmarshalError) {
			return model.UnmarshalError
		}
	}

	for key, value := range mqMap {
		userId, err := strconv.ParseInt(key, 10, 64)

		if err == nil && len(value) == 2 {
			actionType := value[0]
			value := value[1]
			switch actionType {
			case "avatar":
				avatarLogic := logic.NewSetAvatarLogic(l.ctx, l.svcCtx)
				_, err = avatarLogic.SetAvatar(&user.SetAvatarReq{
					UserId: userId,
					Url:    value,
				})
				if err != nil {
					return err
				}

			case "backgroundImage":
				backgroundLogic := logic.NewSetBackgroundImageLogic(l.ctx, l.svcCtx)
				_, err = backgroundLogic.SetBackgroundImage(&user.SetBackgroundImageReq{
					UserId: userId,
					Url:    value,
				})
				if err != nil {
					return err
				}

			case "signature":
				signatureLogic := logic.NewSetSignatureLogic(l.ctx, l.svcCtx)
				_, err = signatureLogic.SetSignature(&user.SetSignatureReq{
					UserId:  userId,
					Content: value,
				})
				if err != nil {
					return err
				}
			}

		} else {
			return err
		}
	}

	return nil
}
