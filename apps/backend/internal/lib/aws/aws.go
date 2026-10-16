package aws

import (
	"fmt"
	"net/http"
	"time"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/server"
)

type AWS struct {
	S3 *S3Client
}

// NewAWS builds the configured S3-compatible client for this application.
func NewAWS(server *server.Server) (*AWS, error) {
	if server == nil || server.Config == nil || server.Config.Storage == nil {
		return nil, fmt.Errorf("storage is not configured")
	}
	storageConfig := server.Config.Storage
	if err := storageConfig.Validate(); err != nil {
		return nil, err
	}

	client := s3.New(s3.Options{
		Region:       storageConfig.Region,
		BaseEndpoint: sdkaws.String(storageConfig.Endpoint),
		UsePathStyle: storageConfig.UsePathStyle,
		Credentials: sdkaws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			storageConfig.AccessKey,
			storageConfig.SecretKey,
			"",
		)),
		HTTPClient:       &http.Client{Timeout: 30 * time.Second},
		RetryMaxAttempts: 3,
	})

	return &AWS{S3: NewS3Client(client, storageConfig.Bucket)}, nil
}
