package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const n = 1000 // 测试数量

func TestIsFavorite(t *testing.T) {
	//-------------------------准备环境-----------------------------//
	// 准备测试用例数据
	for i := 0; i < n; i++ {
		{
			resp, err := logic.SendFavoriteAction(context.Background(), &interaction.FavoriteActionReq{
				UserId:     int64(i),
				VideoId:    int64(i),
				ActionType: 1,
			})
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
				IsSucceed: true,
			}.IsSucceed)
		}
	}
	//-------------------------测试-----------------------------//
	for i := 0; i < n; i++ {
		// 调用被测试的接口方法 查看是否点赞
		{
			resp, err := logic.IsFavorite(context.Background(), &interaction.IsFavoriteReq{
				UserId:  int64(i),
				VideoId: int64(i),
			})

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			// 检查返回结果是否正确
			if !resp.GetIsFavorite() {
				t.Errorf("Unexpected value: %v", resp.GetIsFavorite())
			}
		}
	}
	//-------------------------清除环境-----------------------------//
	//删除数据
	for i := 0; i < n; i++ {
		{
			resp, err := logic.SendFavoriteAction(context.Background(), &interaction.FavoriteActionReq{
				UserId:     int64(i),
				VideoId:    int64(i),
				ActionType: 2,
			})
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
				IsSucceed: true,
			}.IsSucceed)
		}
	}
	//-------------------------测试-----------------------------//
	for i := 0; i < n; i++ {
		// 调用被测试的接口方法  未点赞
		{
			resp, err := logic.IsFavorite(context.Background(), &interaction.IsFavoriteReq{
				UserId:  int64(i),
				VideoId: int64(i),
			})

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			// 检查返回结果是否正确
			if resp.GetIsFavorite() {
				t.Errorf("Unexpected value: %v", resp.GetIsFavorite())
			}
		}
	}
}
