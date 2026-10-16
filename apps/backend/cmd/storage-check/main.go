// storage-check exercises real storage without starting the whole application.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/config"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/lib/aws"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/server"
)

func check() error {
	cfg, err := config.LoadStorageConfig()
	if err != nil {
		return err
	}
	awsClient, err := aws.NewAWS(&server.Server{Config: &config.Config{Storage: cfg}})
	if err != nil {
		return err
	}
	store := awsClient.S3
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	key := "storage-check/" + uuid.NewString() + ".txt"
	const payload = "Task Manager storage is working.\n"
	if err := store.Put(ctx, key, "text/plain", strings.NewReader(payload)); err != nil {
		return fmt.Errorf("upload failed (is bun storage:start running?): %w", err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if err := store.Delete(cleanupCtx, key); err != nil {
			fmt.Fprintln(os.Stderr, "Could not clean up test object:", key, err)
		}
	}()
	body, err := store.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	data, err := io.ReadAll(body)
	body.Close()
	if err != nil || string(data) != payload {
		return fmt.Errorf("downloaded content did not match")
	}
	signedURL, err := store.DownloadURL(ctx, key)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(signedURL)
	if err != nil {
		return err
	}
	data, err = io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil || response.StatusCode != 200 || string(data) != payload {
		return fmt.Errorf("temporary download URL failed")
	}
	// Removing the signature must deny access to the private object.
	response, err = client.Get(strings.SplitN(signedURL, "?", 2)[0])
	if err != nil {
		return err
	}
	response.Body.Close()
	if response.StatusCode != http.StatusForbidden && response.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("unsigned download was not denied: HTTP %d", response.StatusCode)
	}
	if err := store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	if body, err := store.Get(ctx, key); err == nil {
		body.Close()
		return fmt.Errorf("deleted test object still exists")
	} else {
		var apiError smithy.APIError
		if !errors.As(err, &apiError) || apiError.ErrorCode() != "NoSuchKey" {
			return fmt.Errorf("could not confirm deletion: %w", err)
		}
	}
	fmt.Println("PASS: upload, download, temporary URL, private access, and deletion.")
	return nil
}

func main() {
	if err := check(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
