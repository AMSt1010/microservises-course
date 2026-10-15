package part

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

// Записывает информацию о новой детали
func (s *service) CreatePart(ctx context.Context, info model.PartInfo) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if err := info.Validate(); err != nil {
		return "", err
	}

	newUUID := uuid.NewString()
	part := model.Part{
		UUID:      newUUID,
		Info:      info,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.invRepository.CreatePart(ctx, part); err != nil {
		return "", err
	}

	return newUUID, nil
}
