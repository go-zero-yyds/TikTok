package video

import (
	"TikTok/apps/video/rpc/model"
	"TikTok/apps/video/rpc/video"
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestDetail(t *testing.T) {
	videoId := 123225354801
	req := &video.BasicVideoInfoReq{
		VideoId: int64(videoId),
	}

	res, err := client.Detail(context.Background(), req)
	if errors.Is(err, status.New(codes.Unknown, model.ErrVideoNotFound.Error()).Err()) {
		log.Printf("videoID:%d Video Not Found\n", videoId)
		return
	}

	if err != nil {
		log.Println("Failed to  Detail service: ", err)
		return
	}

	testRes := &video.BasicVideoInfoResp{
		Video: &video.BasicVideoInfo{
			Id:       int64(videoId),
			UserId:   12345678,
			PlayUrl:  "12345678",
			CoverUrl: "testCurlfornewsnow2",
			Title:    "testtitel",
		},
	}
	log.Println(res.Video)
	log.Println(res)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, proto.Equal(res, testRes))
}
