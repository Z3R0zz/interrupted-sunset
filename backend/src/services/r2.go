package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var R2 *R2Service

const (
	region = "apac"
)

type R2Service struct {
	S3Client *s3.Client
	Bucket   string
}

func ConnectR2() error {
	accessKey := os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("CLOUDFLARE_R2_BUCKET")
	endpoint := os.Getenv("CLOUDFLARE_R2_ENDPOINT")

	if accessKey == "" || secretKey == "" || bucket == "" || endpoint == "" {
		return errors.New("missing one or more required environment variables for R2 service")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpoint
	})

	R2 = &R2Service{
		S3Client: s3Client,
		Bucket:   bucket,
	}

	return nil
}

func (r *R2Service) GetObject(ctx context.Context, key string) ([]byte, error) {
	output, err := r.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &r.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return data, nil
}

func (r *R2Service) UploadFile(ctx context.Context, key string, file *os.File) error {
	uploader := manager.NewUploader(r.S3Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &r.Bucket,
		Key:    &key,
		Body:   file,
	})
	return err
}

func (r *R2Service) GeneratePresignedURL(ctx context.Context, key string, expireIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(r.S3Client)

	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireIn))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.URL, nil
}
