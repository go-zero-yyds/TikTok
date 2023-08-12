package video

import (
	"TikTok/apps/video/rpc/video"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestSend(t *testing.T) {

	req := &video.PublishActionReq{
		UserId:   12345678,
		PlayUrl:  "12345678",
		CoverUrl: "testCurlfornewsnow2",
		Title:    "testtitel",
	}

	res, err := client.SendPublishAction(context.Background(), req)
	if err != nil {
		fmt.Println("Failed to call login service: ", err)
	}

	testRes := &video.PublishActionResp{
		IsSucceed: true,
	}

	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, proto.Equal(res, testRes))
}
