package storage

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/config"
)

// Store keeps files in a private bucket. Handlers must check task ownership
// before uploading, downloading, deleting, or issuing a temporary download URL.
type Store struct {
	client *s3.Client
	bucket string
}

func New(cfg config.StorageConfig) (*Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client := s3.New(s3.Options{
		Region:           cfg.Region,
		BaseEndpoint:     aws.String(cfg.Endpoint),
		UsePathStyle:     cfg.UsePathStyle,
		Credentials:      aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		HTTPClient:       &http.Client{Timeout: 30 * time.Second},
		RetryMaxAttempts: 3,
	})
	return &Store{client: client, bucket: cfg.Bucket}, nil
}

func (s *Store) Put(ctx context.Context, key, contentType string, body io.ReadSeeker) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key), ContentType: aws.String(contentType), Body: body})
	return err
}

// Get returns a stream that the caller must close.
func (s *Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	return err
}

func (s *Store) DownloadURL(ctx context.Context, key string) (string, error) {
	result, err := s3.NewPresignClient(s.client).PresignGetObject(ctx,
		&s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)},
		s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
