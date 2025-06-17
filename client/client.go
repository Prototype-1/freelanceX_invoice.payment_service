package client

import (
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	profilePb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/user_service"
	timePb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/timeTracker_service"
)

func NewProfileServiceClient() profilePb.ProfileServiceClient {
	conn, err := grpc.NewClient("freelancex_user_service.default.svc.cluster.local:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	return profilePb.NewProfileServiceClient(conn)
}

func NewTimeServiceClient() timePb.TimeLogServiceClient {
	conn, err := grpc.NewClient("freelancex_time_tracker_service.default.svc.cluster.local:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	return timePb.NewTimeLogServiceClient(conn)
}
