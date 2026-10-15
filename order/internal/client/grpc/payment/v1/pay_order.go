package v1

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/client/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
)

func (c *client) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := payment.Validate(); err != nil {
		return "", err
	}

	paymentResponse, err := c.generatedClient.PayOrder(ctx, converter.PaymentToProto(payment))
	if err != nil {
		return "", err
	}

	return paymentResponse.GetTransactionUuid(), nil
}
