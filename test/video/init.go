package video

import (
	"TikTok/apps/video/rpc/video"
	"flag"
	"fmt"

	"google.golang.org/grpc"
)

var client video.VideoClient
var conn grpc.ClientConn
var configFile = flag.String("f", "apps/video/rpc/etc/video.yaml", "the config file")

func init() {

	conn, err := grpc.Dial("127.0.0.1:8002", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect rpc server: " + err.Error())
	}

	client = video.NewVideoClient(conn)
}
