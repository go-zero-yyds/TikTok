package mqs

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/pkg/tool"
	"context"
	"encoding/json"
	"errors"
	"regexp"
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

func (l *PersonalSuccess) Consume(_, val string) error {
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
				err := l.ToSetAvatar(userId, value)
				if err != nil {
					return err
				}

			case "backgroundImage":
				err := l.ToSetBackgroundImage(userId, value)
				if err != nil {
					return err
				}

			case "signature":
				err := l.ToSetSignature(userId, value)
				if err != nil {
					return err
				}

			case "username":
				err := l.ToRegisterBot(userId, value)
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

// ToRegisterBot 注册机器人
func (l *PersonalSuccess) ToRegisterBot(userId int64, username string) error {
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)
	// 用户已注册
	if err == nil {
		return model.DuplicateUsername
	} else if !errors.Is(err, model.ErrNotFound) { // 错误
		return err
	} else { // 注册
		pwdHash, err := tool.HashAndSalt("0") // 加盐加密（机器人默认密码为0，因此无法登录机器人）
		if err != nil {
			return err
		}

		_, errInsert := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
			UserId:          userId,
			Username:        username,
			Password:        pwdHash,
			Avatar:          "avatar/bot",
			BackgroundImage: "backgroundImage/bot",
			Signature:       "滴滴...抱歉，我...不是很擅长说话...会尽力服务...",
		})

		if errInsert != nil {
			return errInsert
		}
	}
	return nil
}

// ToSetSignature 设置签名
func (l *PersonalSuccess) ToSetSignature(userId int64, content string) error {
	//限制签名内容的长度
	matched, err := regexp.MatchString("^.{1,50}$", content)
	if err != nil {
		return err
	}
	if matched == false {
		return nil
	}
	err = l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:    userId,
		Signature: content,
	}, "signature")

	return nil
}

// ToSetBackgroundImage 设置背景图片
func (l *PersonalSuccess) ToSetBackgroundImage(userId int64, url string) error {
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:          userId,
		BackgroundImage: url,
	}, "backgroundImage")

	return err
}

// ToSetAvatar 设置头像
func (l *PersonalSuccess) ToSetAvatar(userId int64, url string) error {
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId: userId,
		Avatar: url,
	}, "avatar")

	return err
}
