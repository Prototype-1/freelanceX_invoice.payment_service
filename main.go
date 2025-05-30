package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/config"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/handler"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/pkg"
	client "github.com/Prototype-1/freelanceX_invoice.payment_service/client"
	invoicepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/invoice_service"
	milestonePb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/milestone"
	paymentPb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/payment"
	"google.golang.org/grpc"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	dbConn, err := pkg.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	invoiceRepo := repository.NewInvoiceRepository(dbConn)
	profileClient := client.NewProfileServiceClient()
	timeTrackerClient := client.NewTimeServiceClient()

	milestoneRepo := repository.NewMilestoneRuleRepository(dbConn)
	milestoneService := service.NewMilestoneRuleService(milestoneRepo)

	paymentRepo := repository.NewPaymentRepository(dbConn)
	paymentService := service.NewPaymentService(paymentRepo, invoiceRepo, milestoneRepo)


invoiceService := service.NewInvoiceService(invoiceRepo)

invoiceHandler := handler.NewInvoiceHandler(
	invoiceRepo,
	invoiceService,
	milestoneService,
	cfg.KafkaBroker,
	cfg.InvoiceKafkaTopic, 
)

invoiceHandler.ProfileClient = profileClient
invoiceHandler.TimeTrackerClient = timeTrackerClient

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	grpcServer := grpc.NewServer()

	invoicepb.RegisterInvoiceServiceServer(grpcServer, invoiceHandler)
	milestoneRuleHandler := handler.NewMilestoneRuleHandler(milestoneService)
	milestonePb.RegisterMilestoneRuleServiceServer(grpcServer, milestoneRuleHandler)
	paymentHandler := handler.NewPaymentServiceServer(paymentService)
	paymentPb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	go func() {
		log.Printf("Starting gRPC server on port %s...", cfg.Port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
