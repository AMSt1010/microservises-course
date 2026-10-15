package v1

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

// Возвращает список деталей с возможностью фильтрации
func (a *api) ListParts(ctx context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	parts, err := a.invService.ListParts(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		return nil, err
	}
	return &invV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
