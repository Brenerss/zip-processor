package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client     *s3.Client
	BucketName string
}

func NewS3Storage() (*S3Storage, error) {
	credProvider := credentials.NewStaticCredentialsProvider(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("MINIO_ENDPOINT"))
		o.UsePathStyle = true
	})

	return &S3Storage{
		Client:     client,
		BucketName: os.Getenv("S3_BUCKET_NAME"),
	}, nil
}

func (s *S3Storage) DownloadFile(objectKey, destinationPath string) error {
	resp, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to search file in S3: %v", err)
	}

	defer resp.Body.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("an error occurred to create local file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("an error occurred to download stream: %v", err)
	}

	return nil
}
