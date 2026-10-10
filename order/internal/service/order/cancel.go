package order

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

// Отменяет заказ.
func (s *service) CancelOrderByUuid(ctx context.Context, orderUUID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if orderUUID == "" {
		return model.ErrOrderUUIDCannotBeEmpty
	}

	order, err := s.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status != nil && *order.Status == model.OrderStatusPAID {
		return model.ErrAlreadyPaid
	}

	if order.Status != nil && *order.Status == model.OrderStatusCANCELLED {
		return model.ErrOrderCancelled
	}

	_, err = s.repo.CancelOrder(ctx, orderUUID)
	if err != nil {
		return err
	}

	return nil
}
