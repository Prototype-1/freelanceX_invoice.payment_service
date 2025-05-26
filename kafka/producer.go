package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/segmentio/kafka-go"
)

type InvoiceEvent struct {
	InvoiceID string `json:"invoice_id"`
	ClientID  string `json:"client_id"`
	EventType string `json:"event_type"`
}

func ProduceInvoiceEvent(broker, topic string, event InvoiceEvent) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	err = w.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Printf("Kafka write error: %v", err)
		return err
	}

	log.Printf("Produced invoice event to Kafka: %+v", event)
	return nil
}
