package robot

import (
	"TikTok/pkg/FileSystem"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logc"
)
var(
	ErrorKqPusher  = errors.New("not KqPusher type")
	ErrorGetAvatar = errors.New("get qq avatar error from http")
	ErrorFileOss   = errors.New("oss error")
	ErrorParam     = errors.New("param error")
)

// 功能是 接收用户信息参数，返回需要发送给user的数据包

// 格式是json包的map[string][]string   【uid】{修改字段，key}
//类型：avatar           头像
//	    backgroundimage 背景图
//      signature       个性签名
type SetPersonInfoRobot struct {
	prologue string
}

func NewSetPersonInfoRobot(KqPusherClient *kq.Pusher)(int64, *SetPersonInfoRobot){
	message := make(map[int64][]string)
	message[0] =[]string{"username" , "@抖音1号"}
	data , err := json.Marshal(message)
	if err != nil{
		panic("robots start error")
	}
	err = KqPusherClient.Push(string(data))
	if err != nil{
		panic("robots start error")
	}
	return 0 , &SetPersonInfoRobot{prologue: "你好呀,我是一个可以修改头像、背景大图、个性签名的机器人呀,请输入help查看命令。"}
}

func (s *SetPersonInfoRobot) DisplayPrologue() string{
	return s.prologue
}

func (s *SetPersonInfoRobot) run(ctx context.Context, userId int64, toUserId int64, content string, v... any) (string ,error){
	data, err := s.deal(ctx, userId, content, v...)
	if err != nil {
		return "" , err
	}
	return data , nil
}

// 返回发送给userid的信息
func (s *SetPersonInfoRobot) deal(ctx context.Context, userId int64, content string,v... any) (string, error) {
	ret := false
	if  len(content) >= 3 && strings.ToLower(content[:3]) == "set" {
		infoMap := s.parse(content)
		//修改头像提取qq号 支持qq
		if value, ok := infoMap["avatar"]; ok {
			if value[0:2] == "qq" && len(v) >= 2 {
				err := s.ToSetAvatar(ctx, userId, value[2:] , v...)
				if err != nil {
					logc.Error(ctx, userId, err)
					return "设置错误,请检查参数是否正确，或等会再试", nil
				}
				ret = true
			}
		}

		//设置背景大图 支持qq
		if value, ok := infoMap["backgroundimage"]; ok {
			if value[0:3] == "url" {
				err := s.ToSetBackgroundImage(ctx, userId, value[3:] , v...)
				if err != nil {
					logc.Error(ctx, userId, err)
					return "设置错误,请检查参数是否正确，或等会再试", nil
				}
				ret = true
			}
		}

		//修改个性签名
		if value, ok := infoMap["signature"]; ok {
			err := s.ToSetSignature(ctx, userId, value , v...)
			if err != nil {
				logc.Error(ctx, userId, err)
				return "设置错误,请检查参数是否正确，或等会再试", nil
			}
			ret = true
		}
	}
	if len(content) >= 4 && strings.ToLower(content[:4]) == "help" {
		return "set \n\t -avatar=qqnum  设置头像为qq号num的头像 \n\t -backgroundimage=urlhttp://xxx  设置背景大图为xxx图片 \n\t" + 
		" -signature=str 设置个性签名为str\n\n" , nil
	}
	if !ret{
		return "啊哦,不认识这个语法 请输入help查看支持命令" , nil
	}
	return "修改已提交", nil
}

func (s *SetPersonInfoRobot) parse(key string) map[string]string {
	infoMap := make(map[string]string)
	// 拆分每个键值对
	keyValuePairs := strings.Split(key, " -")
	for _, pair := range keyValuePairs {
		// 拆分键和值
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			infoMap[key] = value
		}
	}
	return infoMap
}

// 设置头像并返回用于发送给user的oss的key
func (s *SetPersonInfoRobot) ToSetAvatar(ctx context.Context, userId int64, qqnumber string , v... any) error {
	KqPusherClient , ok := v[0].(*kq.Pusher)
	if !ok{
		return ErrorKqPusher
	}
	FS , ok := v[1].(FileSystem.FileSystem)
	if !ok{
		return ErrorFileOss
	}
	key , err := s.getQQAvatar(qqnumber , FS)
	if err != nil || len(key) == 0{ 
		return err
	}
	message := make(map[string][]string)
	message[strconv.FormatInt(userId , 10)] = []string{"avatar" , key}
	data , err := json.Marshal(message)
	if err != nil{
		return err
	}
	//推送消息
	return KqPusherClient.Push(string(data))
}

// 设置背景图片并返回用于发送给user的oss的key
func (s *SetPersonInfoRobot) ToSetBackgroundImage(ctx context.Context, userId int64, url string , v... any) error {
	KqPusherClient , ok := v[0].(*kq.Pusher)
	if !ok{
		return ErrorKqPusher
	}
	FS , ok := v[1].(FileSystem.FileSystem)
	if !ok{
		return ErrorFileOss
	}
	key , err := s.getBackgrounImage(url , FS)
	if err != nil || len(key) == 0{ 
		return err
	}
	message := make(map[string][]string)
	message[strconv.FormatInt(userId , 10)] = []string{"backgroundImage" , key}
	data , err := json.Marshal(message)
	if err != nil{
		return err
	}
	//推送消息
	return KqPusherClient.Push(string(data))
}

// 设置个性签名 个性签名直接是string
func (t *SetPersonInfoRobot) ToSetSignature(ctx context.Context, userId int64, Signature string , v... any) error {
	KqPusherClient , ok := v[0].(*kq.Pusher)
	if !ok{
		return ErrorKqPusher
	}

	message := make(map[string][]string , 1)
	message[strconv.FormatInt(userId , 10)] = []string{"signature",Signature}

	data , err := json.Marshal(message)
	if err != nil{
		return nil
	}
	return KqPusherClient.Push(string(data))
}

func (t *SetPersonInfoRobot) getQQAvatar(qqnumber string , fs FileSystem.FileSystem) (string , error){
	// 构建 QQ 头像 API URL
	qqAvatarURL := fmt.Sprintf("http://q1.qlogo.cn/g?b=qq&nk=%s&s=100", qqnumber)

	// 发起 HTTP 请求获取头像
	resp, err := http.Get(qqAvatarURL)
	if err != nil {
		fmt.Println("Error:", err)
		return "",ErrorGetAvatar
	}
	defer resp.Body.Close()

	// //将头像数据保存到oss
	// err = fs.Upload(resp.Body , "avatar" , qqnumber)
	// if err != nil{
	// 	return "", err
	// }

	return filepath.Join("avatar", qqnumber) , nil
}

func (t *SetPersonInfoRobot) getBackgrounImage(url string ,fs FileSystem.FileSystem) (string , error){
	pattern := `^(https?|http)://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(pattern, url)
	fmt.Println(matched)
	if err != nil || !matched{
		return "" , ErrorParam
	}
	resp, err := http.Get(url)
	if err != nil {
		return "" , err
	}
	defer resp.Body.Close()

	// //将背景图片数据保存到oss
	// err = fs.Upload(resp.Body , "backgroundImage", url )
	// if err != nil{
	// 	return "", err
	// }

	return filepath.Join("backgroundImage", url) , nil
}