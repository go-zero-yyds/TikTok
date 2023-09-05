package robot

import (
	"TikTok/pkg/FileSystem"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/oauth2"

	"github.com/hbagdi/go-unsplash/unsplash"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logc"
)

var (
	ErrorKqPusher  = errors.New("not KqPusher type")
	ErrorGetAvatar = errors.New("get qq avatar error from http")

	ErrorFileOss = errors.New("oss error")
	ErrorParam   = errors.New("param error")
)

const token = "Client-ID x_q4MjED1PCQY4mcHgiFZ1pxSxl7nP_fk3UICVGa01s"

// 功能是 接收用户信息参数，返回需要发送给user的数据包

// 格式是json包的map[string][]string   【uid】{修改字段，key}
// 类型：avatar           头像
//
//		    backgroundimage 背景图
//	     signature       个性签名

type SetPersonInfoRobot struct {
	prologue string
	unsplash *unsplash.Unsplash
	client   *http.Client
}

func NewSetPersonInfoRobot(KqPusherClient *kq.Pusher) (int64, *SetPersonInfoRobot) {
	message := make(map[int64][]string)
	message[1] = []string{"username", "🤖抖音1号"}
	data, err := json.Marshal(message)
	if err != nil {
		panic("robots start error")
	}
	err = KqPusherClient.Push(string(data))
	if err != nil {
		panic("robots start error")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	clnt := oauth2.NewClient(context.TODO(), ts)
	unpas := unsplash.New(clnt)

	return 1, &SetPersonInfoRobot{
		prologue: "滴滴...我是一个能修改头像、背景大图、个性签名的机器人, 请输入help查看命令, 我会尽力为您服务...",
		client:   clnt,
		unsplash: unpas,
	}
}

func (t *SetPersonInfoRobot) DisplayPrologue() string {
	return t.prologue
}

func (t *SetPersonInfoRobot) run(ctx context.Context, userId int64, _ int64, content string, v ...any) (string, error) {
	data, err := t.deal(ctx, userId, content, v...)
	if err != nil {
		return "", err
	}
	return data, nil
}

// 返回发送给userid的信息
func (t *SetPersonInfoRobot) deal(ctx context.Context, userId int64, content string, v ...any) (string, error) {
	ret := false
	if len(content) >= 3 && strings.ToLower(content[:3]) == "set" {
		infoMap := t.parse(content)
		//修改头像提取qq号 支持qq
		if value, ok := infoMap["avatar"]; ok {
			if value[0:2] == "qq" && len(v) >= 2 {
				err := t.ToSetAvatar(userId, value[2:], v...)
				if err != nil {
					logc.Error(ctx, userId, err)
					return "设置错误,请检查参数是否正确,或等会再试", nil
				}
				ret = true
			}
		}

		//设置背景大图 支持qq
		if value, ok := infoMap["backgroundimage"]; ok {
			if value == "random" {
				err := t.ToSetBackgroundImage(ctx, userId, v...)
				if err != nil {
					logc.Error(ctx, userId, err)
					return "设置错误,请检查参数是否正确,或等会再试", nil
				}
				ret = true
			}
		}

		//修改个性签名
		if value, ok := infoMap["signature"]; ok {
			err := t.ToSetSignature(userId, value, v...)
			if err != nil {
				logc.Error(ctx, userId, err)
				return "设置错误,请检查参数是否正确,或等会再试", nil
			}
			ret = true
		}
	}
	if len(content) >= 4 && strings.ToLower(content[:4]) == "help" {
		return "set \n\t --avatar=qq num  设置头像为qq号num的头像 \n\t --backgroundimage=random 设置背景大图为随机图片\n\t" +
			" --signature=str 设置个性签名为str\n\n", nil
	}
	if !ret {
		return "啊哦,不认识这个语法 请输入help查看支持命令", nil
	}
	return "设置成功", nil
}

func (t *SetPersonInfoRobot) parse(key string) map[string]string {
	infoMap := make(map[string]string)
	// 拆分每个键值对
	keyValuePairs := strings.Split(key, " --")
	for _, pair := range keyValuePairs {
		// 拆分键和值
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := strings.ReplaceAll(parts[0], " ", "")
			value := strings.ReplaceAll(parts[1], " ", "")
			infoMap[key] = value
		}
	}
	return infoMap
}

// ToSetAvatar 设置头像并发送给user的oss的key
func (t *SetPersonInfoRobot) ToSetAvatar(userId int64, qqnumber string, v ...any) error {
	if ok, _ := regexp.MatchString("^[0-9]+$", qqnumber); !ok {
		return ErrorParam
	}
	KqPusherClient, ok := v[0].(*kq.Pusher)
	if !ok {
		return ErrorFileOss
	}
	FS, ok := v[1].(FileSystem.FileSystem)
	if !ok {
		return ErrorFileOss
	}
	key, err := t.getQQAvatar(qqnumber, FS)
	if err != nil || len(key) == 0 {
		return err
	}
	message := make(map[string][]string)
	message[strconv.FormatInt(userId, 10)] = []string{"avatar", key}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//推送消息
	return KqPusherClient.Push(string(data))
}

// ToSetBackgroundImage 设置背景图片并发送给user的oss的key
func (t *SetPersonInfoRobot) ToSetBackgroundImage(ctx context.Context, userId int64, v ...any) error {
	KqPusherClient, ok := v[0].(*kq.Pusher)
	if !ok {
		return ErrorKqPusher
	}
	FS, ok := v[1].(FileSystem.FileSystem)
	if !ok {
		return ErrorFileOss
	}

	photos, _, err := t.unsplash.Photos.Random(nil)
	if photos == nil || len(*photos) != 1 {
		logc.Error(ctx, err)
	}
	photo := (*photos)[0]
	url := photo.Links.Download.String()

	key, err := t.getBackgrounImage(url, FS)
	if err != nil || len(key) == 0 {
		return err
	}
	message := make(map[string][]string)
	message[strconv.FormatInt(userId, 10)] = []string{"backgroundImage", key}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//推送消息
	return KqPusherClient.Push(string(data))
}

// ToSetSignature 设置个性签名 个性签名直接是string
func (t *SetPersonInfoRobot) ToSetSignature(userId int64, Signature string, v ...any) error {
	KqPusherClient, ok := v[0].(*kq.Pusher)
	if !ok {
		return ErrorKqPusher
	}

	message := make(map[string][]string, 1)
	message[strconv.FormatInt(userId, 10)] = []string{"signature", Signature}

	data, err := json.Marshal(message)
	if err != nil {
		return nil
	}
	return KqPusherClient.Push(string(data))
}

func toRound(inputFile io.Reader, outputFileFormat string) ([]byte, error) {
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, err
	}

	// 创建一个新的空白图像，大小与原始图像相同
	bounds := img.Bounds()
	outputImg := image.NewRGBA(bounds)

	// 计算圆形的半径
	radius := bounds.Dx() / 2

	// 循环遍历像素，将非圆形部分设置为透明
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			dx := x - bounds.Min.X - radius
			dy := y - bounds.Min.Y - radius
			distance := dx*dx + dy*dy

			if distance <= radius*radius {
				outputImg.Set(x, y, img.At(x, y))
			} else {
				outputImg.Set(x, y, color.Transparent)
			}
		}
	}

	// 创建输出文件
	//outputFileName := "output." + outputFileFormat
	//outputFile, err := os.Create(outputFileName)
	if err != nil {
		return nil, err
	}
	//defer outputFile.Close()
	buf := bytes.NewBuffer(nil)

	// 根据输出图像格式保存图像
	switch outputFileFormat {
	case "png":
		err = png.Encode(buf, outputImg)
	case "jpeg":
		err = jpeg.Encode(buf, outputImg, nil)
	default:
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *SetPersonInfoRobot) getQQAvatar(qqnumber string, fs FileSystem.FileSystem) (string, error) {
	// 构建 QQ 头像 API URL
	qqAvatarURL := fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%s&s=640", qqnumber)

	// 发起 HTTP 请求获取头像
	resp, err := http.Get(qqAvatarURL)
	if err != nil {
		return "", ErrorGetAvatar
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logx.Errorf("HTTP请求返回非200状态码:", resp.Status)
		return "", errors.New("HTTP请求返回非200状态码")
	}
	// 读取响应主体数据到内存缓冲
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}
	if err != nil {
		logx.Errorf("Error reading response body:", err)
		return "", err
	}
	// 检查
	mime := mimetype.Detect(buf.Bytes())

	if !strings.HasPrefix(mime.String(), "image/") {
		return "", errors.New("不是图片")
	}
	p, err := toRound(buf, "png")
	if err != nil {
		return "", err
	}
	//将图片转成sha1码
	// 将图片转成 SHA-256 哈希值
	sha := sha1.New()
	sha.Write(p)
	sha1Value := fmt.Sprintf("%x", sha.Sum(nil))

	key := filepath.Join("avatar", sha1Value)
	//如果oss不存在这个图片
	if ok, _ := fs.FileExists(key); !ok {
		// 将头像数据上传至OSS
		err = fs.Upload(bytes.NewReader(p), "avatar", sha1Value)
		if err != nil {
			return "", err
		}
	}
	return key, nil
}

func (t *SetPersonInfoRobot) getBackgrounImage(url string, fs FileSystem.FileSystem) (string, error) {
	pattern := `^(https?|http)://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil || !matched {
		return "", ErrorParam
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logx.Errorf("HTTP请求返回非200状态码:", resp.Status)
		return "", errors.New("HTTP请求返回非200状态码")
	}
	// 读取响应主体数据到内存缓冲
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}
	if err != nil {
		logx.Errorf("Error reading response body:", err)
		return "", err
	}
	// 检查
	mime := mimetype.Detect(buf.Bytes())

	if !strings.HasPrefix(mime.String(), "image/") {
		return "", errors.New("不是图片")
	}

	//将图片转成sha1码
	// 将图片转成 SHA-1 哈希值
	sha := sha1.New()
	sha.Write(buf.Bytes())
	sha1Value := fmt.Sprintf("%x", sha.Sum(nil))

	key := filepath.Join("backgroundImage", sha1Value)
	// 如果oss不存在这个图片
	if ok, _ := fs.FileExists(key); !ok {
		// 将头像数据上传至OSS
		err = fs.Upload(bytes.NewReader(buf.Bytes()), "backgroundImage", sha1Value)
		if err != nil {
			return "", err
		}
	}
	return key, nil
}
