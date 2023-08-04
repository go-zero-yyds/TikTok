package interactiontest

import (
	"TikTok/apps/interaction/rpc/interaction"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCommentList(t *testing.T){
	UserComId := make(map[int64]int64)
	//测试增加评论
	{
		for i := 0; i < 100; i++ {
			cmt := fmt.Sprintf("user = %d comment videoId = %d ", i, i/10)
			l := new(string)
			*l = cmt
			resp, err := logic.SendCommentAction(context.Background(), &interaction.CommentActionReq{
				UserId:      int64(i),
				VideoId:     int64(i / 10),
				ActionType:  1,
				CommentText: l,
			})
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.Comment.UserId, int64(i))
			assert.Equal(t, resp.Comment.Content, cmt)
			UserComId[int64(i)] = resp.Comment.Id
		}
	}
	//检测评论内容
	{
		for i := 0; i < 9; i++{
			ret , err := logic.GetCommentList(context.Background() , &interaction.CommentListReq{
				UserId: int64(i),
				VideoId: int64(i),
			})
			assert.Equal(t , nil , err)
			for j := 9; j >= 0 ; j--{
				//cmt := fmt.Sprintf("user = %d comment videoId = %d ", i *  10 + j, i)
				//assert.Equal(t , cmt , ret.CommentList[9 - j].Content)
				assert.Equal(t , time.Now().Format("01-02") , ret.CommentList[9 - j].CreateDate)
			}
		}
	}

	//测试根据id删除评论
	{
		for i := 0; i < 100; i++ {
			tmp := UserComId[int64(i)]
			resp, err := logic.SendCommentAction(context.Background(), &interaction.CommentActionReq{
				UserId:      int64(i),
				VideoId:     int64(i / 10),
				ActionType:  2,
				CommentId: &tmp,
				
			})
			assert.Equal(t, err, nil)
			assert.NotEqual(t, resp.Comment.Id , tmp)
		}
	}
	//每个视频没有评论
	{
		for i := 0; i < 10; i++{
			count , err := logic.GetCommentCountByVideoId(context.Background() , &interaction.CommentCountByVideoIdReq{
				VideoId: int64(i),
			})
			assert.Equal(t , nil , err)
			assert.Equal(t , int64(0) , count.CommentCount)
		}

	}
}
