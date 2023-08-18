/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/12 17:21
 */

package social

import (
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestGetMessages(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	log.Println("测试场景：获取消息列表")

	req := &social.MessageChatReq{
		UserId:     111,
		ToUserId:   222,
		PreMsgTime: 1691750037000, //毫秒级
	}

	res, err := logic.GetMessages(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	for i, v := range res.MessageList {
		log.Printf("tested result:[%d]: %+v", i, v)
	}
}