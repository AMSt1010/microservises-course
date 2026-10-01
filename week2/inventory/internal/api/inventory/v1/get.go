package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/converter"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	invV1 "github.com/AMSt1010/microservises-course/week2/shared/pkg/proto/inventory/v1"
)

// Возвращает информацию о детали по её UUID
func (a *api) GetPart(ctx context.Context, req *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	part, err := a.invService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrUUIDCannotBeEmpty) {
			return nil, status.Error(codes.InvalidArgument, "uuid cannot be empty")
		}
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "Part with UUID: %v not found", req.GetUuid())
		}
		return nil, err
	}
	return &invV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
