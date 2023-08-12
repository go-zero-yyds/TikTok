/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/12 17:25
 */

package social

import (
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/social/rpc/socialclient"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	client zrpc.Client
	logic  social.SocialClient
)

func init() {
	client = zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:8004",
	})
	logic = socialclient.NewSocial(client)
}
