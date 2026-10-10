package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
)

// Получает заказ из Map по Uuid заказа
func (r *repository) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	r.mu.RLock()
	order, ok := r.orders[orderUUID]
	r.mu.RUnlock()
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}
