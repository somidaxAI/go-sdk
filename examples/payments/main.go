// Example: settle a payment and verify a webhook using the Somidax Go SDK.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/somidax/go-sdk/somidax"
)

func main() {
	client := somidax.New(os.Getenv("SOMIDAX_API_KEY"))
	ctx := context.Background()

	// Settle a payment for an existing order
	payment, err := client.Payments.Create(ctx, somidax.CreatePaymentRequest{
		OrderID:  "ord_01J9EXAMPLE",
		Amount:   4999, // £49.99 in pence
		Currency: "GBP",
		Method:   somidax.PaymentMethodFiat,
	})
	if err != nil {
		log.Fatalf("create payment: %v", err)
	}
	fmt.Printf("Payment %s — status: %s\n", payment.ID, payment.Status)

	// Start a minimal webhook listener
	secret := os.Getenv("SOMIDAX_WEBHOOK_SECRET")
	http.HandleFunc("/webhooks/somidax", func(w http.ResponseWriter, r *http.Request) {
		body, err := somidax.ReadBody(r)
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}

		if err := somidax.VerifySignature(r, body, secret); err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var payload somidax.WebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received event: %s  id=%s\n", payload.Event, payload.ID)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Webhook listener on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
