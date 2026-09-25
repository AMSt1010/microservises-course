package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// Записывает информацию о новой детали
func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.invRepository.ListParts(ctx, filter)
	if err != nil {
		return []model.Part{}, nil
	}
	return parts, nil
}
