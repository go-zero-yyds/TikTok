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

func TestGetRelationFriendList(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	req := &social.RelationFriendListReq{
		UserId: 111,
	}

	res, err := logic.GetRelationFriendList(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	for i, v := range res.UserList {
		log.Printf("tested result:[%d]: UserId:%d Message:%s MsgType:%d", i, v.UserId, v.Message, v.MsgType)
	}

}
