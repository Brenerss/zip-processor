package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoStorage struct {
	Client    *dynamodb.Client
	TableName string
}

func NewDynamoStorage() (*DynamoStorage, error) {
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

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("DYNAMO_ENDPOINT"))
	})

	return &DynamoStorage{
		Client:    client,
		TableName: "ImagensMetadata",
	}, nil
}

func (storage *DynamoStorage) SaveMetadata(metadata FileMetadata) error {
	item, err := attributevalue.MarshalMap(metadata)
	if err != nil {
		return fmt.Errorf("error to unpack metadata: %v", err)
	}

	_, err = storage.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(storage.TableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("erro to save on dynamodb: %v", err)
	}

	return nil
}
