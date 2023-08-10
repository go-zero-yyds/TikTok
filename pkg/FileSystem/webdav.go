package FileSystem

import (
	"github.com/studio-b12/gowebdav"
	"io"
	"path/filepath"
)

type Webdav struct {
	client             *gowebdav.Client
	downloadLinkPrefix string
	prefix             string
}

func New(URL, user, password, prefix, downloadLinkPrefix string) *Webdav {
	c := gowebdav.NewClient(URL, user, password)
	return &Webdav{
		client:             c,
		downloadLinkPrefix: downloadLinkPrefix,
		prefix:             prefix,
	}
}

// Upload 上传文件, 可覆盖
func (s *Webdav) Upload(file io.Reader, key ...string) error {
	dirPath := key[:len(key)-1]
	err := s.client.MkdirAll(filepath.Join(s.prefix, filepath.Join(dirPath...)), 0644)
	if err != nil {
		return err
	}
	err = s.client.WriteStream(filepath.Join(s.prefix, filepath.Join(key...)), file, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetDownloadLink 获取文件下载链接
func (s *Webdav) GetDownloadLink(key ...string) (string, error) {
	return s.downloadLinkPrefix + filepath.Join(s.prefix, filepath.Join(key...)), nil
}

// Delete 删除文件
func (s *Webdav) Delete(key ...string) error {
	err := s.client.Remove(filepath.Join(s.prefix, filepath.Join(key...)))
	if err != nil {
		return err
	}
	return nil
}

// FileExists 检查是否存在文件
func (s *Webdav) FileExists(key ...string) (bool, error) {
	_, err := s.client.Stat(filepath.Join(s.prefix, filepath.Join(key...)))
	if err != nil {
		return false, nil
	}
	return true, nil
}
