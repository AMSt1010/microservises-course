package grpc

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, payment model.Payment) (string, error)
}
