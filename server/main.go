package main

import (
	"fmt"
	"log"
	"net"

	pb "payment-grpc-app/proto"

	"google.golang.org/grpc"
)

func main() {
	// Step 1: Start a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Step 2: Create a gRPC server instance
	grpcServer := grpc.NewServer()

	// Step 3: Initialize and register the PaymentService
	paymentService := NewPaymentServer()
	pb.RegisterPaymentServiceServer(grpcServer, paymentService)

	// Step 4: Start serving incoming gRPC requests
	fmt.Println("PaymentService gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
