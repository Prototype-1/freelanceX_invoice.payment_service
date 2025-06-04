package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	Port        string
	KafkaBroker        string
	InvoiceKafkaTopic  string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	AppConfig = Config{
		DBHost:      getEnv("DB_HOST", "postgres"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "freelanceX_invoice_service"),
		Port:        getEnv("PORT", "50051"),
		KafkaBroker:        getEnv("KAFKA_BROKER", "kafka:9092"),
		InvoiceKafkaTopic:  getEnv("INVOICE_KAFKA_TOPIC", "invoice-events"),
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
