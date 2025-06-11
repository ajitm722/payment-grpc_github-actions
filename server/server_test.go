package main

import (
	"context"
	"testing"

	pb "github.com/ajitm722/payment-grpc_github-actions/proto"

	"github.com/stretchr/testify/assert"
)

func TestMakePaymentSuccess(t *testing.T) {
	// Initialize the PaymentServer
	paymentServer := NewPaymentServer()

	// Test case: Successful payment
	req := &pb.PaymentRequest{AmountCents: 50000}
	resp, err := paymentServer.MakePayment(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.TransactionId)
	assert.Equal(t, "SUCCESS", resp.Status)
}

func TestMakePaymentFailure(t *testing.T) {
	// Initialize the PaymentServer
	paymentServer := NewPaymentServer()

	// Test case: Failed payment
	req := &pb.PaymentRequest{AmountCents: 150000}
	resp, err := paymentServer.MakePayment(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.TransactionId)
	assert.Equal(t, "FAILED", resp.Status)
}

func TestGetPaymentStatusValidTransaction(t *testing.T) {
	// Initialize the PaymentServer
	paymentServer := NewPaymentServer()

	// Simulate a payment
	req := &pb.PaymentRequest{AmountCents: 50000}
	resp, _ := paymentServer.MakePayment(context.Background(), req)

	// Test case: Valid transaction ID
	statusReq := &pb.StatusRequest{TransactionId: resp.TransactionId}
	statusResp, err := paymentServer.GetPaymentStatus(context.Background(), statusReq)
	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", statusResp.Status)
}

func TestGetPaymentStatusInvalidTransaction(t *testing.T) {
	// Initialize the PaymentServer
	paymentServer := NewPaymentServer()

	// Test case: Invalid transaction ID
	statusReq := &pb.StatusRequest{TransactionId: "invalid-id"}
	statusResp, err := paymentServer.GetPaymentStatus(context.Background(), statusReq)
	assert.NoError(t, err)
	assert.Equal(t, "UNKNOWN", statusResp.Status)
}
