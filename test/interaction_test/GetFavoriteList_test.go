package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFavoriteList(t *testing.T) {
	//添加测试用例
	for i := 0; i < 100; i++ {
		req := &interaction.FavoriteActionReq{
			UserId:     int64(i / 10),
			VideoId:    int64(i),
			ActionType: 1,
		}
		resp, err := logic.SendFavoriteAction(context.Background(), req)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	//每个用户都拥有自身id + 1到id+10的视频
	for i := 0; i < 10; i++ {
		resp, err := logic.GetFavoriteList(context.Background(), &interaction.FavoriteListReq{
			UserId: int64(i),
		})
		assert.Equal(t, err, nil) //比较视频id 大 -> 小 返回
		for j := i*10 + 9; j >= i*10; j-- {
			assert.Equal(t, int64(j), resp.VideoList[j-i*10])
		}
	}
	//删除测试用例
	for i := 0; i < 100; i++ {
		req := &interaction.FavoriteActionReq{
			UserId:     int64(i / 10),
			VideoId:    int64(i),
			ActionType: 2,
		}
		resp, err := logic.SendFavoriteAction(context.Background(), req)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	// 每个用户无视频点赞信息
	for i := 0; i < 10; i++ {
		resp, err := logic.GetFavoriteList(context.Background(), &interaction.FavoriteListReq{
			UserId: int64(i),
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, 0, len(resp.VideoList))
	}
}
