package video

import (
	"TikTok/apps/video/rpc/video"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFeed(t *testing.T) {

	//lastTime := time.Now().Unix()
	lastTime2 := int64(3093515672355000)
	req := &video.FeedReq{
		LatestTime: &lastTime2,
	}

	res, err := client.GetFeed(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println(res)
}
