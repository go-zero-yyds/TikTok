package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFavoriteCountByVideoId(t *testing.T) {
	//添加测试用例
	for i := 0; i < 100; i++ {
		req := &interaction.FavoriteActionReq{
			UserId:     int64(i + 10),
			VideoId:    int64(i / 10),
			ActionType: 1,
		}
		resp, err := logic.SendFavoriteAction(context.Background(), req)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	//目前每个视频点赞为10
	for i := 0; i < 10; i++ {
		resp, err := logic.GetFavoriteCountByVideoId(context.Background(), &interaction.FavoriteCountByVideoIdReq{
			VideoId: int64(i),
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.FavoriteCount, int64(10))
	}
	//删除测试用例
	for i := 0; i < 100; i++ {
		req := &interaction.FavoriteActionReq{
			UserId:     int64(i + 10),
			VideoId:    int64(i / 10),
			ActionType: 2,
		}
		resp, err := logic.SendFavoriteAction(context.Background(), req)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	// 目前每个视频点赞为0
	for i := 0; i < 10; i++ {
		resp, err := logic.GetFavoriteCountByVideoId(context.Background(), &interaction.FavoriteCountByVideoIdReq{
			VideoId: int64(i),
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.FavoriteCount, int64(0))
	}
}
