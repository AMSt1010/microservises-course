package service

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (string, error)
}
