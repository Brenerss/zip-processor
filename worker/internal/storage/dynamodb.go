package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoStorage struct {
	Client    *dynamodb.Client
	TableName string
}

func NewDynamoStorage() (*DynamoStorage, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8002")
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
