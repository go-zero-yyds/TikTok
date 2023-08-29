package mqs

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
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
			}

		} else {
			return err
		}
	}

	return nil
}

// ToSetSignature 设置签名
func (l *PersonalSuccess) ToSetSignature(userId int64, content string) error {
	//限制签名内容的长度
	matched, err := regexp.MatchString("^.{1,50}$", content)

	err = l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:    userId,
		Signature: content,
	}, "signature")
	var data string
	if err != nil || matched != true {
		data = l.Tidy(userId, "signature", "false")
	} else {
		data = l.Tidy(userId, "signature", "true")
	}

	//推送消息
	return l.Push(data)
}

// ToSetBackgroundImage 设置背景图片
func (l *PersonalSuccess) ToSetBackgroundImage(userId int64, url string) error {
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:          userId,
		BackgroundImage: url,
	}, "backgroundImage")
	var data string
	if err != nil {
		data = l.Tidy(userId, "backgroundImage", "false")
	} else {
		data = l.Tidy(userId, "backgroundImage", "true")
	}

	//推送消息
	return l.Push(data)
}

// ToSetAvatar 设置头像
func (l *PersonalSuccess) ToSetAvatar(userId int64, url string) error {
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId: userId,
		Avatar: url,
	}, "avatar")
	var data string
	if err != nil {
		data = l.Tidy(userId, "avatar", "false")
	} else {
		data = l.Tidy(userId, "avatar", "true")
	}

	//推送消息
	return l.Push(data)
}

// Tidy 整合消息
func (l *PersonalSuccess) Tidy(userId int64, key, value string) (data string) {
	message := make(map[string][]string)
	message[strconv.FormatInt(userId, 10)] = []string{key, value}
	marshal, _ := json.Marshal(message)
	data = string(marshal)
	return data
}

// Push 推送消息
func (l *PersonalSuccess) Push(data string) error {
	if err := l.svcCtx.KqPusherClient.Push(data); err != nil {
		logx.Errorf("KqPusherClient Push Error , err :%v", err)
	}
	return nil
}
