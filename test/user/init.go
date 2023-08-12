package user

import (
	"TikTok/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

var client zrpc.Client
var logic userclient.User

func init() {
	client = zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:8081",
	})
	// 创建服务端实例，传入  客户端
	logic = userclient.NewUser(client)
}
