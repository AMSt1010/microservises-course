package service

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.CreateOrderIn) (model.Order, error)
	PayOrder(ctx context.Context, pay model.PayOrderIn) (string, error)
	GetOrderByUuid(ctx context.Context, orderUUID string) (model.Order, error)
	CancelOrderByUuid(ctx context.Context, orderUUID string) error
}
