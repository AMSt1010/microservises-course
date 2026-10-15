package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// ListParts возвращает список деталей с возможностью фильтрации.
func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	parts, err := s.invRepository.ListParts(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}
	return parts, nil
}
