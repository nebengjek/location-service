package minio

import (
	"context"
	"location-service/bin/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	endpoint        string
	accessKey       string
	secretAccessKey string
	useSSL          bool
	minioClient     *minio.Client
)

type MinioClient struct{}

func NewMinio() MinioClient {
	return MinioClient{}
}

func InitMinio() {
	endpoint = config.GetConfig().MinioEndpoint
	accessKey = config.GetConfig().MinioAccessKey
	secretAccessKey = config.GetConfig().MinioSecretKey
	useSSL = config.GetConfig().MinioUseSSL

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		panic(err)
	}

	minioClient = client
}

type UploadObject struct {
	BucketName  string
	ObjectName  string
	FilePath    string
	ContentType string
}

func (s *MinioClient) UploadObject(ctx context.Context, m UploadObject) (result string, err error) {
	info, err := minioClient.FPutObject(ctx, m.BucketName, m.ObjectName, m.FilePath, minio.PutObjectOptions{ContentType: m.ContentType})
	if err != nil {
		return
	}

	result = info.Key
	return
}

type DownloadObject struct {
	BucketName string
	ObjectName string
	SavingPath string
}

func (s *MinioClient) DownloadObject(ctx context.Context, m DownloadObject) error {
	err := minioClient.FGetObject(ctx, m.BucketName, m.ObjectName, m.SavingPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

type RemoveObject struct {
	BucketName string
	ObjectName string
}

func (s *MinioClient) RemoveObject(ctx context.Context, m RemoveObject) error {
	err := minioClient.RemoveObject(ctx, m.BucketName, m.ObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

type IsBucketExists struct {
	BucketName string
}

func (s *MinioClient) IsBucketExists(ctx context.Context, m IsBucketExists) (result bool, err error) {
	result, err = minioClient.BucketExists(ctx, m.BucketName)
	if err != nil {
		return
	}

	return
}

type CreateBucket struct {
	BucketName string
}

func (s *MinioClient) CreateBucket(ctx context.Context, m CreateBucket) error {
	err := minioClient.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return err
	}

	return nil
}

func GetMinioClient() *minio.Client {
	return minioClient
}
