package interaction

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/interactionclient"
	"TikTok/apps/video/rpc/model"
	"TikTok/apps/video/rpc/video"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionRequest) (resp *types.CommentActionResponse, err error) {

	// 参数检查
	if req.ActionType == 1 && (req.CommentText == "" || len(req.CommentText) > 500) { //如为评论则校验评论是否规范
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apivars.ErrTextRuleError),
		}, nil
	}

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)

	_, err = l.svcCtx.VideoRPC.Detail(l.ctx, &video.BasicVideoInfoReq{VideoId: req.VideoID})
	if errors.Is(err, model.ErrVideoNotFound) {
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apivars.ErrVideoNotFound),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	rpcReq := &interaction.CommentActionReq{
		UserId:     tokenID,
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
	}
	if req.CommentID == 0 {
		IPAddr := l.ctx.Value(middleware.IPKey).(string)
		IPAttr, err := l.svcCtx.GeoIPResolver.ResolveIP(IPAddr)
		if err != nil {
			return nil, err
		}
		rpcReq.IPAddr = &IPAddr
		rpcReq.IPAttr = &IPAttr
		rpcReq.CommentText = &req.CommentText
	} else {
		rpcReq.CommentId = &req.CommentID
	}
	sendCommentAction, err := l.svcCtx.InteractionRPC.SendCommentAction(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	apiResp := apivars.Success

	var commentInfo *types.Comment
	if req.CommentID == 0 {
		commentInfo, err = GetCommentInfo(sendCommentAction.Comment, tokenID, l.svcCtx, l.ctx)

		if err == apivars.ErrSomeData {
			return &types.CommentActionResponse{
				RespStatus: types.RespStatus(apiResp),
				Comment:    *commentInfo,
			}, nil
		}

		if err != nil && err != apivars.ErrSomeData {
			return nil, err
		}

		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiResp),
			Comment:    *commentInfo,
		}, nil
	}
	return &types.CommentActionResponse{
		RespStatus: types.RespStatus(apiResp),
	}, nil

}

func GetCommentInfo(comment *interactionclient.Comment, tokenID int64, svcCtx *svc.ServiceContext, ctx context.Context) (*types.Comment, error) {
	res := &types.Comment{}
	timestamp := comment.CreateDate
	if comment.Location != "" {
		res.CreateDate = fmt.Sprintf("%s · IP 属地%s", FormatTimestamp(timestamp), comment.Location)
	} else {
		res.CreateDate = FormatTimestamp(timestamp)
	}
	res.ID = comment.Id
	res.Content = comment.Content
	userInfo, err := user.TryGetUserInfo(tokenID, comment.UserId, svcCtx, ctx)
	if err != nil && err != apivars.ErrSomeData {
		return nil, err
	}
	res.User = *userInfo
	if err == apivars.ErrSomeData {
		return res, apivars.ErrSomeData
	}

	return res, nil
}

func FormatTimestamp(timestamp int64) string {
	currentTime := time.Now()
	timestampTime := time.UnixMilli(timestamp)
	year, month, day := currentTime.Year(), currentTime.Month(), currentTime.Day()
	timestampYear, timestampMonth, timestampDay := timestampTime.Year(), timestampTime.Month(), timestampTime.Day()

	if year == timestampYear && month == timestampMonth && day == timestampDay {
		diff := currentTime.Sub(timestampTime)
		if diff.Hours() >= 1 {
			return fmt.Sprintf("%.0f 小时前", diff.Hours())
		} else if diff.Minutes() >= 1 {
			return fmt.Sprintf("%.0f 分钟前", diff.Minutes())
		} else if diff.Seconds() >= 1 {
			return fmt.Sprintf("%.0f 秒钟前", diff.Seconds())
		} else {
			return "刚刚"
		}
	} else if year == timestampYear && month == timestampMonth && day-1 == timestampDay {
		return fmt.Sprintf("昨天 %02d:%02d", timestampTime.Hour(), timestampTime.Minute())
	} else if year == timestampYear && month == timestampMonth && day-2 == timestampDay {
		return fmt.Sprintf("前天 %02d:%02d", timestampTime.Hour(), timestampTime.Minute())
	} else if year == timestampYear {
		return fmt.Sprintf("%02d-%02d", timestampMonth, timestampDay)
	} else {
		return fmt.Sprintf("%d-%02d-%02d", timestampYear, timestampMonth, timestampDay)
	}
}
