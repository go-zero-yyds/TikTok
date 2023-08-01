package video

import (
	"TikTok/apps/video/rpc/video"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPublishList(t *testing.T) {

	req := &video.PublishListReq{
		UserId: 123456,
	}

	res, err := client.GetPublishList(context.Background(), req)
	if err != nil {
		log.Fatalln("err :", err)
	}

	require.NotNil(t, res)
	log.Println(res)
}
