package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AMSt1010/microservises-course/week2/payment/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/payment/internal/model"
	PaymentV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/payment/v1"
)

// Обрабатывает команду на оплату и возвращает `transaction_uuid`.
func (a *api) PayOrder(ctx context.Context, req *PaymentV1.PayOrderRequest) (*PaymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, converter.PaymentToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, "all payment fields must be filled")
		}

		// 2. Проверяем ошибку контекста (таймаут или отмена клиентом)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, status.FromContextError(err).Err()
		}

		return nil, err
	}

	return &PaymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
