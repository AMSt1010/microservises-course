package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

// Возвращает информацию о заказе.
func (s *service) GetOrderByUuid(ctx context.Context, orderUUID string) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	if orderUUID == "" {
		return model.Order{}, model.ErrOrderUUIDCannotBeEmpty
	}

	order, err := s.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
