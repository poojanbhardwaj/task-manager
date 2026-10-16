package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(client *s3.Client, bucket string) *S3Client {
	return &S3Client{client: client, bucket: bucket}
}

func (s *S3Client) UploadFile(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fileKey := fmt.Sprintf("%s_%d", filepath.Base(fileName), time.Now().UnixNano())

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	contentType := http.DetectContentType(buffer.Bytes())
	if err := s.Put(ctx, fileKey, contentType, bytes.NewReader(buffer.Bytes())); err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return fileKey, nil
}

// Put uploads data to the configured bucket under key.
func (s *S3Client) Put(ctx context.Context, key, contentType string, body io.ReadSeeker) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

// Get returns a stream that the caller must close.
func (s *S3Client) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *S3Client) CreatePresignedUrl(ctx context.Context, objectKey string) (string, error) {
	presignedURL, err := s3.NewPresignClient(s.client).PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(time.Hour),
	)
	if err != nil {
		return "", err
	}
	return presignedURL.URL, nil
}

// DownloadURL creates a temporary URL for the private object.
func (s *S3Client) DownloadURL(ctx context.Context, key string) (string, error) {
	return s.CreatePresignedUrl(ctx, key)
}

func (s *S3Client) DeleteObject(ctx context.Context, objectKey string) error {
	return s.Delete(ctx, objectKey)
}

// Delete removes the object from the configured bucket.
func (s *S3Client) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	return nil
}
