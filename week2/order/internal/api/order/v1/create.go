package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/AMSt1010/microservises-course/week2/order/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	orderV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	if req == nil {
		return nil, a.GenerateError(http.StatusBadRequest, "request body cannot be nil")
	}
	order, err := a.service.CreateOrder(ctx, converter.CreateOrderToModel(*req))
	if err != nil {
		if errors.Is(err, model.ErrUserUUIDCannotBeEmpty) {
			return nil, a.GenerateError(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, model.ErrPartsUUIDsCannotBeEmpty) {
			return nil, a.GenerateError(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, model.ErrInventoryUnavailable) {
			return nil, a.GenerateError(http.StatusServiceUnavailable, "Inventory service is temporarily unavailable")
		}

		if errors.Is(err, model.ErrAtLeastOnePartWasNotFound) {
			return nil, a.GenerateError(http.StatusBadRequest, err.Error())
		}
		return nil, a.GenerateError(http.StatusInternalServerError, err.Error())
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
