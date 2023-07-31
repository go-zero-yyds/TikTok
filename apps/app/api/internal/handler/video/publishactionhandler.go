package video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"rpc/apps/app/api/internal/logic/video"
	"rpc/apps/app/api/internal/svc"
	"rpc/apps/app/api/internal/types"
)

func PublishActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishActionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewPublishActionLogic(r.Context(), svcCtx)
		resp, err := l.PublishAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
