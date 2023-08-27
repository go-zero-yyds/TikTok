package interaction

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/interactionclient"
	"TikTok/apps/video/rpc/video"
	"context"
	"fmt"
	"strconv"
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
	// todo: add your logic here and delete this line
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.VideoRPC.Detail(l.ctx, &video.BasicVideoInfoReq{VideoId: req.VideoID})
	if err != nil {
		return nil, err
	}
	IPAddr := l.ctx.Value(middleware.IPKey).(string)
	IPAttr, err := l.svcCtx.GeoIPResolver.ResolveIP(IPAddr)
	if err != nil {
		return nil, err
	}
	rpcReq := &interaction.CommentActionReq{
		UserId:     tokenID,
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
		IPAddr:     &IPAddr,
		IPAttr:     &IPAttr,
	}
	if req.CommentID == 0 {
		rpcReq.CommentText = &req.CommentText
	} else {
		rpcReq.CommentId = &req.CommentID
	}
	sendCommentAction, err := l.svcCtx.InteractionRPC.SendCommentAction(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	apiResp := apiVars.Success

	var commentInfo *types.Comment
	if req.CommentID == 0 {
		commentInfo, err = GetCommentInfo(sendCommentAction.Comment, tokenID, l.svcCtx, l.ctx)

		if err == apiVars.SomeDataErr {
			return &types.CommentActionResponse{
				RespStatus: types.RespStatus(apiResp),
				Comment:    *commentInfo,
			}, nil
		}

		if err != nil && err != apiVars.SomeDataErr {
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
	timestamp, err := strconv.ParseInt(comment.CreateDate, 10, 64)
	if err != nil {
		return nil, err
	}
	if comment.Location != "" {
		res.CreateDate = fmt.Sprintf("%s · IP 属地%s", FormatTimestamp(timestamp), comment.Location)
	} else {
		res.CreateDate = FormatTimestamp(timestamp)
	}
	res.ID = comment.Id
	res.Content = comment.Content
	userInfo, err := user.TryGetUserInfo(tokenID, comment.UserId, svcCtx, ctx)
	if err != nil && err != apiVars.SomeDataErr {
		return nil, err
	}
	res.User = *userInfo
	if err == apiVars.SomeDataErr {
		return res, apiVars.SomeDataErr
	}

	return res, nil
}

func FormatTimestamp(timestamp int64) string {
	currentTime := time.Now()
	timestampTime := time.Unix(timestamp, 0)
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
