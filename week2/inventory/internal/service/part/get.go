package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// Записывает информацию о новой детали
func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.invRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
