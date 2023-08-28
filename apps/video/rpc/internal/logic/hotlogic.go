package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"TikTok/apps/video/rpc/internal/svc"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 获取视频得分窗口期的key
	// 以分钟为单位，redis设置hashtable
	SCORE_WINDOW_KEY = "video:hot:hash:window:%s"
	// 热门视频的key
	HOT_VIDEO_KEY = "video:hot:hash"
	// 热门视频排行榜的key
	SCORE_RANK_KEY = "video:hot:zset:score"
	// 视频最近更新事件的集合key
	SCORE_LAST_UPDATED_KEY = "video:hot:set:scorelastupdated"
	// 视频最近更新事件的集合（任务计算临时使用）key
	SCORE_LAST_UPDATED_TEMP_KEY = "video:hot:set:scorelastupdated:temp"
	// 热门视频（临时）存储集合的key
	HOT_VIDEO_TEMP_KEY = "video:hot:hash:temp"
	// 热门视频（临时） 排行榜的key
	SCORE_RANK_TEMP_KEY = "video:hot:zset:score:temp"
	// 任务锁的key,防止多个任务同时执行
	SCORE_RANK_TASK_LOCK_KEY = "video:hot:string:task:lock"
)

type VideoScore struct {
	VideoId int64 `json:"video_id" redis:"video_id"`
	Score   int64 `json:"score" redis:"score"`
}

type ScoreStorage struct {
	logx.Logger
	RedisClient *redis.Client
}

type HotVideoLogic struct {
	logx.Logger
	Storage      *ScoreStorage // 缓存
	ExpireMinute time.Duration // window key 过期时间
	HotScore     int           // 热门视频的得分标准
	Windowsize   int           // 计算多少分钟的累计热度值
	ctx          context.Context
}

// 热门视频
func NewHotVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HotVideoLogic {
	return &HotVideoLogic{
		HotScore:     svcCtx.Config.HotVideoConf.HotScore,
		Windowsize:   svcCtx.Config.HotVideoConf.Windowsize,
		ExpireMinute: time.Duration(svcCtx.Config.HotVideoConf.ExpireMinute) * time.Minute,
		Storage:      NewScoreStorage(svcCtx.RedisClient, logx.WithContext(ctx)),
		ctx:          ctx,
		Logger:       logx.WithContext(ctx),
	}
}

// 是否热门视频
func (p *HotVideoLogic) IsHotVideo(VideoId int64) (bool, error) {
	return p.Storage.IsHotVideo(p.ctx, strconv.Itoa(int(VideoId)))
}

// 计算热门视频，定时任务
func (p *HotVideoLogic) ScoreCalculation(ctime time.Time) error {
	p.Logger.Infof("received task:%s", ctime)
	start := time.Now()
	defer func() {
		p.Logger.Infof("ScoreCalculation cost:%s", time.Since(start))
	}()

	err := p.Storage.ScoreCalculation(p.ctx, ctime, p.Windowsize, p.HotScore)
	if err != nil {
		p.Logger.Errorf("ScoreCalculation failed %s", err.Error())
		return err
	}
	return nil
}

// 视频得分累计
// ctime: 事件发生时间,用于得分窗口期的计算
func (p *HotVideoLogic) ScoresIncr(videosScore []VideoScore, ctime time.Time) error {
	err := p.Storage.ScoresIncr(p.ctx, videosScore, p.ExpireMinute, ctime)
	if err != nil {
		p.Logger.Errorf("SetScore VideosScoreIncr failed %s videosScore:%v", err.Error(), videosScore)
		return err
	}
	return nil
}

// 清除所有缓存,用于测试
func (p *HotVideoLogic) CleanAll(ctime time.Time) {
	p.Storage.CleanAll(p.ctx, ctime, p.Windowsize)
}

func NewScoreStorage(RedisClient *redis.Client, logger logx.Logger) *ScoreStorage {
	return &ScoreStorage{RedisClient: RedisClient, Logger: logger}
}

// 获取视频得分窗口期的key
// 以分钟为单位，redis设置hashtable
func (s *ScoreStorage) GetScoreWindowKey(ctime time.Time) string {
	return fmt.Sprintf(SCORE_WINDOW_KEY, ctime.Format("200601021504"))
}

// 热门视频的key
func (s *ScoreStorage) GetHotVideoKey() string {
	return HOT_VIDEO_KEY
}

// 热门视频排行榜的key
func (s *ScoreStorage) GetScoreRankKey() string {
	return SCORE_RANK_KEY
}

// 视频最近更新事件的集合key
func (s *ScoreStorage) GetScoreLastUpdatedKey() string {
	return SCORE_LAST_UPDATED_KEY
}

// 热门视频（临时）存储集合的key
func (s *ScoreStorage) GetHotVideoTempKey() string {
	return HOT_VIDEO_TEMP_KEY
}

// 热门视频（临时） 排行榜的key
func (s *ScoreStorage) GetScoreRankTempKey() string {
	return SCORE_RANK_TEMP_KEY
}

// 任务锁的key,防止多个任务同时执行
// func (s *ScoreStorage) GetScoreRankTaskLockKey() string {
// 	return SCORE_RANK_TASK_LOCK_KEY
// }

// 视频最近更新事件的集合（任务计算临时使用）
func (s *ScoreStorage) GetScoreLastUpdatedTempKey() string {
	return SCORE_LAST_UPDATED_TEMP_KEY
}

// 视频得分累计
// ctime: 事件发生时间,用于得分窗口期的计算
// expiration: 过期时间
func (s *ScoreStorage) ScoresIncr(ctx context.Context, videos []VideoScore, expiration time.Duration, ctime time.Time) error {
	pipe := s.RedisClient.Pipeline()
	for _, video := range videos {
		// 窗口期视频得分累计
		pipe.HIncrBy(ctx, s.GetScoreWindowKey(ctime), strconv.Itoa(int(video.VideoId)), video.Score)
		// 记录视频最近更新事件，方便定时任务计算视频得分
		pipe.SAdd(ctx, s.GetScoreLastUpdatedKey(), video.VideoId)
	}
	// 设置key过期时间
	pipe.Expire(ctx, s.GetScoreWindowKey(ctime), expiration)
	_, err := pipe.Exec(ctx)
	return err
}

// 获取视频得分窗口期的keys
// 以分钟为单位，redis设置hashtable
// Windowsize: 最近几分钟
// ctime: 事件发生时间,用于得分窗口期的计算
func (s *ScoreStorage) GetKeysWithCtimeWindowsize(ctime time.Time, Windowsize int) []string {
	keys := make([]string, 0, Windowsize)
	for i := 0; i < Windowsize; i++ {
		keys = append(keys, s.GetScoreWindowKey(ctime))
		ctime = ctime.Add(-time.Minute)
	}
	return keys
}

// 计算热门视频
func (s *ScoreStorage) ScoreCalculation(ctx context.Context, ctime time.Time, Windowsize int, pv int) error {
	// 任务计算的时候会把SCORE_LAST_UPDATED_KEY重命名为SCORE_LAST_UPDATED_TEMP_KEY
	// 保证任务计算的时候，不会有新的视频事件加入,以后可以考虑使用更优方法
	// 在集群模式下， 和key必须newkey位于同一个哈希槽中，
	// 这意味着实际上只有具有相同哈希标签的键才能在集群中可靠地重命名。
	// https://redis.io/commands/rename
	if _, err := s.RedisClient.Rename(ctx, s.GetScoreLastUpdatedKey(), s.GetScoreLastUpdatedTempKey()).Result(); err != nil {
		return err
	}

	// 获取视频得分窗口期的keys
	keys := s.GetKeysWithCtimeWindowsize(ctime, Windowsize)
	updated := false
	for {
		// 视频最近有接收新事件的集合，弹出一个待计算的视频ID,从集合中随机获取一个元素 set
		videoId, err := s.RedisClient.SPop(ctx, s.GetScoreLastUpdatedTempKey()).Result()
		if err != nil {
			if err == redis.Nil {
				break
			}
			return err
		}

		// 批量获取视频所有窗口得分 hash
		pipe := s.RedisClient.Pipeline()
		for _, key := range keys {
			pipe.HGet(ctx, key, videoId)
		}
		values, err := pipe.Exec(ctx)
		if err != nil && err != redis.Nil {
			return err
		}
		// 视频得分累计
		totalScore := 0
		for _, value := range values {
			if value == nil {
				continue
			}
			score, err := value.(*redis.StringCmd).Int()
			if err != nil && err != redis.Nil {
				return err
			}
			totalScore += score
		}

		// 总得分是否满足热门视频的标准
		if totalScore < pv {
			s.Logger.Infof("ScoreCalculation videoId:%s is not hot video", videoId)
			continue
		}

		pipe = s.RedisClient.Pipeline()
		// 写入热门视频临时排行榜 有序集合Sorted Set
		pipe.ZAdd(ctx,
			s.GetScoreRankTempKey(),
			redis.Z{Score: float64(totalScore), Member: videoId},
		)
		// 写入热门视频临时存储集合 hash
		pipe.HSet(ctx, s.GetHotVideoTempKey(), videoId, totalScore)
		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}
		s.Logger.Infof("ScoreCalculation videoId:%s is hot video", videoId)
		updated = true
	}

	if !updated {
		return nil
	}

	pipe := s.RedisClient.Pipeline()
	// 把临时排行榜交换到正式排行榜
	pipe.Rename(ctx, s.GetScoreRankTempKey(), s.GetScoreRankKey())
	// 把临时存储集合交换到正式热门集合
	pipe.Rename(ctx, s.GetHotVideoTempKey(), s.GetHotVideoKey())
	// 最后保留热门排行榜单 top200
	pipe.ZRemRangeByRank(ctx, s.GetScoreRankKey(), 0, -201)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// 是否热门视频
func (s *ScoreStorage) IsHotVideo(ctx context.Context, videoId string) (bool, error) {
	// Redis查询哈希表
	return s.RedisClient.HExists(ctx, s.GetHotVideoKey(), videoId).Result()
}

// 清除所有缓存,用于测试
func (s *ScoreStorage) CleanAll(ctx context.Context, ctime time.Time, windowsize int) {
	keys := []string{
		s.GetScoreRankKey(),
		s.GetHotVideoKey(),
		s.GetScoreLastUpdatedKey(),
		s.GetScoreRankTempKey(),
		s.GetHotVideoTempKey(),
		s.GetScoreLastUpdatedTempKey(),
	}
	for i := 0; i < windowsize; i++ {
		keys = append(keys, s.GetScoreWindowKey(ctime.Add(-time.Minute*time.Duration(i))))
	}
	pipe := s.RedisClient.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	pipe.Exec(ctx)
}
