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

func TestGetRelationFollowList(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	req := &social.RelationFollowListReq{
		UserId: 111,
	}

	res, err := logic.GetRelationFollowList(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println("tested result:", res.UserList)
}
