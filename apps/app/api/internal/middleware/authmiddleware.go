package middleware

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/app/api/utils/auth"
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type AuthMiddleware struct {
	JwtAuth auth.JwtAuth
}

func NewAuthMiddleware(jwtAuth auth.JwtAuth) *AuthMiddleware {
	return &AuthMiddleware{
		JwtAuth: jwtAuth,
	}
}

const TokenIDKey = "TokenIDKey"

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if token = r.URL.Query().Get("token"); token == "" {
			token = r.PostFormValue("token")
		}
		if token == "" {
			httpx.OkJsonCtx(r.Context(), w, types.RespStatus(apivars.NotLogged))
			return
		}
		tokenID, err := m.JwtAuth.ParseToken(token)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.RespStatus(apivars.TokenSignatureInvalid))
			return
		}
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, TokenIDKey, tokenID)
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}
