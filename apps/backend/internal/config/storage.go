package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// StorageConfig works with local SeaweedFS and hosted S3-compatible providers.
type StorageConfig struct {
	Endpoint     string `koanf:"endpoint" validate:"required,url"`
	Region       string `koanf:"region" validate:"required"`
	Bucket       string `koanf:"bucket" validate:"required"`
	AccessKey    string `koanf:"access_key" validate:"required"`
	SecretKey    string `koanf:"secret_key" validate:"required"`
	UsePathStyle bool   `koanf:"use_path_style"`
}

func (c StorageConfig) Validate() error {
	if err := validator.New().Struct(c); err != nil {
		return fmt.Errorf("invalid storage configuration: %w", err)
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("storage endpoint must be an http(s) URL without credentials, query, or fragment")
	}
	return nil
}

// LoadStorageConfig lets the storage check run independently of PostgreSQL,
// Redis, Clerk, and the other application integrations.
func LoadStorageConfig() (*StorageConfig, error) {
	get := func(key string) string { return os.Getenv("TASK_MANAGER_STORAGE." + key) }
	pathStyle, err := strconv.ParseBool(get("USE_PATH_STYLE"))
	if err != nil {
		return nil, fmt.Errorf("TASK_MANAGER_STORAGE.USE_PATH_STYLE must be true or false")
	}
	c := &StorageConfig{Endpoint: get("ENDPOINT"), Region: get("REGION"), Bucket: get("BUCKET"), AccessKey: get("ACCESS_KEY"), SecretKey: get("SECRET_KEY"), UsePathStyle: pathStyle}
	return c, c.Validate()
}
