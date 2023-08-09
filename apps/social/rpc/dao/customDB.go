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

// QueryUsersOfFriendListByUserId query：userId ==> friendList ==> friendList.userId.fields.values
func (db *CustomDB) QueryUsersOfFriendListByUserId(ctx context.Context, userId int64) (socialList []Social, err error) {
	tableNameA := "friend"
	tableNameB := "social"
	query := fmt.Sprintf("SELECT s.user_id FROM %s AS f JOIN %s AS s ON f.to_user_id = s.user_id WHERE f.user_id = ? AND f.status = ?", tableNameA, tableNameB)
	err = db.conn.QueryRowsPartialCtx(ctx, &socialList, query, userId, 1)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return socialList, err
}

// QueryUsersOfFollowerListByUserId query：userId ==> followerList ==> followerList.userId.fields.values
func (db *CustomDB) QueryUsersOfFollowerListByUserId(ctx context.Context, userId int64) (userList []int64, err error) {
	tableNameA := "follow"
	tableNameB := "social"
	//query := fmt.Sprintf("SELECT s.id, s.user_id, s.follow_count, s.follower_count, s.total_favorited, s.work_count, s.favorite_count FROM %s AS f JOIN %s AS s ON f.to_user_id = s.user_id WHERE f.to_user_id = ?", tableNameA, tableNameB)
	query := fmt.Sprintf("SELECT s.user_id FROM %s AS f JOIN %s AS s ON f.user_id = s.user_id WHERE f.to_user_id = ?", tableNameA, tableNameB)
	err = db.conn.QueryRowsPartialCtx(ctx, &userList, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return userList, err
}

// QueryUsersOfFollowListByUserId query：userId ==> followList ==> followList.userId.fields.values
func (db *CustomDB) QueryUsersOfFollowListByUserId(ctx context.Context, userId int64) (userList []int64, err error) {
	tableNameA := "follow"
	tableNameB := "social"
	//query := fmt.Sprintf("SELECT s.id, s.user_id, s.follow_count, s.follower_count, s.total_favorited, s.work_count, s.favorite_count FROM %s AS f JOIN %s AS s ON f.to_user_id = s.user_id WHERE f.user_id = ?", tableNameA, tableNameB)
	query := fmt.Sprintf("SELECT s.user_id FROM %s AS f JOIN %s AS s ON f.to_user_id = s.user_id WHERE f.user_id = ?", tableNameA, tableNameB)
	err = db.conn.QueryRowsPartialCtx(ctx, &userList, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return userList, err
}

// QueryUserIdExistsInSocial query：userId ==> exist
func (db *CustomDB) QueryUserIdExistsInSocial(ctx context.Context, userId int64) (exists bool, err error) {
	tableName := "social"
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id = ?)", tableName)
	err = db.conn.QueryRowPartialCtx(ctx, &exists, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
		return false, err
	}
	return exists, err
}

// QueryFieldByUserIdInSocial query：userId in table
func (db *CustomDB) QueryFieldByUserIdInSocial(ctx context.Context, userId int64, fieldName string) (social Social, err error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ?", fieldName, "social")
	err = db.conn.QueryRowPartialCtx(ctx, &social, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId, fieldName)
		return social, err
	}
	return social, err
}

// QueryRecordByUserIdAndToUserIdInFollow query：userId & followId & status ==> record
func (db *CustomDB) QueryRecordByUserIdAndToUserIdInFollow(ctx context.Context, userId, toUserId int64) (follow Follow, err error) {
	tableName := "follow"
	query := fmt.Sprintf("SELECT id,user_id,to_user_id,status FROM %s WHERE user_id = ? AND to_user_id = ?", tableName)
	err = db.conn.QueryRowPartialCtx(ctx, &follow, query, userId, toUserId)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
		return follow, err
	}
	return follow, err
}

// InsertRecordByUserIdAndToUserIdInFollow insert： user record with toUser in follow
func (db *CustomDB) InsertRecordByUserIdAndToUserIdInFollow(ctx context.Context, userId, toUserId int64) error {
	var follow Follow
	tableName := "follow"
	insert := fmt.Sprintf("insert into %s (`user_id`, `to_user_id`, `status`) VALUES (?, ?, ?)", tableName)
	err := db.conn.QueryRowPartialCtx(ctx, &follow, insert, userId, toUserId, 0)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
		return err
	}
	return err
}

// InsertRecordByUserIdAndToUserIdInFriend insert： user record with toUser in friend
func (db *CustomDB) InsertRecordByUserIdAndToUserIdInFriend(ctx context.Context, userId, toUserId int64) error {
	var friend Friend
	tableName := "follow"
	insert := fmt.Sprintf("INSERT INTO %s (`user_id`, `to_user_id`, `status`) VALUES (?, ?, ?)", tableName)
	err := db.conn.QueryRowPartialCtx(ctx, &friend, insert, userId, toUserId, 1)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
		return err
	}
	return err
}

// DeleteRecordByUserIdAndToUserIdInFriend delete： user record with toUser in friend
func (db *CustomDB) DeleteRecordByUserIdAndToUserIdInFriend(ctx context.Context, userId int64, toUserId int64) error {
	var friend Friend
	tableName := "friend"
	insert := fmt.Sprintf("DELETE FROM %s WHERE to_user_id IN (?, ?) AND `user_id` IN (?, ?);", tableName)
	err := db.conn.QueryRowPartialCtx(ctx, &friend, insert, userId, toUserId, userId, toUserId)
	if err != nil {
		logc.Error(ctx, err, userId)
		return err
	}
	return err
}

// UpdateStatusByUserIdAndToUserIdInFollow update： user.status with toUser
func (db *CustomDB) UpdateStatusByUserIdAndToUserIdInFollow(ctx context.Context, userId int64, toUserId int64, actionType byte) error {
	var follow Follow
	tableName := "follow"
	update := fmt.Sprintf("UPDATE %s SET status = ? WHERE user_id = ? AND to_user_id = ?", tableName)
	err := db.conn.QueryRowPartialCtx(ctx, &follow, update, actionType, userId, toUserId)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
		return err
	}
	return err
}

// AutoIncrementUpdateFieldByUserIdAndToUserIdInTable update： userId.field in table
func (db *CustomDB) AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(ctx context.Context, userId int64, tableName string, fieldName string, value string) error {
	var obj interface{}
	update := fmt.Sprintf("UPDATE %s SET %s = %s + ? WHERE user_id = ?", tableName, fieldName, fieldName)
	err := db.conn.QueryRowPartialCtx(ctx, &obj, update, value, userId)
	if err != nil {
		logc.Error(ctx, err, fieldName, value, userId)
		return err
	}
	return err
}
