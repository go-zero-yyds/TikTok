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

func TestGetFollowCount(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	log.Println("测试场景：获取关注数")

	//创建数据库连接
	conn := GetTestDB()

	//获取social表当前存在的用户的ID
	var userList []int64
	_ = conn.QueryRowsPartialCtx(context.Background(), &userList, "SELECT user_id FROM `tiktok_social`.`social`")
	log.Println(userList)

	for i := 0; i < len(userList); i++ {
		req := &social.FollowCountReq{
			UserId: userList[i],
		}

		res, err := logic.GetFollowCount(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		require.NotNil(t, res)
		log.Printf("tested result[%d]:%d", userList[i], res.FollowCount)
	}

}
