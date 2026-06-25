package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/open-data-brazil/open-data-agro/internal/config"
)

// BronzeStore uploads and deletes bronze parquet objects.
type BronzeStore interface {
	Put(ctx context.Context, key string, data []byte, contentType string) error
	Delete(ctx context.Context, key string) error
	Backend() string
}

// NewBronzeStore returns R2 when configured, otherwise a local lake fallback.
func NewBronzeStore(cfg config.Config) (BronzeStore, error) {
	if cfg.R2Enabled() {
		return newR2BronzeStore(cfg)
	}
	return newLocalBronzeStore(cfg.LakeLocalRoot), nil
}

type localBronzeStore struct {
	root string
}

func newLocalBronzeStore(root string) BronzeStore {
	return &localBronzeStore{root: root}
}

func (s *localBronzeStore) Backend() string {
	return "local"
}

func (s *localBronzeStore) Put(_ context.Context, key string, data []byte, _ string) error {
	path := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create bronze dir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write bronze file: %w", err)
	}
	return nil
}

func (s *localBronzeStore) Delete(_ context.Context, key string) error {
	path := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete bronze file: %w", err)
	}
	return nil
}

type r2BronzeStore struct {
	client *s3.Client
	bucket string
}

func newR2BronzeStore(cfg config.Config) (BronzeStore, error) {
	endpoint := strings.TrimRight(cfg.R2Endpoint, "/")
	client := s3.New(s3.Options{
		Region: "auto",
		BaseEndpoint: aws.String(endpoint),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			cfg.R2AccessKeyID,
			cfg.R2SecretAccessKey,
			"",
		)),
		UsePathStyle: true,
	})

	return &r2BronzeStore{
		client: client,
		bucket: cfg.R2Bucket,
	}, nil
}

func (s *r2BronzeStore) Backend() string {
	return "r2"
}

func (s *r2BronzeStore) Put(ctx context.Context, key string, data []byte, contentType string) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("r2 put object: %w", err)
	}
	return nil
}

func (s *r2BronzeStore) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2 delete object: %w", err)
	}
	return nil
}
