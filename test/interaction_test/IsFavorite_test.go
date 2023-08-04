package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestIsFavorite(t *testing.T) {
	// 准备测试用例数据
	userId := 1
	videoId := 1
	//添加测试数据
	{
		resp, err := logic.SendFavoriteAction(context.Background(), &interaction.FavoriteActionReq{
			UserId:     int64(userId),
			VideoId:    int64(videoId),
			ActionType: 1,
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	// 调用被测试的接口方法 点赞
	{
		resp, err := logic.IsFavorite(context.Background(), &interaction.IsFavoriteReq{
			UserId:  int64(userId),
			VideoId: int64(videoId),
		})

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// 检查返回结果是否正确
		if !resp.GetIsFavorite() {
			t.Errorf("Unexpected value: %v", resp.GetIsFavorite())
		}
	}
	{
		resp, err := logic.SendFavoriteAction(context.Background(), &interaction.FavoriteActionReq{
			UserId:     int64(userId),
			VideoId:    int64(videoId),
			ActionType: 2,
		})
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.IsSucceed, interaction.FavoriteActionResp{
			IsSucceed: true,
		}.IsSucceed)
	}
	// 调用被测试的接口方法  未点赞
	{
		resp, err := logic.IsFavorite(context.Background(), &interaction.IsFavoriteReq{
			UserId:  int64(userId),
			VideoId: int64(videoId),
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
