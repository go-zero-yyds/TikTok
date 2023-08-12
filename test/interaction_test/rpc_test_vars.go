package interactiontest

import (
	"TikTok/apps/interaction/rpc/interactionclient"

	"github.com/zeromicro/go-zero/zrpc"
)

var client zrpc.Client
var logic interactionclient.Interaction

func init() {
	client = zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:8003",
	})
	// 创建服务端实例，传入  客户端
	logic = interactionclient.NewInteraction(client)
}
