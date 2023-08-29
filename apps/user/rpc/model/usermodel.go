package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UpdateByUserId(ctx context.Context, data *User, field string) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *defaultUserModel) UpdateByUserId(ctx context.Context, user *User, field string) error {
	data, err := m.FindOne(ctx, user.UserId)
	if err != nil {
		return err
	}

	if data != nil {
		switch field {
		case "avatar":
			data.Avatar = user.Avatar
		case "backgroundImage":
			data.BackgroundImage = user.BackgroundImage
		case "signature":
			data.Signature = user.Signature
		}
		err := m.Update(ctx, data)
		if err != nil {
			return err
		}
	}

	return err
}
