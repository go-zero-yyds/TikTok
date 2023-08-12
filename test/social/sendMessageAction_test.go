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

func TestSendMessageAction(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	req := &social.MessageActionReq{
		UserId:     111,
		ToUserId:   222,
		ActionType: 1,
		Content:    "hello FxShadow and his group",
	}

	res, err := logic.SendMessageAction(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println("tested result:", res.IsSucceed)
}
