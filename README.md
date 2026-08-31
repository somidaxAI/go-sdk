# Somidax Go SDK

Official Go client for the [Somidax API](https://somidax.net/developers).

[![CI](https://github.com/somidaxAI/go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/somidaxAI/go-sdk/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/somidaxAI/go-sdk/somidax.svg)](https://pkg.go.dev/github.com/somidaxAI/go-sdk/somidax)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Installation

```bash
go get github.com/somidaxAI/go-sdk
```

Requires Go 1.22+.

## Quick Start

```go
import "github.com/somidaxAI/go-sdk/somidax"

client := somidax.New("sk_live_your_api_key")

order, err := client.Orders.Create(ctx, somidax.CreateOrderRequest{
    CustomerID: "cust_01J9X4MZABCD",
    Currency:   "GBP",
    Items: []somidax.CreateOrderItem{
        {ProductID: "prod_01J9TRAINERS", Quantity: 1},
    },
})
```

## Services

| Service | Description |
|---|---|
| `client.Orders` | Create, retrieve, list, update, cancel orders |
| `client.Catalog` | Manage products and variants |
| `client.Payments` | Somidax Pay settlements and refunds |
| `client.Rewards` | $SMDX token earn, redeem, balances |
| `client.Analytics` | Merchant KPIs and top products |
| `client.Webhooks` | Register endpoints and verify signatures |

## Configuration

```go
client := somidax.New(
    "sk_live_your_api_key",
    somidax.WithBaseURL("https://sandbox.api.somidax.net/v1"),
    somidax.WithTimeout(10 * time.Second),
)
```

## Webhook Verification

```go
body, _ := somidax.ReadBody(r)
if err := somidax.VerifySignature(r, body, webhookSecret); err != nil {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
}
```

## Error Handling

```go
var apiErr *somidax.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
}
```

## Token Contracts

| Chain | Address |
|---|---|
| Ethereum ERC-20 | `0x7e8539D1E5cB91d63E46B8e188403b3f262a949B` |
| BNB Chain BEP-20 | `0xea8c5b9c537f3ebbcc8f2df0573f2d084e9e2bdb` |

## Links

- [API Reference](https://somidax.net/developers)
- [Windows SDK](https://github.com/somidaxAI/go-sdk-windows)
- [TypeScript SDK](https://github.com/somidaxAI/sdk-typescript)
- [Python SDK](https://github.com/somidaxAI/sdk-python)

## License

MIT © 2026 Somidax
