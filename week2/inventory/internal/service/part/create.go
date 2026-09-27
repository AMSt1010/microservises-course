package part

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// Записывает информацию о новой детали
func (s *service) CreatePart(ctx context.Context, info model.PartInfo) (string, error) {
	UUID, err := s.invRepository.CreatePart(ctx, info)
	if err != nil {
		return "", err
	}
	return UUID, nil
}
