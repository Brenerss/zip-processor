package queue

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/processorsystem/internal/processor"
	"github.com/processorsystem/internal/storage"
	"github.com/rabbitmq/amqp091-go"
)

type ProcessorTask struct {
	FileID   int64  `json:"id"`
	FilePath string `json:"path"`
}

type WebhookPayload struct {
	Status string `json:"status"`
}

func StartConsumer(db *storage.DynamoStorage, s3Client *storage.S3Storage) error {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task ProcessorTask

			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				d.Ack(false)
				continue
			}

			localZip := fmt.Sprintf("/tmp/%d.zip", task.FileID)

			err = s3Client.DownloadFile(task.FilePath, localZip)
			if err != nil {
				log.Printf("failed to download S3 file: %v", err)
				d.Nack(false, true)
				continue
			}

			err = processor.ProcessZip(localZip, task.FileID, db)
			if err != nil {
				log.Printf("failed to process zip file: %v", err)
				os.Remove(localZip)
				d.Nack(false, true)
				continue
			}

			os.Remove(localZip)

			success := notifyWebhook(task.FileID)
			if success {
				d.Ack(false)
			} else {
				log.Printf("critical error when notifying webhook, returning to the queue")
				d.Nack(false, true)
			}
		}
	}()

	<-forever

	return nil
}

func notifyWebhook(fileID int64) bool {
	payload := WebhookPayload{Status: "completed"}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to notify: %v", err)
		return false
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + "." + string(payloadBytes)
	secret := os.Getenv("WORKER_API_KEY")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	url := fmt.Sprintf("%s/attachments/%d/status", os.Getenv("WEBHOOK_BASE_ENDPOINT"), fileID)

	client := &http.Client{Timeout: 10 * time.Second}

	tries := 5
	backoff := 2 * time.Second

	for i := 1; i <= tries; i++ {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("an error occurred to send the HTTP: %v", err)
			return false
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Webhook-Timestamp", timestamp)
		req.Header.Set("X-Webhook-Signature", signature)

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			return true
		}

		statusError := "network/timeout error"
		if resp != nil {
			statusError = fmt.Sprintf("HTTP %d", resp.StatusCode)
			resp.Body.Close()
		}

		log.Printf("Webhook failed (tries %d/%d) - %s. Trying again in %v...", i, tries, statusError, backoff)
		time.Sleep(backoff)

		backoff *= 2
	}

	return false
}
