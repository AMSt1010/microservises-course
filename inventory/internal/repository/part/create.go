package part

import (
	"context"
	"log"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	repoConverter "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/converter"
)

// Записывает информацию о новой детали
func (r *repository) CreatePart(ctx context.Context, part model.Part) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	repoPart := repoConverter.PartToRepoModel(part)
	r.data[repoPart.UUID] = repoPart

	partsFilter := newPartsFilterStorage(repoPart)
	r.filterIndices[repoPart.UUID] = partsFilter

	log.Printf("Сохранена деталь с Uuid %v", repoPart.UUID)
	return nil
}
