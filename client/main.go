package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "payment-grpc-app/proto"
)

// makePayment sends a payment request to the server and returns the transaction ID.
func makePayment(client pb.PaymentServiceClient, user string, amountCents int64) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.MakePayment(ctx, &pb.PaymentRequest{
		UserId:      user,
		AmountCents: amountCents,
	})
	if err != nil {
		log.Fatalf("MakePayment failed for user %s: %v", user, err)
	}

	fmt.Printf("Payment submitted for user %s\nAmount: $%.2f\nTransaction ID: %s\nStatus: %s\n\n",
		user, float64(amountCents)/100, resp.TransactionId, resp.Status)

	return resp.TransactionId
}

// confirmPaymentStatus checks and logs the current status of a transaction.
func confirmPaymentStatus(client pb.PaymentServiceClient, txID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.GetPaymentStatus(ctx, &pb.StatusRequest{
		TransactionId: txID,
	})
	if err != nil {
		log.Fatalf("GetPaymentStatus failed for transaction %s: %v", txID, err)
	}

	fmt.Printf("Confirmed status for transaction %s: %s\n\n", txID, resp.Status)
}

// queryTransaction performs a standalone status lookup for any transaction ID.
func queryTransaction(client pb.PaymentServiceClient, txID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.GetPaymentStatus(ctx, &pb.StatusRequest{
		TransactionId: txID,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to query transaction %s: %v", txID, err)
		return
	}

	fmt.Printf("Queried status for transaction %s: %s\n\n", txID, resp.Status)
}

func main() {
	// Connect to the gRPC server running locally
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	// 1. Make and confirm a successful payment
	txID1 := makePayment(client, "ajit", 2499) // $24.99
	confirmPaymentStatus(client, txID1)
	queryTransaction(client, txID1)

	// 2. Make and confirm a failed payment (amount over $1000)
	txID2 := makePayment(client, "bob", 150000) // $1500.00
	confirmPaymentStatus(client, txID2)
	queryTransaction(client, txID2)

	// 3. Query a non-existent transaction ID
	queryTransaction(client, "non-existent-id-1234")
}

