package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	if err := ctx.Err(); err != nil {
		return model.Part{}, err
	}

	if uuid == "" {
		return model.Part{}, model.ErrUUIDCannotBeEmpty
	}

	part, err := s.invRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
