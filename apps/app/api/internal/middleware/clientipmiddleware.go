package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const IPKey = "IPaddr"

type ClientIPMiddleware struct {
}

func NewClientIPMiddleware() *ClientIPMiddleware {
	return &ClientIPMiddleware{}
}

func (m *ClientIPMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error parsing remote address:", err)
			return
		}
		fmt.Println("host=", host)
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, IPKey, host)
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}
