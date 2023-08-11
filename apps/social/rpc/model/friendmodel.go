package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FriendModel = (*customFriendModel)(nil)

type (
	// FriendModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendModel.
	FriendModel interface {
		friendModel
	}

	customFriendModel struct {
		*defaultFriendModel
	}
)

// NewFriendModel returns a model for the database table.
func NewFriendModel(conn sqlx.SqlConn) FriendModel {
	return &customFriendModel{
		defaultFriendModel: newFriendModel(conn),
	}
}
