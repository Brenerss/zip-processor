package main

import (
	"log"

	"github.com/processorsystem/internal/queue"
	"github.com/processorsystem/internal/storage"
)

func main() {
	log.Println("Starting worker processor...")

	db, err := storage.NewDynamoStorage()
	if err != nil {
		log.Fatalf("failed to connect dynamodb: %v", err)
	}

	s3Client, err := storage.NewS3Storage()
	if err != nil {
		log.Fatalf("failed to connect S3 Client: %v", err)
	}

	err = queue.StartConsumer(db, s3Client)
	if err != nil {
		log.Fatalf("failed to start rabbitmq: %v", err)
	}
}
