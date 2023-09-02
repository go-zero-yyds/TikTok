package robot

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-queue/kq"
)

// 添加其他机器人时要写run 和 DisplayPrologue 函数
//    在NewBossRobot中添加映射key就是bots的userid


type (
	Temp interface {
		//用户注册后机器人首先发送的第一条消息
		DisplayPrologue() string
		//处理信息的函数
		run(ctx context.Context, userId int64, toUserId int64, content string, v... any ) (string, error) 
	}

	BossRobot struct {
		Robots map[int64]Temp //存储机器人们的切片
	}
)

// ！每增加一个bot需在这里添加信息
func NewBossRobot(KqPusherClient *kq.Pusher) *BossRobot {
	ret := &BossRobot{}
	ret.Robots = make(map[int64]Temp, 1)
	k , v := NewSetPersonInfoRobot(KqPusherClient)
	ret.Robots[k] = v
	return ret
}

// 对外接口
// v 一般是 kpusher(queue)    FS(oss)
//返回值：
//		bool代表机器人是否执行操作
//		string代表机器人回发的消息
//		error代表出错
func (b *BossRobot) ProcessIfMessageForRobot(ctx context.Context, userId int64, toUserId int64, content string, v... any) (bool, string, error) {
	flag, err := b.isRobot(ctx, toUserId)
	if !flag || err != nil {
		return false, "", nil
	}
	num, err := b.parseRobot(ctx, toUserId)
	if err != nil {
		return false, "", nil
	}
	fmt.Println(len(v))
	data, err := b.Robots[num].run(ctx, userId, toUserId, content, v...)
	if err != nil {
		return false, "", nil
	}
	return true, data, nil
}

//下面俩个函数目前状况功能相似，现在留有接口，后续可以方便更改

// 查询是否是给bot发的信息
// 后续优化
func (b *BossRobot) isRobot(ctx context.Context, toUserId int64) (bool, error) {
	_ , exists := b.Robots[toUserId]
	return exists , nil
}

// 解析是给哪个机器人发送的信息，调用对应函数进行设置
func (b *BossRobot) parseRobot(ctx context.Context, toUserId int64) (int64, error) {
	return toUserId, nil
}
