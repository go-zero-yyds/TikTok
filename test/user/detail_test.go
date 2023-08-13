package user

import (
	"TikTok/apps/user/rpc/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetail(t *testing.T) {
	userId := 1690288877564071936
	req := &user.BasicUserInfoReq{
		UserId: int64(userId),
	}

	res, err := logic.Detail(context.Background(), req)

	testRes := &user.BasicUserInfoResp{
		User: &user.BasicUserInfo{
			Id:   int64(userId),
			Name: "testuser",
		},
	}

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.User.Id, testRes.User.Id)
	assert.Equal(t, res.User.Name, testRes.User.Name)
}
