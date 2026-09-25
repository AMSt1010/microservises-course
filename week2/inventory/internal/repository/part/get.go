package part

import (
	"context"
	"log"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/converter"
)

// Возвращает информацию о детали по её UUID
func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if uuid == "" {
		return model.Part{}, model.ErrUUIDCannotBeEmpty
	}

	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	log.Printf("Деталь с Uuid %s получена: %+v", uuid, part)

	return converter.PartToModel(part), nil
}
