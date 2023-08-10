package FileSystem

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"path/filepath"
	"time"
)

type S3 struct {
	URL    string
	Bucket string
	client *s3.Client
	ctx    context.Context
}

func NewS3(URL, Bucket, AwsAccessKeyId, AwsSecretAccessSecret string) *S3 {
	return NewS3Ctx(context.TODO(), URL, Bucket, AwsAccessKeyId, AwsSecretAccessSecret)
}
func NewS3Ctx(ctx context.Context, URL, Bucket, AwsAccessKeyId, AwsSecretAccessKe string) *S3 {
	var res S3
	res.ctx = ctx
	res.Bucket = Bucket
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "",
			URL:               URL,
			SigningRegion:     "",
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(res.ctx,
		config.WithRegion(""),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(AwsAccessKeyId, AwsSecretAccessKe, ""),
		),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic(err)
	}
	res.client = s3.NewFromConfig(cfg)
	_, err = res.bucketExists()
	if err != nil {
		panic(err)
	}
	return &res
}

// Upload 上传文件, 可覆盖
func (s *S3) Upload(file io.Reader, key ...string) error {
	_, err := s.client.PutObject(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filepath.Join(key...)),
		Body:   file,
	}, s3.WithAPIOptions(
		v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
	))
	if err != nil {
		return err
	}

	return nil
}

// GetDownloadLink 获取文件下载链接
func (s *S3) GetDownloadLink(key ...string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(filepath.Join(key...)),
		},
		s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}

// Delete 删除文件
func (s *S3) Delete(key ...string) error {
	_, err := s.client.DeleteObject(s.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filepath.Join(key...)),
	})
	if err != nil {
		return err
	}
	return nil
}

// bucketExists 检查是否存在桶
func (s *S3) bucketExists() (bool, error) {
	_, err := s.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				logc.Debugf(s.ctx, "Bucket %v is available.\n", s.Bucket)
				exists = false
				err = nil
			default:
				logc.Debugf(s.ctx, "Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", s.Bucket, err)
			}
		}
	} else {
		logc.Debugf(s.ctx, "Bucket %v exists and you already own it.", s.Bucket)
	}

	return exists, err
}

// FileExists 检查是否存在文件
func (s *S3) FileExists(key ...string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filepath.Join(key...)),
	})

	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				logc.Debugf(s.ctx, "Bucket %v is available.\n", s.Bucket)
				exists = false
				err = nil
			default:
				logc.Debugf(s.ctx, "Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", s.Bucket, err)
			}
		}
	} else {
		logc.Debugf(s.ctx, "Bucket %v exists and you already own it.", s.Bucket)
	}

	return exists, err
}
