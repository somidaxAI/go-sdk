// Example: create and cancel an order using the Somidax Go SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/somidax/go-sdk/somidax"
)

func main() {
	client := somidax.New(os.Getenv("SOMIDAX_API_KEY"))
	ctx := context.Background()

	// Create an order
	order, err := client.Orders.Create(ctx, somidax.CreateOrderRequest{
		CustomerID: "cust_01J9X4MZABCD",
		Currency:   "GBP",
		Items: []somidax.CreateOrderItem{
			{ProductID: "prod_01J9TRAINERS", Quantity: 1},
			{ProductID: "prod_01J9CAP", VariantID: "var_black_m", Quantity: 2},
		},
	})
	if err != nil {
		log.Fatalf("create order: %v", err)
	}
	fmt.Printf("Order created: %s — %s\n", order.ID, order.Status)

	// Retrieve it
	fetched, err := client.Orders.Get(ctx, order.ID)
	if err != nil {
		log.Fatalf("get order: %v", err)
	}
	fmt.Printf("Fetched: %s  total=%d %s\n", fetched.ID, fetched.Total.Amount, fetched.Total.Currency)

	// Cancel it
	cancelled, err := client.Orders.Cancel(ctx, order.ID)
	if err != nil {
		log.Fatalf("cancel order: %v", err)
	}
	fmt.Printf("Cancelled: %s — %s\n", cancelled.ID, cancelled.Status)
}
