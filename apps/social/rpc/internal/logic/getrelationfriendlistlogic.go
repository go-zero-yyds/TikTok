package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/model"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type GetRelationFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFriendListLogic {
	return &GetRelationFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRelationFriendList 获取好友列表（只要双方互相关注了自动变成好友）
func (l *GetRelationFriendListLogic) GetRelationFriendList(in *social.RelationFriendListReq) (*social.RelationFriendListResp, error) {

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId)
	}

	friendIdList, err := l.svcCtx.CustomDB.QueryFriendIdListByUserIdInFriend(l.ctx, in.UserId)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	messageList, err := l.svcCtx.CustomDB.QueryMessageByUserIdAndUserListInMessage(l.ctx, in.UserId, friendIdList)
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFriendListResp{UserList: nil}, nil
	}

	FriendUserList := make([]*social.FriendUser, len(friendIdList))

	respChan := make(chan []*social.FriendUser, 100)
	resultChan := make(chan *social.FriendUser, len(friendIdList))

	var wg sync.WaitGroup

	for _, v := range messageList {
		wg.Add(1)

		go func(v model.Message) {
			defer wg.Done()

			friendUser := &social.FriendUser{}
			if v.FromUserId != in.UserId {
				friendUser.UserId = v.FromUserId
				friendUser.MsgType = 0
			} else {
				friendUser.UserId = v.ToUserId
				friendUser.MsgType = 1
			}

			friendUser.Message = v.Content

			// 通过通道发送处理后的结果
			resultChan <- friendUser
		}(v)
	}

	// 启动一个goroutine从通道接收结果并更新FriendUserList
	go func() {
		wg.Wait()
		close(resultChan)

		for friendUser := range resultChan {
			FriendUserList = append(FriendUserList, friendUser)
		}

		respChan <- FriendUserList[len(friendIdList):]
	}()

	return &social.RelationFriendListResp{UserList: <-respChan}, nil
}
