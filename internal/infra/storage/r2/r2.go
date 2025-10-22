package r2

import (
	"context"
	"fmt"
	"io"
	"time"

	appcfg "github.com/prawirdani/golang-restapi/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Cloudflare R2 Storage, using AWS S3 Sdk
type R2 struct {
	bucket string
	client *s3.Client
}

func NewR2Storage(cfg appcfg.R2Config) (*R2, error) {
	r2Cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(r2Cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(
			fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
		)
	})
	return &R2{
		bucket: cfg.Bucket,
		client: client,
	}, nil
}

// Put implements storage.Storage.
func (r *R2) Put(
	ctx context.Context,
	path string,
	reader io.Reader,
	contentType string,
) error {
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	return err
}

// Get implements storage.Storage.
func (r *R2) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// Delete implements storage.Storage.
func (r *R2) Delete(ctx context.Context, path string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	return err
}

// Exists implements storage.Storage.
func (r *R2) Exists(ctx context.Context, path string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		// TODO: Check if it's a "not found" error
		return false, nil
	}
	return true, nil
}

// GetURL implements storage.Storage.
func (r *R2) GetURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(r.client)
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}
