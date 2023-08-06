package dao

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SocialModel = (*customSocialModel)(nil)

type (
	// SocialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSocialModel.
	SocialModel interface {
		socialModel
	}

	customSocialModel struct {
		*defaultSocialModel
	}
)

// NewSocialModel returns a model for the database table.
func NewSocialModel(conn sqlx.SqlConn) SocialModel {
	return &customSocialModel{
		defaultSocialModel: newSocialModel(conn),
	}
}
