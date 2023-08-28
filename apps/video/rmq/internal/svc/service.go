package svc

import (
	"TikTok/apps/video/rmq/internal/config"
	"TikTok/apps/video/rpc/videoclient"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	chanCount   = 2
	bufferCount = 3
)

type Service struct {
	Config   config.Config
	VideoRPC videoclient.Video

	msgsChan []chan *KafkaData
	waiter   sync.WaitGroup
	ctx      context.Context
}

type KafkaData struct {
	VideoId int64 `json:"video_id"`
}

func NewService(c config.Config) *Service {
	s := &Service{
		Config:   c,
		msgsChan: make([]chan *KafkaData, chanCount),
		VideoRPC: videoclient.NewVideo(zrpc.MustNewClient(c.VideoRPC)),
		ctx:      context.Background(),
	}
	for i := 0; i < chanCount; i++ {
		ch := make(chan *KafkaData, bufferCount)
		s.msgsChan[i] = ch
		s.waiter.Add(1)
		go s.consumeDTM(ch)
	}
	return s
}

func (s *Service) consumeDTM(ch chan *KafkaData) {
	defer s.waiter.Done()
	for {
		msg, ok := <-ch
		if !ok {
			log.Fatal("video rmq exit")
		}

		//校验视频是否热门
		resp, err := s.VideoRPC.CheckHotVideo(s.ctx,
			&videoclient.CheckHotVideoReq{
				VideoId: msg.VideoId,
			},
		)
		if err != nil {
			logx.Errorf("RPC CheckHotVideo error: %v", err)
			return
			// continue
		}

		if resp.IsHotVideo == 1 {
			logx.Infof("video %d is hot", msg.VideoId)
		} else {
			logx.Infof("video %d is not hot", msg.VideoId)
		}
		//TODO:
	}
}

func (s *Service) Consume(_ string, value string) error {
	logx.Infof("Consume value: %s\n", value)
	var data []*KafkaData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	for _, d := range data {
		s.msgsChan[rand.Intn(len(s.msgsChan))] <- d
	}
	return nil
}
