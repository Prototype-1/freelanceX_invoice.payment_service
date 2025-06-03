# Invoice & Payment Service — FreelanceX

## Overview
[Add overview here: responsible for generating invoices for completed projects, tracking payments, and managing transaction history between clients and freelancers.]

## Tech Stack
- Go (Golang)
- gRPC
- PostgreSQL
- Protocol Buffers
- Razorpay test
- Kafka producer(invoice)

## Folder Structure

.
├── config/ # Config loader and .env setup
├── handler/ # gRPC service implementations
├── model/ # DB models for invoice and payments
├── proto/ # Protobuf service definitions
├── repository/ # Data access layer
├── service/ # Business logic layer
├── kafka/ # DB migrations or utilities
├── main.go # Service entrypoint
└── go.mod


## Setup

### 1. Clone & Navigate
```bash
git clone github.com/Prototype-1/freelanceX_invoice.payment_service.git
cd freelancex_project/invoice.payment_service
```

## Install Dependencies

go mod tidy

### Create .env File

PORT=50056
DB_URL=postgres://username:password@localhost:5432/invoicedb
RAZORPAY_KEY_ID
RAZORPAY_KEY_SECRET
KAFKA_BROKER=localhost
INVOICE_KAFKA_TOPIC

### Run Migrations

go run scripts/migrate.go

## Start the Service

go run main.go

## Proto Definitions

Located in proto/invoice/invoice.proto

Regenerate:

protoc --go_out=. --go-grpc_out=. proto/invoice/invoice.proto

### Notes

Each invoice is linked to a completed project and freelancer.

Payments can be tracked manually or integrated with payment gateways (e.g. Stripe, Razorpay).

Future support for auto-reminders and partial payments planned.


#### Maintainers

aswin100396@gmail.com

