package user

import (
	"TikTok/apps/user/rpc/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// 测试用户名正确的情况
	req := &user.RegisterReq{
		// Username: "existing_user",
		Username: "testuser",
		Password: "123456",
	}

	resp, err := logic.Register(context.Background(), req)

	assert.NotNil(t, resp)
	assert.NoError(t, err)

	// 测试用户名重复的情况
	resp, err = logic.Register(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Code(100) desc = 用户已存在", err.Error())

}
