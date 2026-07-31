# Distributed Image Processing System (WIP)

![Laravel](https://img.shields.io/badge/API_Gateway-Laravel-FF2D20?style=for-the-badge&logo=laravel)
![Golang](https://img.shields.io/badge/Worker-Golang-00ADD8?style=for-the-badge&logo=go)
![RabbitMQ](https://img.shields.io/badge/Message_Broker-RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq)
![PostgreSQL](https://img.shields.io/badge/SQL-PostgreSQL-316192?style=for-the-badge&logo=postgresql)
![DynamoDB](https://img.shields.io/badge/NoSQL-DynamoDB-4053D6?style=for-the-badge&logo=amazon-dynamodb)
![MinIO](https://img.shields.io/badge/Storage-MinIO_(S3)-C7202C?style=for-the-badge&logo=minio)

## 📌 About the Project
This project is a Proof of Concept (PoC) of an Event-Driven Architecture designed to solve a classic scalability problem: **handling heavy file processing in web requests**.

The system allows users to upload `.zip` files containing large batches of images. Instead of locking up the web server by unzipping and processing the files synchronously, the application delegates the heavy lifting to a high-performance asynchronous Worker, notifying the client in real-time once the process is complete.

## 🏗️ Architecture Design

1. **API Gateway (Laravel):** Receives the ZIP file incredibly fast, uploads it directly to the Storage (MinIO/S3), registers the transaction in **PostgreSQL**, and publishes a message to **RabbitMQ** returning a `202 Accepted` response.
2. **Message Broker (RabbitMQ):** Acts as a shock absorber, applying backpressure to ensure that request spikes do not crash the processing server.
3. **Worker (Golang):** Consumes the queue. Utilizes *Goroutines* to download the ZIP from S3, extract, and process dozens of images in parallel.
4. **Polyglot Persistence:** Saves the process status in **PostgreSQL** (ACID consistency) and extracts dynamic metadata from each image (GPS, camera model, dimensions) to save in **DynamoDB** (schema-less NoSQL).
5. **Real-Time (SSE):** Laravel keeps a Server-Sent Events tunnel open with **React**, pushing the status update to the user's screen in the exact millisecond the Go worker finishes the job.

## 🧠 Architectural Decisions (Trade-offs)
* **Why not process it in Laravel?** To prevent exhausting the web server's RAM and avoid `504 Gateway Timeout` errors when unzipping and processing heavy files.
* **Why RabbitMQ?** It guarantees resilience. If the Go Worker crashes, Laravel continues to receive uploads normally. When the Worker is back online, it processes the accumulated tasks without any data loss.
* **Why DynamoDB alongside PostgreSQL?** Different images have different metadata. We avoided creating a sparse relational table (filled with `NULL` values) and delegated the unstructured data to a high-performance NoSQL database.
* **Why SSE and not WebSockets?** The workflow required strictly unidirectional communication (Server -> Client) for status updates. SSE meets this need while consuming significantly fewer network and infrastructure resources than a bidirectional WebSocket tunnel.

## ⚙️ How to Run Locally

### Prerequisites
* Docker and Docker Compose
* PHP 8.2+ and Composer
* Golang 1.21+
* Node.js 18+

### 1. Spin up the Infrastructure (Databases, Queue, and Storage)
```bash
docker compose up -d
