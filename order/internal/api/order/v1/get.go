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

func (a *api) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, err := a.service.GetOrderByUuid(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderUUIDCannotBeEmpty) {
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
		return nil, a.GenerateError(http.StatusInternalServerError, err.Error())
	}
	return &orderV1.GetOrderResponse{
		Order: converter.OrderToAPIOrder(order),
	}, nil
}
