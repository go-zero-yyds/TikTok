package user

import (
	"TikTok/apps/user/rpc/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// 测试用户名存在的情况
	req := &user.LoginReq{
		Username: "testuser",
		Password: "123456",
	}

	resp, err := logic.Login(context.Background(), req)

	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(1690288877564071936), resp.UserId)

	// 测试用户名不存在的情况
	req = &user.LoginReq{
		Username: "non_existing_user",
		Password: "123456",
	}

	resp, err = logic.Login(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Code(100) desc = 用户不存在", err.Error())
}
