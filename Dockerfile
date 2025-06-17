FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o invoice_payment_service main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/invoice_payment_service .
COPY .env .

EXPOSE 50056

CMD ["./invoice_payment_service"]
