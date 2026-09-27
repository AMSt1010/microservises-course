package v1

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

// Записывает информацию о новой детали
func (a *api) CreatePart(ctx context.Context, req *invV1.CreatePartRequest) (*invV1.CreatePartResponse, error) {
	UUID, err := a.InvService.CreatePart(ctx, converter.PartInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}
	return &invV1.CreatePartResponse{
		Uuid: UUID,
	}, nil
}
