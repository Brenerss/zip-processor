# 🚀 Sistema Distribuído de Processamento de Imagens (Event-Driven)

![Laravel](https://img.shields.io/badge/API_Gateway-Laravel-FF2D20?style=for-the-badge&logo=laravel)
![Golang](https://img.shields.io/badge/Worker-Golang-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/Frontend-React_TS-20232A?style=for-the-badge&logo=react)
![RabbitMQ](https://img.shields.io/badge/Message_Broker-RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq)
![PostgreSQL](https://img.shields.io/badge/SQL-PostgreSQL-316192?style=for-the-badge&logo=postgresql)
![DynamoDB](https://img.shields.io/badge/NoSQL-DynamoDB-4053D6?style=for-the-badge&logo=amazon-dynamodb)
![MinIO](https://img.shields.io/badge/Storage-MinIO_(S3)-C7202C?style=for-the-badge&logo=minio)

## 📌 Sobre o Projeto
Este projeto é uma prova de conceito (PoC) de uma arquitetura baseada em eventos (Event-Driven Architecture) projetada para resolver um problema clássico de escalabilidade: **o processamento pesado de arquivos em requisições web**.

O sistema permite que os usuários façam upload de arquivos `.zip` contendo grandes lotes de imagens. Em vez de travar o servidor web descompactando e processando os arquivos de forma síncrona, a aplicação delega o trabalho pesado para um Worker assíncrono em alta performance, notificando o cliente em tempo real quando o processo é concluído.

## 🏗️ Desenho da Arquitetura

1. **API Gateway (Laravel):** Recebe o arquivo ZIP de forma extremamente rápida, faz o upload direto para o Storage (MinIO/S3), registra a transação no **PostgreSQL** e publica uma mensagem no **RabbitMQ** retornando um `202 Accepted`.
2. **Message Broker (RabbitMQ):** Atua como amortecedor (*Shock Absorber*), garantindo que picos de requisições não derrubem o servidor de processamento (Backpressure).
3. **Worker (Golang):** Consome a fila. Utiliza *Goroutines* para baixar o ZIP do S3, extrair e processar dezenas de imagens em paralelo.
4. **Persistência Poliglota:** Salva o status do processo no **PostgreSQL** (consistência ACID) e extrai os metadados dinâmicos de cada imagem (GPS, modelo da câmera, dimensões) para salvar no **DynamoDB** (NoSQL schema-less).
5. **Real-Time (SSE):** O Laravel mantém um túnel de Server-Sent Events aberto com o **React**, empurrando a atualização de status para a tela do usuário no milissegundo em que o Go finaliza o trabalho.

## 🧠 Decisões Arquiteturais (Trade-offs)
* **Por que não processar no Laravel?** Para evitar o esgotamento da memória RAM do servidor web e evitar erros de `504 Gateway Timeout` durante a descompactação de arquivos pesados.
* **Por que RabbitMQ?** Garante resiliência. Se o Worker em Go cair, o Laravel continua recebendo uploads normalmente. Quando o Worker voltar, ele processa as tarefas acumuladas sem perda de dados.
* **Por que DynamoDB junto com PostgreSQL?** Diferentes imagens possuem diferentes metadados. Evitamos criar uma tabela relacional esparsa (cheia de valores `NULL`) e delegamos dados não-estruturados para um banco NoSQL de alta performance.
* **Por que SSE e não WebSockets?** O fluxo exigia uma comunicação estritamente unidirecional (Servidor -> Cliente) para atualização de status. O SSE atende a essa necessidade consumindo muito menos recursos de rede e infraestrutura do que um túnel bidirecional WebSocket.

## ⚙️ Como Executar Localmente

### Pré-requisitos
* Docker e Docker Compose
* PHP 8.2+ e Composer
* Golang 1.21+
* Node.js 18+

### 1. Subir a Infraestrutura (Bancos, Fila e Storage)
```bash
docker compose up -d
