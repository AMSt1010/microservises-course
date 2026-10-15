package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AMSt1010/microservises-course/week2/order/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

// POST /api/v1/orders/{order_uuid}/pay
// Проводит оплату ранее созданного заказа.
func (a *api) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	if req == nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "The request body is invalid or malformed",
		}, nil
	}
	transactionUUID, err := a.service.PayOrder(ctx, model.PayOrderIn{
		PaymentMethod: converter.PaymentMethodToModel(req.GetPaymentMethod()),
		OrderUUID:     params.OrderUUID,
	})
	if err != nil {
		if errors.Is(err, model.ErrOrderUUIDCannotBeEmpty) || errors.Is(err, model.ErrInsufficientPaymentMethod) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order by uuid '%s' not found", params.OrderUUID),
			}, nil
		}

		if errors.Is(err, model.ErrAlreadyPaid) || errors.Is(err, model.ErrOrderCancelled) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "The order has already been paid for or cancelled",
			}, nil
		}

		if errors.Is(err, model.ErrPaymentUnavailable) {
			return &orderV1.BadGatewayError{
				Code:    http.StatusBadGateway,
				Message: "Bad Gateway: Invalid response from upstream service",
			}, nil
		}

		return nil, errors.New("an unexpected error occurred. Please try again later")
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
