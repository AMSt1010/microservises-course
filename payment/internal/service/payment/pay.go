package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if err := payment.Validate(); err != nil {
		return "", err
	}

	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно , transaction_uuid: %v\nUUID заказа: %v\nUUID пользователя: %v\nСпособ оплаты: %v\n", transactionUUID, payment.OrderUUID, payment.UserUUID, payment.PaymentMethod)
	return transactionUUID, nil
}
