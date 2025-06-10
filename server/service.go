package main

import (
	"context"
	"sync"

	pb "payment-grpc-app/proto"

	"github.com/google/uuid"
)

// Constants for payment statuses (avoids string duplication)
const (
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
	StatusUnknown = "UNKNOWN"
)

// PaymentServer implements the PaymentService gRPC server interface.
type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	mu          sync.Mutex
	statusStore map[string]string // transaction ID -> status
}

// NewPaymentServer returns a new PaymentServer instance.
func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		statusStore: make(map[string]string),
	}
}

// MakePayment simulates a payment transaction.
// If amount_cents ≤ 100000 (i.e., $1000), the payment is marked as SUCCESS.
// If amount_cents > 100000, it's marked as FAILED.
func (s *PaymentServer) MakePayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	txID := uuid.New().String()
	status := StatusSuccess // assume success by default

	if req.AmountCents > 100000 {
		status = StatusFailed
	}

	s.statusStore[txID] = status

	return &pb.PaymentResponse{
		TransactionId: txID,
		Status:        status,
	}, nil
}

// GetPaymentStatus returns the status of a transaction.
// If the transaction ID doesn't exist, it returns "UNKNOWN".
func (s *PaymentServer) GetPaymentStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, ok := s.statusStore[req.TransactionId]
	if !ok {
		status = StatusUnknown
	}
	return &pb.StatusResponse{Status: status}, nil
}
