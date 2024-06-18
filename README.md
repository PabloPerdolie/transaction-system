# Transaction System

This repository contains a transaction system service built using GoLang, Kafka, Docker, and PostgreSQL. It is designed to handle high-throughput transaction processing with robust scalability.

## Features

- High-throughput transaction processing
- Microservices architecture
- Kafka for messaging
- PostgreSQL for persistent storage
- Docker for containerization

## Technologies Used

- **Programming Language**: Go (Golang)
- **Messaging**: Kafka
- **Database**: PostgreSQL
- **Containerization**: Docker

## Getting Started
### To start kafka and zookeeper
```bash
    docker compose up 
```
### To start API-Gateway
```bash
    go run ./api-gateway/main/main
```
### To start transaction-system
```bash
    go run ./transaction-system/cmd/main/main
```
### Prerequisites

- Docker and Docker Compose
- Go 1.16+
- Kafka and Zookeeper

### Installation

Clone the repository:
```bash
   git clone https://github.com/PabloPerdolie/transaction-system.git
```
