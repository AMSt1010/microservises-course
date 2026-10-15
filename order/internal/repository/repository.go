package repository

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

type OrderRepository interface {
	// Создает новый заказ в Map или обновляет информацию о заказе
	CreateOrder(ctx context.Context, order model.Order) error

	// Получает заказ из Map по Uuid заказа
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)

	// Помечает заказ как оплаченный с использованием сравнения и замены (Compare-And-Swap)
	PayOrder(ctx context.Context, orderUUID, transactionUUID string, method model.PaymentMethod) (model.Order, error)

	// CancelOrder отменяет заказ, если он еще не оплачен
	CancelOrder(ctx context.Context, orderUUID string) (model.Order, error)
}
