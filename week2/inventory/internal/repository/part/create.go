package part

import (
	"context"
	"log"
	"time"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
	repoConverter "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/converter"
	repoModel "github.com/AMSt1010/microservises-course/week2/inventory/internal/repository/model"
	"github.com/google/uuid"
)

// Записывает информацию о новой детали
func (r *repository) CreatePart(ctx context.Context, info model.PartInfo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newUUID := uuid.NewString()

	part := repoModel.Part{
		UUID:      newUUID,
		Info:      repoConverter.PartInfoToRepoModel(info),
		CreatedAt: time.Now(),
	}

	r.data[newUUID] = part

	partsFilter := newPartsFilterStorage(part)
	r.filterIndices[newUUID] = partsFilter

	log.Printf("Создана деталь с Uuid %v", newUUID)
	return newUUID, nil
}
