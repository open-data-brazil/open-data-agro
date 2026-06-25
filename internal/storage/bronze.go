package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/open-data-brazil/open-data-agro/internal/config"
)

// BronzeStore uploads, lists, and deletes bronze lake objects.
type BronzeStore interface {
	Put(ctx context.Context, key string, data []byte, contentType string) error
	Delete(ctx context.Context, key string) error
	ListPrefix(ctx context.Context, prefix string) ([]string, error)
	Backend() string
}

// NewBronzeStore selects the storage driver from STORAGE_MODE.
func NewBronzeStore(cfg config.Config) (BronzeStore, error) {
	switch cfg.StorageMode {
	case config.StorageModeLocal:
		return newLocalBronzeStore(cfg.LakeLocalRoot), nil
	case config.StorageModeMinIO:
		return newS3BronzeStore(s3BronzeConfig{
			endpoint:        cfg.MinIOEndpoint,
			accessKeyID:     cfg.MinIOAccessKey,
			secretAccessKey: cfg.MinIOSecretKey,
			bucket:          cfg.MinIOBucket,
			usePathStyle:    true,
		}, config.StorageModeMinIO)
	case config.StorageModeR2:
		return newS3BronzeStore(s3BronzeConfig{
			endpoint:        cfg.R2Endpoint,
			accessKeyID:     cfg.R2AccessKeyID,
			secretAccessKey: cfg.R2SecretAccessKey,
			bucket:          cfg.R2Bucket,
			usePathStyle:    false,
		}, config.StorageModeR2)
	default:
		return nil, fmt.Errorf("unsupported STORAGE_MODE %q", cfg.StorageMode)
	}
}

type localBronzeStore struct {
	root string
}

func newLocalBronzeStore(root string) BronzeStore {
	return &localBronzeStore{root: root}
}

func (s *localBronzeStore) Backend() string {
	return config.StorageModeLocal
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

func (s *localBronzeStore) ListPrefix(_ context.Context, prefix string) ([]string, error) {
	root := filepath.Join(s.root, filepath.FromSlash(strings.TrimSuffix(prefix, "/")))
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat prefix root: %w", err)
	}

	var keys []string
	walkRoot := filepath.Join(s.root, filepath.FromSlash(prefix))
	err := filepath.WalkDir(walkRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(s.root, path)
		if err != nil {
			return err
		}
		keys = append(keys, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list local prefix: %w", err)
	}
	return keys, nil
}

type s3BronzeConfig struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	bucket          string
	usePathStyle    bool
}

type s3BronzeStore struct {
	client  *s3.Client
	bucket  string
	backend string
}

func newS3BronzeStore(cfg s3BronzeConfig, backend string) (BronzeStore, error) {
	endpoint := strings.TrimRight(cfg.endpoint, "/")
	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(endpoint),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			cfg.accessKeyID,
			cfg.secretAccessKey,
			"",
		)),
		UsePathStyle: cfg.usePathStyle,
	})

	return &s3BronzeStore{
		client:  client,
		bucket:  cfg.bucket,
		backend: backend,
	}, nil
}

func (s *s3BronzeStore) Backend() string {
	return s.backend
}

func (s *s3BronzeStore) Put(ctx context.Context, key string, data []byte, contentType string) error {
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
		return fmt.Errorf("%s put object: %w", s.backend, err)
	}
	return nil
}

func (s *s3BronzeStore) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("%s delete object: %w", s.backend, err)
	}
	return nil
}

func (s *s3BronzeStore) ListPrefix(ctx context.Context, prefix string) ([]string, error) {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	var keys []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s list prefix: %w", s.backend, err)
		}
		for _, obj := range page.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}
	return keys, nil
}
