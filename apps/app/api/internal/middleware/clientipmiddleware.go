package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const IPKey = "IPAddr"

type ClientIPMiddleware struct {
}

func NewClientIPMiddleware() *ClientIPMiddleware {
	return &ClientIPMiddleware{}
}

func (m *ClientIPMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, IPKey, ClientIP(r))
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
