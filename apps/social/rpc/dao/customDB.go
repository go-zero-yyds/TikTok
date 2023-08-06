package dao

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/net/context"
)

type CustomDB struct {
	message MessageModel
	follow  FollowModel
	social  SocialModel
	conn    sqlx.SqlConn
}

// NewCustomDB 初始化自定义dao数据库
func NewCustomDB(conn sqlx.SqlConn) *CustomDB {
	return &CustomDB{
		message: NewMessageModel(conn),
		follow:  NewFollowModel(conn),
		social:  NewSocialModel(conn),
		conn:    conn,
	}
}

// IsFollowByUserIdAndFolloweeId 查询：userId & followId & status==1
func (db *CustomDB) IsFollowByUserIdAndFolloweeId(userId, toUserId int64) (bool, error) {
	var follow Follow
	tableName := "follow"
	//todo 查询结果扫描不到结构体里
	query := fmt.Sprintf("select user_id from %s where user_id = ? and to_user_id = ? and status = 1", tableName)
	err := db.conn.QueryRowCtx(context.Background(), &follow, query, userId, toUserId)
	if err != nil {
		logc.Error(context.Background(), userId, toUserId, err)
		return false, err
	}
	return true, nil
}

// FollowAction 新增/修改： A关注/取关B
//func (db *CustomDB) FollowAction(userId, videoId int64) (bool, error) {
//
//	exist, err := db.follow.
//	if err != nil && err != dao.ErrNotFound {
//		logc.Error(context.Background(), userId, videoId, err)
//		return false, err
//	}
//	return exist, nil
//}
