package repository

import (
	"context"

	"github.com/AMSt1010/microservises-course/week2/inventory/internal/model"
)

type InventoryRepository interface {
	// Записывает информацию о новой детали
	CreatePart(ctx context.Context, info model.PartInfo) (string, error)
	// Возвращает информацию о детали по её UUID
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	// Возвращает список деталей с возможностью фильтрации
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
