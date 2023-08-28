package video

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/video/rpc/video"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionRequest, r *http.Request) (resp *types.PublishActionResponse, err error) {

	// 参数检查
	if req.Token == "" {
		return &types.PublishActionResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	file, err := l.Upload(r, "data")
	if err != nil {
		return nil, err
	}
	mime := mimetype.Detect(file)

	if !strings.HasPrefix(mime.String(), "video/") {
		return nil, errors.New("不是视频")
	}

	name := uuid.New().String()
	err = l.uploadVideoToOSS(name, mime.Extension(), file)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.VideoRPC.SendPublishAction(l.ctx, &video.PublishActionReq{
		UserId:   tokenID,
		PlayUrl:  filepath.Join("video", name+mime.Extension()),
		CoverUrl: filepath.Join("img", name+".jpeg"),
		Title:    req.Title,
	})
	if err != nil {
		return nil, err
	}

	return &types.PublishActionResponse{
		RespStatus: types.RespStatus(apiVars.Success),
	}, nil
}

func (l *PublishActionLogic) uploadVideoToOSS(name, extension string, file []byte) error {
	img, err := ExampleReadFrameAsJpeg(bytes.NewReader(file))
	if err != nil {
		return err
	}
	err = l.svcCtx.FS.Upload(bytes.NewReader(file), "video", name+extension)
	if err != nil {
		return err
	}

	err = l.svcCtx.FS.Upload(img, "img", name+".jpeg")
	if err != nil {
		return err
	}
	logc.Infof(l.ctx, "视频%v上传成功", name)
	return nil
}

//const maxFileSize = 10 << 20 // 10 MB

func (l *PublishActionLogic) Upload(r *http.Request, key string) ([]byte, error) {
	//_ = r.ParseMultipartForm(maxFileSize)
	file, handler, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer func(file multipart.File) {
		_ = file.Close()

	}(file)
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	res := buf.Bytes()
	return res, nil
}

// ExampleReadFrameAsJpeg 获取视频略缩图，随机截取关键帧。
func ExampleReadFrameAsJpeg(inFile io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("pipe:0").
		Filter("select", ffmpeg.Args{"eq(pict_type\\,I)"}).
		Filter("random", ffmpeg.Args{"seed=42"}).
		Output("pipe:1", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput().
		WithInput(inFile).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
