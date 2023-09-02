package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"testing"
)

// 点赞

func TestSendFavoriteActionLogic(t *testing.T) {
	// 测试场景1: 成功情况
	req := &interaction.FavoriteActionReq{
		UserId:     1,
		VideoId:    1,
		ActionType: 1,
	}
	resp, err := logic.SendFavoriteAction(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !resp.IsSucceed {
		t.Fatalf("Expected IsSucceed to be true, but got false")
	}

	// 测试场景2: 再点赞成功情况
	req = &interaction.FavoriteActionReq{
		UserId:     1,
		VideoId:    1,
		ActionType: 1,
	}
	resp, err = logic.SendFavoriteAction(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.IsSucceed {
		t.Fatalf("Expected IsSucceed to be true, but got false")
	}
	// 测试场景3：取消操作

	req = &interaction.FavoriteActionReq{
		UserId:     1,
		VideoId:    1,
		ActionType: 2,
	}
	resp, err = logic.SendFavoriteAction(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !resp.IsSucceed {
		t.Fatalf("Expected IsSucceed to be true, but got false")
	}
}
