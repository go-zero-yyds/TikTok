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

	log.Println("测试场景：关注对方并自动使其变为好友")
	//模拟 ? 个UID并发送请求
	uIdCount := 100 //生成UID的次数（即循环次数）

	//创建数据库连接
	conn := GetTestDB()

	for i := 0; i < uIdCount; i++ {
		uid := i + 1
		uid2 := i + 2

		//a关注b
		//开始模拟
		req := &social.RelationActionReq{
			UserId:     int64(uid),
			ToUserId:   int64(uid2),
			ActionType: 1,
		}

		res, err := logic.SendRelationAction(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		//b关注a
		req = &social.RelationActionReq{
			UserId:     int64(uid2),
			ToUserId:   int64(uid),
			ActionType: 1,
		}

		res, err = logic.SendRelationAction(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		//可以发现friend表中自动将互关的人变成了好友（需要注释掉下方的DELETE才能展现出来）
		//模拟完成，删除测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`follow` WHERE user_id = ? AND to_user_id = ?", uid, uid2)
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`follow` WHERE user_id = ? AND to_user_id = ?", uid2, uid)
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`friend` WHERE user_id = ? AND to_user_id = ?", uid2, uid)
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`social` WHERE user_id in (?,?)", uid, uid2)
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`message` WHERE from_user_id = ? AND to_user_id = ? AND content = ?", uid2, uid, "我们成为好友啦，快来聊天吧！")

		require.NotNil(t, res)
		log.Printf("tested result%d:%v", i, res.IsSucceed)
	}
}
