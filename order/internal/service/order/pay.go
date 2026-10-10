package order

import (
	"context"
	"fmt"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

// Проводит оплату ранее созданного заказа.
func (s *service) PayOrder(ctx context.Context, pay model.PayOrderIn) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if err := pay.Validate(); err != nil {
		return "", err
	}

	order, err := s.repo.GetOrder(ctx, pay.OrderUUID)
	if err != nil {
		return "", err
	}

	if order.Status != nil {
		if *order.Status == model.OrderStatusPAID {
			return "", model.ErrAlreadyPaid
		}

		if *order.Status == model.OrderStatusCANCELLED {
			return "", model.ErrOrderCancelled
		}
	}

	payment := model.Payment{
		OrderUUID:     order.OrderUUID,
		UserUUID:      order.UserUUID,
		PaymentMethod: pay.PaymentMethod,
	}

	transactionUuid, err := s.paymentClient.PayOrder(ctx, payment)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("%w: %w", model.ErrPaymentUnavailable, err)
	}

	_, err = s.repo.PayOrder(ctx, pay.OrderUUID, transactionUuid, pay.PaymentMethod)
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
