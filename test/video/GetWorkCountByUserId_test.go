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

func TestGetWork(t *testing.T) {

	req := &video.WorkCountByUserIdReq{
		UserId: 123456,
	}

	res, err := client.GetWorkCountByUserId(context.Background(), req)
	if err != nil {
		fmt.Println("Failed to call login service: ", err)
	}

	testRes := &video.WorkCountByUserIdResp{
		WorkCount: 3,
	}

	require.NoError(t, err)
	require.NotNil(t, res)
	log.Println(res.WorkCount)
	require.True(t, proto.Equal(res, testRes))
}
