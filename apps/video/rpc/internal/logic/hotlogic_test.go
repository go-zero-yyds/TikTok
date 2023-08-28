package logic

import (
	"context"
	"testing"
	"time"

	"TikTok/apps/video/rpc/internal/config"
	"TikTok/apps/video/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = "../../etc/video.yaml"

func SetupAndTeardown(t *testing.T) (*HotVideoLogic, func(t *testing.T, ctime time.Time)) {
	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx, _ := svc.NewServiceContext(c)
	ctx.Config.HotVideoConf.HotScore = 1000
	ctx.Config.HotVideoConf.Windowsize = 3
	hotVideoLogic := NewHotVideoLogic(context.TODO(), ctx)
	return hotVideoLogic, func(t *testing.T, ctime time.Time) {
		// return
		hotVideoLogic.CleanAll(ctime)
	}
}

func TestHotVideoLogic(t *testing.T) {
	ctx := context.Background()
	ctime, _ := time.Parse(
		"2006-01-02 15:04:05  -0700 MST",
		"2023-08-21 20:21:32 +0800 CST",
	)
	ctime2, _ := time.Parse(
		"2006-01-02 15:04:05  -0700 MST",
		"2023-08-21 20:18:32 +0800 CST",
	)
	vid1 := int64(123)
	vid2 := int64(456)
	vid3 := int64(789)
	hotVideoLogic, cleanup := SetupAndTeardown(t)
	defer cleanup(t, ctime)

	tests := []struct {
		name          string
		wantErr       error
		ctx           context.Context
		ctime         time.Time
		hotVideoLogic *HotVideoLogic
		videosScore   []VideoScore
	}{
		{
			name:          "ctime",
			ctime:         ctime,
			wantErr:       nil,
			ctx:           ctx,
			hotVideoLogic: hotVideoLogic,
			videosScore: []VideoScore{
				VideoScore{VideoId: vid1, Score: 1001},
				VideoScore{VideoId: vid2, Score: 999},
				VideoScore{VideoId: vid3, Score: 2000},
			},
		},
		{
			name:          "ctime2-vid:456",
			ctime:         ctime2,
			wantErr:       nil,
			ctx:           ctx,
			hotVideoLogic: hotVideoLogic,
			videosScore: []VideoScore{
				VideoScore{VideoId: vid2, Score: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.hotVideoLogic.ScoresIncr(tt.videosScore, tt.ctime); err != tt.wantErr {
				t.Errorf("ScoresIncr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := hotVideoLogic.ScoreCalculation(ctime); err != nil {
		t.Errorf("ScoreCalculation() error = %v", err)
	}

	TestsIshotVideo := []struct {
		name          string
		wantErr       error
		ctx           context.Context
		hotVideoLogic *HotVideoLogic
		VideoId       int64
		want          bool
	}{
		{
			name:          "vid:123",
			wantErr:       nil,
			ctx:           ctx,
			hotVideoLogic: hotVideoLogic,
			VideoId:       vid1,
			want:          true,
		},
		{
			name:          "vid:456",
			wantErr:       nil,
			ctx:           ctx,
			hotVideoLogic: hotVideoLogic,
			VideoId:       vid2,
		
			want: false,
		},
	}

	for _, tt := range TestsIshotVideo {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := tt.hotVideoLogic.IsHotVideo(tt.VideoId)
			if err != tt.wantErr {
				t.Errorf("IshotVideo() VideoId:%v error = %v, wantErr %v", tt.VideoId, err, tt.wantErr)
			}
			if ret != tt.want {
				t.Errorf("IshotVideo() VideoId:%v ret = %v, want %v", tt.VideoId, ret, tt.want)
			}
		})
	}
}
