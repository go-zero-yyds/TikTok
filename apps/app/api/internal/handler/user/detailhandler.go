package user

import (
	"TikTok/apps/app/api/apiVars"
	"net/http"

	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req)
		if err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, types.RespStatus(apiVars.InternalError))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
