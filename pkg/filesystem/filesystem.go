package filesystem

import "io"

type FileSystem interface {
	Upload(file io.Reader, key ...string) error
	GetDownloadLink(key ...string) (string, error)
	Delete(key ...string) error
	FileExists(key ...string) (bool, error)
}
