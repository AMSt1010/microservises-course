package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	"github.com/AMSt1010/microservises-course/week2/order/internal/repository/converter"
)

// Создает новый заказ в Map или обновляет информацию о заказе
func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	repoOrder := converter.OrderToRepoModel(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[repoOrder.OrderUUID]; ok {
		return model.ErrAlreadyExists
	}
	r.orders[repoOrder.OrderUUID] = repoOrder
	return nil
}
