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
	"time"
)

func TestGetRelationFriendList(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	log.Println("测试场景：获取好友列表")

	//模拟 ? 个UID并发送请求
	uIdCycleCount := 5 //测试循环的次数

	//创建数据库连接
	conn := GetTestDB()

	for i := 0; i < uIdCycleCount; i++ {
		uid := i + 1
		uid2 := i + 2

		//添加好友数据进行测试
		_, err := conn.ExecCtx(context.Background(), "INSERT INTO `tiktok_social`.`friend` (`user_id`, `to_user_id`, `status`,`version`) VALUES (?, ?, 1,0)", uid, uid2)

		//添加成为好友自动发消息的内容
		_, err = conn.ExecCtx(context.Background(), "INSERT INTO `tiktok_social`.`message` (`from_user_id`, `to_user_id`, `content`,`created_time`) VALUES (?, ?, ?,?)", uid, uid2, "我们成为好友啦，快来聊天吧！", time.Now())

		//开始模拟
		req := &social.RelationFriendListReq{
			UserId: int64(uid),
		}

		res, err := logic.GetRelationFriendList(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		//模拟完成，删除friend测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`friend` WHERE user_id = ? AND to_user_id = ?", uid, uid2)
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`message` WHERE from_user_id = ? AND to_user_id = ? AND content = ?", uid, uid2, "我们成为好友啦，快来聊天吧！")

		//模拟完成，删除social测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`social` WHERE user_id = ?", uid)

		require.NotNil(t, res)
		for i, v := range res.UserList {
			log.Printf("tested result:[%d]: UserId:%d Message:%s MsgType:%d", i, v.UserId, v.Message, v.MsgType)
		}

		log.Printf("tested result%d:%v", i, res.UserList)
	}
}
