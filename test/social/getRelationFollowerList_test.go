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

func TestGetRelationFollowerList(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	log.Println("测试场景：获取粉丝列表（该测试可能会有bug，请以实际正式环境/apiFox为准）")

	//模拟 ? 个UID并发送请求
	uIdCycleCount := 10 //测试循环的次数

	//创建数据库连接
	conn := GetTestDB()

	for i := 0; i < uIdCycleCount; i++ {
		uid := i + 1
		uid2 := i + 2

		//添加被关注数据进行测试
		_, err := conn.ExecCtx(context.Background(), "INSERT INTO `tiktok_social`.`follow` (`user_id`, `to_user_id`, `status`,`version`) VALUES (?, ?, 1,0)", uid2, uid)

		//开始模拟
		req := &social.RelationFollowerListReq{
			UserId: int64(uid),
		}

		res, err := logic.GetRelationFollowerList(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		////模拟完成，删除follow被关注测试数据记录
		//_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`follow` WHERE user_id = ? AND to_user_id = ?", uid2, uid)
		//
		////模拟完成，删除social测试数据记录
		//_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`social` WHERE user_id = ? OR user_id = ?", uid2, uid)

		require.NotNil(t, res)
		log.Printf("tested result%d:%v", i, res.UserList)
	}

}
