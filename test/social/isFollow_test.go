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

func TestIsFollow(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	req := &social.IsFollowReq{
		UserId:   111,
		ToUserId: 222,
	}

	res, err := logic.IsFollow(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println("tested result:", res.IsFollow)
}
