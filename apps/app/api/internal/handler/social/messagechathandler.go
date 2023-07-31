package social

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"rpc/apps/app/api/internal/logic/social"
	"rpc/apps/app/api/internal/svc"
	"rpc/apps/app/api/internal/types"
)

func MessageChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := social.NewMessageChatLogic(r.Context(), svcCtx)
		resp, err := l.MessageChat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
