package video

import (
	"TikTok/apps/video/rpc/video"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestDetail(t *testing.T) {
	videoId := 1686274172445003776
	req := &video.BasicVideoInfoReq{
		VideoId: int64(videoId),
	}

	res, err := client.Detail(context.Background(), req)
	if err != nil {
		fmt.Println("Failed to  Detail service: ", err)
	}

	testRes := &video.BasicVideoInfoResp{
		Video: &video.BasicVideoInfo{
			Id:       int64(videoId),
			UserId:   123456,
			PlayUrl:  "testPurl",
			CoverUrl: "testCurl",
			Title:    "啦啦啦",
		},
	}
	log.Println(res.Video)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, proto.Equal(res, testRes))
}
