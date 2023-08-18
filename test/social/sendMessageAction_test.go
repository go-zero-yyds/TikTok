/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/12 17:21
 */

package social

import (
	"TikTok/apps/social/rpc/social"
	"TikTok/test/social/utils"
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSendMessageAction(t *testing.T) {
	log.Println("--------------------Testing--------------------")

	log.Println("测试场景：在成为好友的情况下模拟发送消息")
	//模拟 ? 个UID并发送请求
	uIdCycleCount := 5 //测试循环的次数
	uIdLength := 19    //UID长度

	//创建数据库连接
	conn := GetTestDB()

	for i := 0; i < uIdCycleCount; i++ {
		uid := i + 1
		uid2 := i + 2
		//随机内容
		content := utils.GenerateContent(uIdLength)

		//将uid插入friend表
		_, err := conn.ExecCtx(context.Background(), "INSERT INTO `tiktok_social`.`friend` (`user_id`, `to_user_id`, `status`,`version`) VALUES (?, ?, 1,0)", uid, uid2)

		//开始模拟
		req := &social.MessageActionReq{
			UserId:     int64(uid),
			ToUserId:   int64(uid2),
			ActionType: 1,
			Content:    content,
		}

		res, err := logic.SendMessageAction(context.Background(), req)
		if err != nil {
			log.Fatalln("err :", err)
		}

		//模拟完成，删除friend测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`friend` WHERE user_id = ? AND to_user_id = ?", uid, uid2)

		//模拟完成，删除message测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`message` WHERE from_user_id = ? AND to_user_id = ? AND content = ?", uid, uid2, "我们成为好友啦，快来聊天吧！")

		//模拟完成，删除social测试数据记录
		_, err = conn.ExecCtx(context.Background(), "DELETE FROM `tiktok_social`.`social` WHERE user_id = ?", uid)

		require.NotNil(t, res)
		log.Printf("tested result%d:%v", i, res.IsSucceed)
	}
}
