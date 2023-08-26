package interaction

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/interactionclient"
	"TikTok/apps/video/rpc/video"
	"context"
	"fmt"
	"regexp"
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

	// 参数检查
	matched, err := regexp.MatchString("^\\d+$", strconv.FormatInt(req.VideoID, 10)) //是否为纯数字
	if strconv.FormatInt(req.VideoID, 10) == "" || matched == false {
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiVars.VideoIdRuleError),
		}, nil
	} else if (req.ActionType != 1) && (req.ActionType != 2) { //是否有除1或2的数字
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiVars.ActionTypeRuleError),
		}, nil
	}

	matched, err = regexp.MatchString("^\\d+$", strconv.FormatInt(req.CommentID, 10)) //是否为纯数字
	if req.ActionType == 1 && req.CommentText == "" {                                 //如为评论则校验评论是否规范
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiVars.TextIsNull),
		}, nil
	} else if req.ActionType == 2 && matched == false { //如为删评则校验评论id是否规范
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiVars.CommentIdRuleError),
		}, nil
	}

	if req.Token == "" {
		return &types.CommentActionResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.VideoRPC.Detail(l.ctx, &video.BasicVideoInfoReq{VideoId: req.VideoID})
	if err != nil {
		return nil, err
	}
	rpcReq := &interaction.CommentActionReq{
		UserId:     tokenID,
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
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
	res.CreateDate = FormatTimestamp(timestamp)
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
			return fmt.Sprintf("%.0f小时前", diff.Hours())
		} else if diff.Minutes() >= 1 {
			return fmt.Sprintf("%.0f分钟前", diff.Minutes())
		} else {
			return fmt.Sprintf("%.0f秒钟前", diff.Seconds())
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
