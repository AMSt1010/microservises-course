package part

import (
	"context"
	"log"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	"github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	if err := ctx.Err(); err != nil {
		return model.Part{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	log.Printf("Деталь с Uuid %s получена: %+v", uuid, part)

	return converter.PartToModel(part), nil
}
