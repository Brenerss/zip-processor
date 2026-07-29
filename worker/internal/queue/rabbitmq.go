package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/processorsystem/internal/processor"
	"github.com/processorsystem/internal/storage"
	"github.com/rabbitmq/amqp091-go"
)

type ProcessorTask struct {
	FileID   int64  `json:"id"`
	FilePath string `json:"path"`
}

func StartConsumer(db *storage.DynamoStorage, s3Client *storage.S3Storage) error {
	conn, err := amqp091.Dial("amqp://root:password@localhost:5672/")
	if err != nil {
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("image-processor", true, false, false, false, nil)
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	fmt.Printf("%v", msgs)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task ProcessorTask
			json.Unmarshal(d.Body, &task)

			fmt.Printf("task: %v", task)

			localZip := fmt.Sprintf("/tmp/%d.zip", task.FileID)

			err := s3Client.DownloadFile(task.FilePath, localZip)
			if err != nil {
				log.Printf("failed to download S3 file: %v", err)
				continue
			}

			err = processor.ProcessZip(localZip, task.FileID, db)
			if err != nil {
				log.Printf("failed to process zip file: %v", err)
			}

			os.Remove(localZip)
		}
	}()

	<-forever

	return nil
}
