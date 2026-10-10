package v1

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/order/internal/client/converter"
	"github.com/AMSt1010/microservises-course/week2/order/internal/model"
	inventoryV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	parts, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.PartsToModel(parts.GetParts()), nil
}
