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

func TestSendRelationAction(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	req := &social.RelationActionReq{
		UserId:     333,
		ToUserId:   222,
		ActionType: 1,
	}

	res, err := logic.SendRelationAction(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println("tested result:", res.IsSucceed)
}
