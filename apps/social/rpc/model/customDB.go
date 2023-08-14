package model

import (
	"TikTok/apps/social/rpc/social"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

type CustomDB struct {
	Message MessageModel
	Follow  FollowModel
	Social  SocialModel
	Conn    sqlx.SqlConn
}

// NewCustomDB 初始化数据库
func NewCustomDB(conn sqlx.SqlConn) *CustomDB {
	return &CustomDB{
		Message: NewMessageModel(conn),
		Follow:  NewFollowModel(conn),
		Social:  NewSocialModel(conn),
		Conn:    conn,
	}
}

// QueryMessageByUserIdAndToUserIdInMessage query：userId & toUserId & (<time) ==> message
func (db *CustomDB) QueryMessageByUserIdAndToUserIdInMessage(ctx context.Context, userId int64, toUserId int64, time string) (messageList []*social.Message, err error) {
	tableName := "message"
	query := fmt.Sprintf("SELECT id, from_user_id, to_user_id, content, created_time FROM %s WHERE ((from_user_id = ? AND to_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)) AND created_time <= ? ORDER BY created_time DESC;", tableName)
	err = db.Conn.QueryRowsPartialCtx(ctx, &messageList, query, userId, toUserId, userId, toUserId, time)
	if err != nil {
		logc.Error(ctx, err, messageList, query, userId, toUserId, userId, toUserId, time)
	}
	return messageList, err
}

// QueryMessageByUserIdAndUserListInMessage query：userId in (userList) ==> message
func (db *CustomDB) QueryMessageByUserIdAndUserListInMessage(ctx context.Context, userId int64, userList []int64) (messageList []Message, err error) {
	tableName := "message"
	userListString := IntListToString(userList)
	query := fmt.Sprintf("SELECT m.id, m.from_user_id, m.to_user_id, m.content, m.created_time FROM message AS m WHERE (m.from_user_id IN (%s) OR m.to_user_id IN (%s)) AND ((m.from_user_id = ? AND m.to_user_id IN (%s)) OR (m.to_user_id = ? AND m.from_user_id IN (%s))) AND m.created_time = (SELECT MAX(created_time) FROM %s WHERE (from_user_id = m.from_user_id AND to_user_id = m.to_user_id) OR (from_user_id = m.to_user_id AND to_user_id = m.from_user_id));", userListString, userListString, userListString, userListString, tableName)
	err = db.Conn.QueryRowsPartialCtx(ctx, &messageList, query, userId, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return messageList, err
}

// QueryFriendIdListByUserIdInFriend query：userId ==> friendId
func (db *CustomDB) QueryFriendIdListByUserIdInFriend(ctx context.Context, userId int64) (friendIdList []int64, err error) {
	tableName := "friend"
	query := fmt.Sprintf("SELECT CASE WHEN user_id = ? THEN to_user_id WHEN to_user_id = ? THEN user_id END AS friend_id FROM %s WHERE user_id = ? OR to_user_id = ? AND status = 1", tableName)
	err = db.Conn.QueryRowsPartialCtx(ctx, &friendIdList, query, userId, userId, userId, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return friendIdList, err
}

// QueryUsersOfFollowerListByUserId query：userId ==> followerList ==> followerList.userId.fields.values
func (db *CustomDB) QueryUsersOfFollowerListByUserId(ctx context.Context, userId int64) (userList []int64, err error) {
	tableNameA := "follow"
	tableNameB := "social"
	//query := fmt.Sprintf("SELECT s.id, s.user_id, s.follow_count, s.follower_count, s.total_favorited, s.work_count, s.favorite_count FROM %s AS f JOIN %s AS s ON f.to_user_id = s.user_id WHERE f.to_user_id = ?", tableNameA, tableNameB)
	query := fmt.Sprintf("SELECT s.user_id FROM %s AS f JOIN %s AS s ON f.user_id = s.user_id WHERE f.to_user_id = ?", tableNameA, tableNameB)
	err = db.Conn.QueryRowsPartialCtx(ctx, &userList, query, userId)
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
	err = db.Conn.QueryRowsPartialCtx(ctx, &userList, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	return userList, err
}

// QueryUserIdIsExistInSocial query：userId ==> exist
func (db *CustomDB) QueryUserIdIsExistInSocial(ctx context.Context, userId int64) (exist bool, err error) {
	tableName := "social"
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id = ?)", tableName)
	err = db.Conn.QueryRowPartialCtx(ctx, &exist, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId)
		return false, err
	}
	return exist, err
}

// QueryFieldByUserIdInSocial query：userId in table
func (db *CustomDB) QueryFieldByUserIdInSocial(ctx context.Context, userId int64, fieldName string) (social Social, err error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ?", fieldName, "social")
	err = db.Conn.QueryRowPartialCtx(ctx, &social, query, userId)
	if err != nil {
		logc.Error(ctx, err, userId, fieldName)
	}
	return social, err
}

// QueryIsExistByUserIdAndToUserIdInFriend query：userId & followId ==> exist
func (db *CustomDB) QueryIsExistByUserIdAndToUserIdInFriend(ctx context.Context, userId, toUserId int64) (exist bool, err error) {
	tableName := "friend"
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE (user_id = ? AND to_user_id = ?) OR (user_id = ? AND to_user_id = ?)) AS result", tableName)
	err = db.Conn.QueryRowPartialCtx(ctx, &exist, query, userId, toUserId, toUserId, userId)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
		return false, err
	}
	return exist, err
}

// QueryRecordByUserIdAndToUserIdAndStatusInFollow query：userId & followId & status by status ==> record
func (db *CustomDB) QueryRecordByUserIdAndToUserIdAndStatusInFollow(ctx context.Context, userId, toUserId int64, status int8) (follow Follow, err error) {
	tableName := "follow"
	query := fmt.Sprintf("SELECT id,user_id,to_user_id,status FROM %s WHERE user_id = ? AND to_user_id = ? AND status = ?", tableName)
	err = db.Conn.QueryRowPartialCtx(ctx, &follow, query, userId, toUserId, status)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
	}
	return follow, err
}

// QueryRecordByUserIdAndToUserIdInFollow query：userId & followId & status ==> record
func (db *CustomDB) QueryRecordByUserIdAndToUserIdInFollow(ctx context.Context, userId, toUserId int64) (follow Follow, err error) {
	tableName := "follow"
	query := fmt.Sprintf("SELECT id,user_id,to_user_id,status FROM %s WHERE user_id = ? AND to_user_id = ?", tableName)
	err = db.Conn.QueryRowPartialCtx(ctx, &follow, query, userId, toUserId)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
	}
	return follow, err
}

// InsertRecordByUserIdAndToUserIdInFollow insert： user record with toUser in follow
func (db *CustomDB) InsertRecordByUserIdAndToUserIdInFollow(ctx context.Context, userId, toUserId int64) error {
	var follow Follow
	tableName := "follow"
	insert := fmt.Sprintf("insert into %s (`user_id`, `to_user_id`, `status`) VALUES (?, ?, ?)", tableName)
	err := db.Conn.QueryRowPartialCtx(ctx, &follow, insert, userId, toUserId, 0)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// InsertRecordByUserIdAndToUserIdInFriend insert： user record with toUser in friend
func (db *CustomDB) InsertRecordByUserIdAndToUserIdInFriend(ctx context.Context, userId, toUserId int64, status int8) error {
	var friend Friend
	tableName := "friend"
	insert := fmt.Sprintf("INSERT INTO %s (`user_id`, `to_user_id`, `status`) VALUES (?, ?, ?)", tableName)
	err := db.Conn.QueryRowPartialCtx(ctx, &friend, insert, userId, toUserId, status)
	if err != nil {
		logc.Error(ctx, err, userId, toUserId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// InsertRecordByUserIdAndToUserIdAndContentInMessage insert： message record with toUser in message
func (db *CustomDB) InsertRecordByUserIdAndToUserIdAndContentInMessage(ctx context.Context, fromUserId, toUserId int64, content string) error {
	tableName := "message"
	var message Message
	currentTime := time.Now()
	timeStr := currentTime.Format("2006-01-02 15:04:05")
	insert := fmt.Sprintf("INSERT INTO %s (from_user_id, to_user_id, content, created_time) VALUES (?, ?, ?, ?)", tableName)
	err := db.Conn.QueryRowCtx(ctx, &message, insert, fromUserId, toUserId, content, timeStr)
	if err != nil {
		logc.Error(ctx, err, fromUserId, toUserId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// InsertRecordByUserIdInSocial insert： user in social
func (db *CustomDB) InsertRecordByUserIdInSocial(ctx context.Context, userId int64, followCount int64, followerCount int64) error {
	var social Social
	tableName := "social"
	insert := fmt.Sprintf("INSERT INTO %s (`user_id`, `follow_count`, `follower_count`) VALUES (?, ?, ?)", tableName)
	err := db.Conn.QueryRowPartialCtx(ctx, &social, insert, strconv.FormatInt(userId, 10), followCount, followerCount)
	if err != nil {
		logc.Error(ctx, err, userId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// UpdateRecordByUserIdAndToUserIdInFriend update： user record with toUser in friend
func (db *CustomDB) UpdateRecordByUserIdAndToUserIdInFriend(ctx context.Context, userId, toUserId int64, status int8) error {
	var friend Friend
	tableName := "friend"
	update := fmt.Sprintf("UPDATE %s SET status = ? WHERE to_user_id IN (?, ?) AND `user_id` IN (?, ?)", tableName)
	err := db.Conn.QueryRowPartialCtx(ctx, &friend, update, status, userId, toUserId, userId, toUserId)
	if err != nil && err != sql.ErrNoRows {
		logc.Error(ctx, err, userId, toUserId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// UpdateStatusByUserIdAndToUserIdInFollow update： user.status with toUser
func (db *CustomDB) UpdateStatusByUserIdAndToUserIdInFollow(ctx context.Context, userId int64, toUserId int64, actionType byte) error {
	var follow Follow
	tableName := "follow"
	update := fmt.Sprintf("UPDATE %s SET status = ? WHERE user_id = ? AND to_user_id = ?", tableName)
	err := db.Conn.QueryRowPartialCtx(ctx, &follow, update, actionType, userId, toUserId)
	if err != nil && err != sql.ErrNoRows {
		logc.Error(ctx, err, userId, toUserId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// AutoIncrementUpdateFieldByUserIdAndToUserIdInTable update： userId.field in table
func (db *CustomDB) AutoIncrementUpdateFieldByUserIdAndToUserIdInTable(ctx context.Context, userId int64, tableName string, fieldName string, value string) error {
	var obj interface{}
	update := fmt.Sprintf("UPDATE %s SET %s = %s + ? WHERE user_id = ?", tableName, fieldName, fieldName)
	err := db.Conn.QueryRowPartialCtx(ctx, &obj, update, value, userId)
	if err != nil && err != sql.ErrNoRows {
		logc.Error(ctx, err, fieldName, value, userId)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// IntListToString int列表转字符串
func IntListToString(list []int64) string {
	str := make([]string, len(list))
	for i, num := range list {
		str[i] = fmt.Sprintf("%d", num)
	}

	result := strings.Join(str, ", ")
	return result
}

func (db *CustomDB) TransactionUpdateSocialCount(ctx context.Context, session sqlx.Session, userId int64, toUserId int64, actionType byte, count int8) (err error) {

	//update status in follow
	var follow Follow
	update := fmt.Sprintf("UPDATE %s SET status = ? WHERE user_id = ? AND to_user_id = ?", "follow")
	err = session.QueryRowPartialCtx(ctx, &follow, update, actionType, userId, toUserId)
	if err != sql.ErrNoRows {
		return err
	}

	//update followCount in follow
	var obj interface{}
	update = fmt.Sprintf("UPDATE %s SET %s = %s + ? WHERE user_id = ?", "social", "follow_count", "follow_count")
	err = session.QueryRowPartialCtx(ctx, &obj, update, count, userId)
	if err != sql.ErrNoRows {
		return err
	}

	//update followerCount in follow
	update = fmt.Sprintf("UPDATE %s SET %s = %s + ? WHERE user_id = ?", "social", "follower_count", "follower_count")
	err = session.QueryRowPartialCtx(ctx, &obj, update, count, toUserId)
	if err != sql.ErrNoRows {
		return err
	}

	return err
}
