package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCancelByUuid(ctx context.Context, params orderV1.OrderCancelByUuidParams) (orderV1.OrderCancelByUuidRes, error) {
	err := a.service.CancelOrderByUuid(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderUUIDCannotBeEmpty) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrAlreadyPaid) || errors.Is(err, model.ErrOrderCancelled) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "The order has already been paid for or cancelled",
			}, nil
		}

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order by uuid '%s' not found", params.OrderUUID),
			}, nil
		}

		return nil, errors.New("an unexpected error occurred. Please try again later")
	}

	return &orderV1.OrderCancelByUuidNoContent{}, nil
}
