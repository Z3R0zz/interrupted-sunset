package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		return nil, fmt.Errorf("failed to get object from R2: %w", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return data, nil
}
